package config

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() *redis.Client {
	addr := AppConfig.RedisAddr

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	log.Println("Connected to Redis successfully")
	RedisClient = client
	return client
}
