package ratelimit_middleware

import (
	"go-upcycle_connect-backend/utils/response"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	visitors map[string]int
	mu       sync.Mutex
	limit    int
	interval time.Duration
}

func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]int),
		limit:    limit,
		interval: interval,
	}
}

func (rl *RateLimiter) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rl.mu.Lock()
		defer rl.mu.Unlock()

		if count, ok := rl.visitors[ip]; ok {
			if count >= rl.limit {
				response.NewErrorMessage(w, "Hop hop hop, on se calme", http.StatusTooManyRequests)
				return
			}
			rl.visitors[ip] = count + 1
		} else {
			rl.visitors[ip] = 1
			go func() {
				time.Sleep(rl.interval)
				rl.mu.Lock()
				defer rl.mu.Unlock()
				delete(rl.visitors, ip)
			}()
		}

		next.ServeHTTP(w, r)
	}
}
