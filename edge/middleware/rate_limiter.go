package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	capacity     int
	refillAmount int
	refillPeriod time.Duration
}

func NewRateLimiter(capacity, refillAmount int, refillPeriod time.Duration) *RateLimiter {
	return &RateLimiter{
		capacity:     capacity,
		refillAmount: refillAmount,
		refillPeriod: refillPeriod,
	}
}
func (r *RateLimiter) Allow(ip string, rdb *redis.Client) bool {
	key := fmt.Sprintf("ratelimit:%s", ip)
	countCmd := rdb.Incr(rdb.Context(), key)
	if countCmd.Err() != nil {
		log.Println("Error storing rate limit in Redis:", countCmd.Err())
		return false
	}

	count, err := countCmd.Result()
	if err != nil {
		log.Println("Error retrieving rate limit from Redis:", err)
		return false
	}
	if count > int64(r.capacity) {
		log.Println(count)

		return false
	} else {
		expiredTime := time.Duration(r.refillAmount) * r.refillPeriod

		expireCmd := rdb.Expire(rdb.Context(), key, expiredTime)
		if expireCmd.Err() != nil {
			log.Println("Error setting expiration for rate limit key in Redis:", expireCmd.Err())
		}
	}

	return true
}
