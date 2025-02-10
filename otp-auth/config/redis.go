package config

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic("Redis connection failed: " + err.Error())
	}
	fmt.Println("✅ Redis connected!")
}
