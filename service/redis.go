package service

import (
	"github.com/Kucoin/kumex-level3-sdk/utils/log"
	"github.com/go-redis/redis"
)

type Redis struct {
	redisPool *redis.Client
}

const RedisKeyPrefix = "kumexMarket:rpcKey:"

func NewRedis(addr, password string, db int, rpcKey string, symbol string, rpcPort string) *Redis {
	if addr == "" {
		return nil
	}

	log.Warn("connect to redis: " + addr)
	redisPool := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
		//DialTimeout:  10 * time.Second,
		//ReadTimeout:  30 * time.Second,
		//WriteTimeout: 30 * time.Second,
		//PoolSize:     10,
		//PoolTimeout:  30 * time.Second,
	})

	//redisKey如: kucoinMarket:rpcKey:KCS-USDT:rpcKey
	if err := redisPool.Set(RedisKeyPrefix+symbol+":"+rpcKey, rpcPort, 0).Err(); err != nil {
		panic("connect to redis failed: " + err.Error())
	}

	return &Redis{
		redisPool: redisPool,
	}
}

func (r *Redis) Publish(channel string, message interface{}) error {
	return r.redisPool.Publish(channel, message).Err()
}
