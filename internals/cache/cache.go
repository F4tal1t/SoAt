package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var cache *redis.Client

func Client() *redis.Client {
	if cache == nil {
		fmt.Println("Warning: Redis client is not initialized")
	}
	return cache
}

func Connect() {
	redisAddr := os.Getenv("REDIS_URL")
	fmt.Printf("DEBUG: Attempting to connect to Redis with URL: '%s'\n", redisAddr)
	if redisAddr == "" {
		fmt.Println("REDIS_URL environment variable is not set, skipping Redis connection")
		return
	}

	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		fmt.Println("Failed to parse Redis URL, skipping Redis connection")
		fmt.Println(err)
		return
	}

	rdb := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("Failed to connect to Redis, continuing without cache")
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to Redis successfully")
	cache = rdb
}
