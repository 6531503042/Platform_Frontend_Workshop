package utils

import (
	"backend/config"
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// Global Redis client
var RedisClient *redis.Client

// InitRedis initializes the Redis client
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
		DB:       config.RedisDB,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	log.Println("Connected to Redis!")
}