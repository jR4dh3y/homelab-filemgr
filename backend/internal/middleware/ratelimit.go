// Package middleware provides HTTP middleware for the file manager API.
package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter provides per-IP rate limiting for HTTP requests.
// It uses a token bucket algorithm to allow bursts while enforcing
// a sustained rate limit.
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      float64 // requests per second
	burst    int     // max burst size
}

// NewRateLimiter creates a new rate limiter with the specified requests per second.
// The burst size is set to 2x the RPS to allow short bursts.
// If rps is 0 or negative, rate limiting is effectively disabled (very high limit).
func NewRateLimiter(rps float64) *RateLimiter {
	if rps <= 0 {
		rps = 1000 // Effectively disabled
	}
	burst := int(rps * 2)
	if burst < 1 {
		burst = 1
	}
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
		burst:    burst,
	}
}

// getLimiter returns the rate limiter for the given IP address,
// creating one if it doesn't exist.
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()

	if exists {
		return limiter
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check after acquiring write lock
	if limiter, exists = rl.limiters[ip]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
	rl.limiters[ip] = limiter
	return limiter
}

// Allow returns true if the request from the given IP should be allowed.
func (rl *RateLimiter) Allow(ip string) bool {
	return rl.getLimiter(ip).Allow()
}

// RateLimit returns a middleware that limits requests per IP address.
// Requests that exceed the rate limit receive a 429 Too Many Requests response.
func RateLimit(rps float64) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(rps)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			if !limiter.Allow(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"Too many requests","code":"rate_limit_exceeded"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP address from the request.
// It checks X-Forwarded-For and X-Real-IP headers first (for proxied requests),
// then falls back to RemoteAddr.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (may contain multiple IPs)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list (original client)
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr (strip port)
	addr := r.RemoteAddr
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i]
		}
	}
	return addr
}

// StartCleanup starts a background goroutine that periodically removes
// stale rate limiters to prevent memory growth.
func (rl *RateLimiter) StartCleanup(interval time.Duration, stopCh <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				rl.cleanup()
			}
		}
	}()
}

// cleanup removes all rate limiters. Since limiters are created lazily
// and are lightweight, this is a simple way to prevent unbounded growth.
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	// Create a new map instead of deleting entries (more efficient)
	rl.limiters = make(map[string]*rate.Limiter)
}
