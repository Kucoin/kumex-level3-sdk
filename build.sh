#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o kumex-market-for-linux kumex_market.go

CGO_ENABLED=0 go build -ldflags '-s -w' -o kumex-market-for-mac kumex_market.go

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o kumex-market-for-windows.exe kumex_market.go

