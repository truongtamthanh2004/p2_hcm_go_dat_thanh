package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

// var (
// 	Rdb *redis.Client
// 	Ctx = context.Background()
// )

// func InitRedis() {
// 	Rdb = redis.NewClient(&redis.Options{
// 		Addr:     os.Getenv("REDIS_ADDR"),
// 		Password: os.Getenv("REDIS_PASSWORD"),
// 		DB:       0,
// 	})

// 	// test connection
// 	_, err := Rdb.Ping(Ctx).Result()
// 	if err != nil {
// 		log.Fatalf("failed to connect to redis: %v", err)
// 	}
// 	log.Println("Connected to Redis")
// }

// // Set cache
// func Set(key string, value string, ttl time.Duration) error {
// 	return Rdb.Set(Ctx, key, value, ttl).Err()
// }

// // Get cache
// func Get(key string) (string, error) {
// 	return Rdb.Get(Ctx, key).Result()
// }
var Rdb *redis.Client

func InitRedis(ctx context.Context) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// test connection
	if _, err := Rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("✅ Connected to Redis")
}

// Set cache
func Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return Rdb.Set(ctx, key, value, ttl).Err()
}

// Get cache
func Get(ctx context.Context, key string) (string, error) {
	return Rdb.Get(ctx, key).Result()
}
