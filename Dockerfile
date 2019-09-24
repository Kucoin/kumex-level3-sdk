FROM golang:1.13-stretch as builder

RUN export GO111MODULE=on \
    && export GOPROXY=https://goproxy.io \
    && mkdir -p /go/src/github.com/Kucoin/kumex-level3-sdk

COPY . /go/src/github.com/Kucoin/kumex-level3-sdk

RUN cd /go/src/github.com/Kucoin/kumex-level3-sdk \
    && CGO_ENABLED=0 go build -ldflags '-s -w' -o /go/bin/kumex_market kumex_market.go

FROM debian:stretch

RUN apt-get update \
    && apt-get install ca-certificates -y

COPY --from=builder /go/bin/kumex_market /usr/local/bin/

# .env => /app/.env
WORKDIR /app
VOLUME /app

EXPOSE 9090

COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["kumex_market", "-c", "/app/.env", "-symbol", "XBTUSDM", "-p", "9090", "-rpckey", "XBTUSDM"]
