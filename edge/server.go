package edge

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sinamehrabi/rate_limiter/common"
	"github.com/sinamehrabi/rate_limiter/common/data"
	"github.com/sinamehrabi/rate_limiter/edge/middleware"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func StartServer(ctx context.Context, config *common.Config) {

	// Create a rate limiter that allows 10 requests per second
	rateLimiter := middleware.NewRateLimiter(2, 50, time.Minute)

	// Get redis client from context
	rdb, ready := ctx.Value(data.RedisClientKey).(*redis.Client)
	if !ready {
		log.Fatal("Redis client not found! Please try again!")
	}

	// Reverse proxy target URL
	targetHttpConfig := config.Server.TargetHttp
	targetURL, err := url.Parse(fmt.Sprintf("%s:%s", targetHttpConfig.Host, targetHttpConfig.Port))
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
			log.Printf("Incoming request from %s: %s %s | Reverse to -> %s:%s", ip, r.Method, r.URL.Path, targetHttpConfig.Host, targetHttpConfig.Port)
			proxy.ServeHTTP(w, r)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	})

	// Start the HTTP server
	httpConfig := config.Server.HTTP
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", httpConfig.Host, httpConfig.Port),
		Handler: handler,
	}

	fmt.Printf("Reverse proxy with rate limiter using Redis is running on %s:%s ...\n", httpConfig.Host, httpConfig.Port)

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
