package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/Kucoin/kumex-market/builder"
	"github.com/Kucoin/kumex-market/events"
	"github.com/Kucoin/kumex-market/utils/log"
)

//Server is api server
type Server struct {
	level3Builder *builder.Builder
	eventWatcher  *events.Watcher

	apiPort string
	token   string
}

//InitRpcServer init rpc server
func InitRpcServer(apiPort string, token string, level3Builder *builder.Builder, watcher *events.Watcher) {
	if apiPort == "" || token == "" {
		panic(fmt.Sprintf("missing configuration，apiPort: %s, token: %s", apiPort, token))
	}

	apiPort = ":" + apiPort
	if err := rpc.Register(&Server{
		level3Builder: level3Builder,
		eventWatcher:  watcher,

		apiPort: apiPort,
		token:   token,
	}); err != nil {
		panic("rpc service registration failed")
	}

	log.Warn("start running rpc server, port: %s", apiPort)

	listener, err := net.Listen("tcp", apiPort)
	if err != nil {
		panic("api server run failed, error: %s" + err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}

//TokenMessage is token type message
type TokenMessage struct {
	Token string `json:"token"`
}

//Response is api response
type Response struct {
	Code  string      `json:"code"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func (s *Server) checkToken(token string) string {
	if token != s.token {
		return s.failure(TokenErrorCode, "error token")
	}

	return ""
}

func (s *Server) success(data interface{}) string {
	response, _ := json.Marshal(&Response{
		Code:  "0",
		Data:  data,
		Error: "",
	})

	return string(response)
}

const (
	ServerErrorCode = "10"
	TokenErrorCode  = "20"
	TickerErrorCode = "30"
)

func (s *Server) failure(code string, err string) string {
	response, _ := json.Marshal(&Response{
		Code:  code,
		Data:  "",
		Error: err,
	})

	return string(response)
}
