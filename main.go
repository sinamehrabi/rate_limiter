package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis server address
		Password: "",               // Update with Redis server password if applicable
		DB:       0,                // Redis database index
	})
	// Create a rate limiter that allows 10 requests per second
	rateLimiter := NewRateLimiter(2, 50, time.Second)
	// Reverse proxy target URL
	targetURL, err := url.Parse("http://localhost:2222")
	if err != nil {
		log.Fatal(err)
	}
	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	// Define the handler function that includes rate limiting
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // In a production environment, use a more reliable way to obtain the client IP
		if rateLimiter.Allow(ip, rdb) {
			// Pass the request to the reverse proxy if allowed
			proxy.ServeHTTP(w, r)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	})
	// Start the HTTP server
	server := &http.Server{
		Addr:    ":8002",
		Handler: handler,
	}
	fmt.Println("Reverse proxy with rate limiter using Redis is running on :8002...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
