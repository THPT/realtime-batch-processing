package infra

import (
	"fmt"
	"realtime-batch-processing/config"

	"gopkg.in/redis.v5"
)

var Redis *redis.Client

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisHost + ":" + config.Config.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := Redis.Ping().Result()
	if err != nil {
		panic(err)
	}
}
func CloseRedis() {
	if Redis != nil {
		if err := Redis.Close(); err != nil {
			fmt.Println("[ERROR] Cannot close Redis connection, err:", err)
		}
	}
}
