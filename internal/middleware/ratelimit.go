package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu       sync.Mutex
	counters map[string]*counter
	limit    int
}

type counter struct {
	count    int
	resetAt  time.Time
}

func RateLimit(perMinute int) func(http.Handler) http.Handler {
	rl := &rateLimiter{counters: make(map[string]*counter), limit: perMinute}
	go rl.cleanup()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if !rl.allow(ip) {
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	c, ok := rl.counters[ip]
	now := time.Now()
	if !ok || now.After(c.resetAt) {
		rl.counters[ip] = &counter{count: 1, resetAt: now.Add(time.Minute)}
		return true
	}
	if c.count >= rl.limit {
		return false
	}
	c.count++
	return true
}

func (rl *rateLimiter) cleanup() {
	for range time.Tick(5 * time.Minute) {
		rl.mu.Lock()
		now := time.Now()
		for ip, c := range rl.counters {
			if now.After(c.resetAt) {
				delete(rl.counters, ip)
			}
		}
		rl.mu.Unlock()
	}
}
