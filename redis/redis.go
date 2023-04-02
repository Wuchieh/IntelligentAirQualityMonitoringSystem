package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var (
	CTX = context.Background()

	Redis *redis.Client

	Reset = "\033[0m"
	Green = "\033[32m"
)

type setting struct {
	RedisAddr     string `json:"redisAddr"`
	RedisPassword string `json:"redisPassword"`
	RedisDB       int    `json:"redisDB"`
}

func init() {
	var s setting
	if file, err := os.ReadFile("setting.json"); err != nil {
		panic(err)
	} else {
		err = json.Unmarshal(file, &s)
		if err != nil {
			panic(err)
		}
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:     s.RedisAddr,
		Password: s.RedisPassword, // no password set
		DB:       s.RedisDB,       // use default DB
	})

	ping := Redis.Ping(CTX)
	_, err := ping.Result()

	if err != nil {
		panic(err)
	}
	log.Println(Green, "Redis is connect", Reset)
}
