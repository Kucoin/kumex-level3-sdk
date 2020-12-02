# DEPRECATED(弃用)

This repository is deprecated, please use [https://github.com/Kucoin/kucoin-level3-sdk](https://github.com/Kucoin/kucoin-level3-sdk) instead.

---

# KuMEX Level3 market

## guide
  [guide](docs/guide_CN.md)

## Installation

1. install dependencies

```
go get github.com/JetBlink/orderbook
go get github.com/go-redis/redis
go get github.com/gorilla/websocket
go get github.com/joho/godotenv
go get github.com/pkg/errors
go get github.com/shopspring/decimal
```

2. build

```
CGO_ENABLED=0 go build -ldflags '-s -w' -o kumex_market kumex_market.go
```

or you can download the latest available release.

## Usage

1. vim `.env`:
    ```
    # enable debug log
    API_DEBUG_MODE=true
    
    API_BASE_URI=https://api.kumex.com
    
    ENABLE_ORDER_BOOK=true
    
    ENABLE_EVENT_WATCHER=true
    REDIS_HOST=127.0.0.1:6379
    
    RPC_TOKEN=your-rpc-token
    ```

1. Run Command：

    ```
    ./kumex_market -c .env -symbol XBTUSDM -p 9090 -rpckey XBTUSDM
    ```
    

## Docker Usage

1. Build docker image

   ```
   docker build -t kumex_market .
   ```

1. [vim .env](#usage)

1. Run

  ```
  docker run --rm -it -v $(pwd)/.env:/app/.env --net=host kumex_market
  ```

## RPC Method

> endpoint : 127.0.0.1:9090
> the sdk rpc is based on golang jsonrpc 1.0 over tcp.

see:[python jsonrpc client demo](./demo/python-demo/Level3/rpc.py)

* Get Part Order Book
    ```
    {"method": "Server.GetPartOrderBook", "params": [{"token": "your-rpc-token", "number": 1}], "id": 0}
    ```
    
* Get Full Order Book
    ```
    {"method": "Server.GetOrderBook", "params": [{"token": "your-rpc-token"}], "id": 0}
    ```

* Add Event ClientOids To Channels
    ```
    {"method": "Server.AddEventClientOidsToChannels", "params": [{"token": "your-rpc-token", "data": {"clientOid": ["channel-1", "channel-2"]}}], "id": 0}
    ```

* Add Event OrderIds To Channels
    ```
    {"method": "Server.AddEventOrderIdsToChannels", "params": [{"token": "your-rpc-token", "data": {"orderId": ["channel-1", "channel-2"]}}], "id": 0}
    ```
## Python-Demo

> the demo including orderbook display、L3 monitor order and sand order ,you can add your own functionality on this basis or rewriting strategy.

see:[python use_level3 demo](./demo/python-demo/demo)

- Run KuMEXOrderBook.py
    ```
    command: python KuMEXOrderBook.py
    describe: display orderbook
    ```
- Run orderMonitor.py
    ```
    command: python orderMonitor.py
    describe: use level3 (Add Event ClientOids To Channels) monitor yourself order,if your order is match,It will be sold at a profitable price
    ```
- Run demo.py
    ```
    command: python demo.py
    describe: you sand order always buy1,until your order is match.(easy strategy)
    ```

