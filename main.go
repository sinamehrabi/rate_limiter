package main

import (
	"context"
	"flag"
	"github.com/go-redis/redis/v8"
	"github.com/sinamehrabi/rate_limiter/common"
	"github.com/sinamehrabi/rate_limiter/common/data"
	"github.com/sinamehrabi/rate_limiter/edge"
	"log"
)

func main() {

	ctx, _ := context.WithCancel(context.Background())
	configPath := flag.String("config", "", "path to config file")

	flag.Parse()

	if configPath == nil || *configPath == "" {
		flag.Usage()
		log.Fatal(ctx, "config file not specified")
	}

	config, err := common.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(ctx, "failed to load config")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port, // Update with your Redis server address
		Username: config.Redis.Username,                       // Update with Redis server username if applicable
		Password: config.Redis.Password,                       // Update with Redis server password if applicable
		DB:       config.Redis.DB,                             // Redis database index
	})

	_, err = rdb.Ping(ctx).Result()

	if err != nil {
		log.Fatal(ctx, err)
	}

	// Create a new context with the Redis client
	ctx = context.WithValue(ctx, data.RedisClientKey, rdb)

	edge.StartServer(ctx, config)
}
