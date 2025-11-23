package database

import (
	"context"
	"fmt"
	"lomi-backend/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func ConnectRedis(cfg *config.Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	fmt.Println("✅ Connected to Redis")
}

// PublishMessage publishes a message to a Redis channel
func PublishMessage(channel string, message interface{}) error {
	return RedisClient.Publish(ctx, channel, message).Err()
}

// SubscribeChannel subscribes to a Redis channel
func SubscribeChannel(channel string) *redis.PubSub {
	return RedisClient.Subscribe(ctx, channel)
}

// SetCache sets a key-value pair with expiration
func SetCache(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// GetCache gets a value by key
func GetCache(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// DeleteCache deletes a key
func DeleteCache(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

