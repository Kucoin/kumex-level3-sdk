# DEPRECATED(弃用)

此项目不再维护，请采用 [https://github.com/Kucoin/kucoin-level3-sdk](https://github.com/Kucoin/kucoin-level3-sdk)

---

# KuMEX Level3 market

## 入门文档
  [入门文档](docs/guide_CN.md)

## 安装

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

或者直接下载已经编译完成的二进制文件

## 用法

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

1. 运行命令：

    ```
    ./kumex_market -c .env -symbol XBTUSDM -p 9090 -rpckey XBTUSDM
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

> python的demo包含了一个本地orderbook的展示，一个简单的下单并且追踪订单的策略.

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
