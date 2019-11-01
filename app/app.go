package app

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Kucoin/kumex-level3-sdk/api"
	"github.com/Kucoin/kumex-level3-sdk/builder"
	"github.com/Kucoin/kumex-level3-sdk/events"
	"github.com/Kucoin/kumex-level3-sdk/kumex"
	"github.com/Kucoin/kumex-level3-sdk/service"
	"github.com/Kucoin/kumex-level3-sdk/utils/log"
)

type App struct {
	debug bool

	apiService *kumex.KuMEX
	symbol     string

	enableOrderBook bool
	level3Builder   *builder.Builder

	enableEventWatcher bool
	eventWatcher       *events.Watcher
	redisPool          *service.Redis

	rpcPort  string
	rpcToken string
}

func NewApp(symbol string, rpcPort string, rpcKey string) *App {
	if symbol == "" {
		panic("symbol is required")
	}

	if rpcPort == "" {
		panic("rpcPort is required")
	}

	redisHost := os.Getenv("REDIS_HOST")

	if rpcKey == "" && redisHost != "" {
		panic("rpckey is required")
	}

	apiService := kumex.NewKuMEX(os.Getenv("API_BASE_URI"), false, 30*time.Second)
	level3Builder := builder.NewBuilder(apiService, symbol)

	var redisPassword = os.Getenv("REDIS_PASSWORD")

	redisPool := service.NewRedis(redisHost, redisPassword, rpcKey, symbol, rpcPort)
	eventWatcher := events.NewWatcher(redisPool)

	return &App{
		debug: os.Getenv("API_DEBUG_MODE") == "true",

		apiService: apiService,
		symbol:     symbol,

		enableOrderBook: os.Getenv("ENABLE_ORDER_BOOK") == "true",
		level3Builder:   level3Builder,

		enableEventWatcher: os.Getenv("ENABLE_EVENT_WATCHER") == "true",
		redisPool:          redisPool,
		eventWatcher:       eventWatcher,

		rpcPort:  rpcPort,
		rpcToken: os.Getenv("RPC_TOKEN"),
	}
}

func (app *App) Run() {
	if app.enableOrderBook {
		go app.level3Builder.ReloadOrderBook()
	}

	if app.enableEventWatcher {
		go app.eventWatcher.Run()
	}

	//rpc server
	go api.InitRpcServer(app.rpcPort, app.rpcToken, app.level3Builder, app.eventWatcher)

	app.websocket()
}

func (app *App) writeMessage(msgRawData json.RawMessage) {
	if app.debug {
		log.Info("raw message : %s", kumex.ToJsonString(msgRawData))
	}
	if app.enableOrderBook {
		app.level3Builder.Messages <- msgRawData
	}

	if app.enableEventWatcher {
		app.eventWatcher.Messages <- msgRawData
	}

	const msgLenLimit = 50
	if len(app.level3Builder.Messages) > msgLenLimit ||
		len(app.eventWatcher.Messages) > msgLenLimit {
		log.Error(
			"msgLenLimit: app.level3Builder.Messages: %d, app.eventWatcher.Messages: %d, app.verify.Messages: %d",
			len(app.level3Builder.Messages),
			len(app.eventWatcher.Messages),
		)
	}
}

func (app *App) websocket() {
	apiService := app.apiService

	tk, err := apiService.WebSocketPublicToken()
	if err != nil {
		panic(err)
	}

	c := apiService.NewWebSocketClient(tk)

	mc, ec, err := c.Connect()
	if err != nil {
		panic(err)
	}

	ch := kumex.NewSubscribeMessage("/contractMarket/level3:"+app.symbol, false)
	if err := c.Subscribe(ch); err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-ec:
			c.Stop()
			panic(err)

		case msg := <-mc:
			app.writeMessage(msg.RawData)
		}
	}
}
