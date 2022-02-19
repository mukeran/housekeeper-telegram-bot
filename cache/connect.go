package cache

import (
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

var (
	Redis *redis.Pool
)

func Connect() {
	Redis = &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			redisAddress := os.Getenv("REDIS_ADDR")
			if redisAddress == "" {
				redisAddress = "127.0.0.1:6379"
			}
			return redis.Dial("tcp", redisAddress)
		},
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
	}
}
