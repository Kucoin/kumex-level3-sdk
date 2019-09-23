package main

import (
	"flag"

	"github.com/Kucoin/kumex-market/app"
	"github.com/Kucoin/kumex-market/utils/log"
	"github.com/joho/godotenv"
)

func main() {
	envFile := flag.String("c", ".env", ".env file")
	symbol := flag.String("symbol", "", "SYMBOL")
	rpcPort := flag.String("p", "", "rpc port")
	rpcKey := flag.String("rpckey", "", "market maker redis rpckey")

	flag.Parse()

	loadEnv(*envFile)

	app.NewApp(*symbol, *rpcPort, *rpcKey).Run()
}

func loadEnv(file string) {
	err := godotenv.Load(file)
	if err != nil {
		log.Error("Error loading .env file: %s", file)
	}
}
