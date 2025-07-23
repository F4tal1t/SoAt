package cache

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var cache *redis.Client

func Client() *redis.Client {
	return cache
}

func Connect() {
	redisAddr := os.Getenv("REDIS_URL")
	fmt.Printf("DEBUG: Attempting to connect to Redis with URL: '%s'\n", redisAddr)
	if redisAddr == "" {
		panic("REDIS_URL environment variable is not set")
	}

	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		fmt.Println("Failed to parse Redis URL")
		panic(err)
	}

	rdb := redis.NewClient(opt)
	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("Failed to connect to redis")
		panic(err)
	}

	fmt.Println("Connected to Redis successfully")
	cache = rdb
}
