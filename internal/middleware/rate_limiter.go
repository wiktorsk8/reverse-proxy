package middleware

import (
	"net/http"
	"sync"

	"github.com/wiktorsk8/reverse-proxy/internal/config"
	"github.com/wiktorsk8/reverse-proxy/internal/tools"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     sync.RWMutex
	config config.RateLimit
}

func NewRateLimiterMiddleware(config config.RateLimit) *RateLimiter {
	return &RateLimiter{
		ips:    make(map[string]*rate.Limiter),
		config: config,
	}
}

func (r *RateLimiter) addIp(ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	limiter := rate.NewLimiter(rate.Limit(0.5), r.config.Burst)
	r.ips[ip] = limiter
}

func (r *RateLimiter) getLimiter(ip string) *rate.Limiter {
	for {
		r.mu.RLock()
		limiter, ok := r.ips[ip]
		r.mu.RUnlock()
		if ok {
			return limiter
		}
		r.addIp(ip)
	}
}

func (r *RateLimiter) GetMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ip := tools.GetIpFromRequest(req)
			limiter := r.getLimiter(ip)
			if !limiter.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
			}
			next.ServeHTTP(w, req)
		})
	}
}
