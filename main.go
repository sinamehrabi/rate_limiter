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
