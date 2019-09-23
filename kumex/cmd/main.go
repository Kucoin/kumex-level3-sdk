package main

import (
	"log"
	"time"

	"github.com/JetBlink/orderbook/base"
	"github.com/Kucoin/kumex-market/kumex"
	"github.com/Kucoin/kumex-market/utils/recovery"
)

func main() {
	defer recovery.Recover(func(stack string) {
		log.Println(stack)
	})
	apiService := kumex.NewKuMEX("http://api.kumex.com", true, 30*time.Second)
	//if err := apiService.AtomicFullOrderBook("XBTUSDM"); err != nil {
	//	panic(err)
	//}

	tk, err := apiService.WebSocketPublicToken()
	if err != nil {
		panic(err)
	}

	c := apiService.NewWebSocketClient(tk)

	mc, ec, err := c.Connect()
	if err != nil {
		panic(err)
	}

	ch := kumex.NewSubscribeMessage("/contractMarket/level3:XBTUSDM", false)
	if err := c.Subscribe(ch); err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-ec:
			c.Stop() // Stop subscribing the WebSocket feed
			panic(err)

		case msg := <-mc:
			log.Println(base.ToJsonString(msg))
		}
	}
}
