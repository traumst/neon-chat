package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	h "neon-chat/src/utils/http"
)

func ThrottlingUserMiddleware(rate int, burst int) Middleware {
	limiter := newUserLimiter(rate, burst)

	return Middleware{
		Name: "ThrottlingUser",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				key := getClientIP(r)
				if key == "" {
					panic("rate limit user key is empty")
				}
				if !limiter.isAllowed(key) {
					log.Println("Too Many User Requests", key)
					h.SetRetryAfterHeader(&w, 5)
					http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(w, r)
			})
		}}
}

type userLimiter struct {
	mu sync.Mutex
	// normal RPS
	rps int
	// max burst RPS
	burst int
	// user token bucket
	tokens    map[string]int
	lastCheck map[string]time.Time
}

func newUserLimiter(rate int, burst int) *userLimiter {
	if rate <= 0 || burst <= 0 || rate > burst {
		panic("invalid rate or burst")
	}
	return &userLimiter{
		rps:       rate,
		burst:     burst,
		tokens:    make(map[string]int),
		lastCheck: make(map[string]time.Time),
	}
}

func (rl *userLimiter) isAllowed(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastCheck[key])
	rl.lastCheck[key] = now

	rl.tokens[key] += int(elapsed.Seconds()) * rl.rps
	if rl.tokens[key] > rl.burst {
		rl.tokens[key] = rl.burst
	}

	if rl.tokens[key] > 0 {
		rl.tokens[key]--
		return true
	}

	return false
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
