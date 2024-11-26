package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	h "neon-chat/src/utils/http"
)

func ThrottlingTotalMiddleware(rate int, burst int) Middleware {
	limiter := newRateLimiter(rate, burst)

	return Middleware{
		Name: "ThrottlingTotal",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !limiter.isAllowed() {
					log.Println("Too Many Total Requests")
					h.SetRetryAfterHeader(&w, 1)
					http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(w, r)
			})
		}}
}

type rateLimiter struct {
	mu sync.Mutex
	// normal RPS
	rps int
	// max burst RPS
	burst int
	// token bucket
	tokens    int
	lastCheck time.Time
}

func newRateLimiter(rate int, burst int) *rateLimiter {
	return &rateLimiter{
		rps:       rate,
		burst:     burst,
		tokens:    burst,
		lastCheck: time.Now(),
	}
}

func (rl *rateLimiter) isAllowed() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastCheck)
	rl.lastCheck = now

	rl.tokens += int(elapsed.Seconds()) * rl.rps
	if rl.tokens > rl.burst {
		rl.tokens = rl.burst
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}
