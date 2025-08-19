package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type visitor struct {
	tokens     float64
	lastRefill time.Time
	lastSeen   time.Time
}
type RateLimiter struct {
	visitors    map[string]*visitor
	rate        int
	burst       int
	mu          sync.Mutex
	cleanupDone chan bool
}

func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors:    make(map[string]*visitor),
		rate:        rate,
		burst:       burst,
		cleanupDone: make(chan bool),
	}

	go rl.startCleanup()

	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			tokens:     float64(rl.burst),
			lastRefill: time.Now(),
			lastSeen:   time.Now(),
		}
		rl.visitors[ip] = v
	}

	now := time.Now()
	v.lastSeen = now
	elapsed := now.Sub(v.lastRefill)
	refillTokens := elapsed.Seconds() * float64(rl.rate)

	if refillTokens > 0 {
		v.tokens += refillTokens
		if v.tokens > float64(rl.burst) {
			v.tokens = float64(rl.burst)
		}
		v.lastRefill = now
	}

	if v.tokens >= 1 {
		v.tokens -= 1
		return true
	}
	return false
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		if !rl.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {

		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	ip := r.RemoteAddr
	if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
		ip = ip[:colonIndex]
	}

	return ip
}

func (rl *RateLimiter) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute) // Cleanup every 5 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanupVisitors()
		case <-rl.cleanupDone:
			return
		}
	}
}

func (rl *RateLimiter) cleanupVisitors() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	inactiveThreshold := 30 * time.Minute // Clean inactive visitors after 30 minutes

	for ip, visitor := range rl.visitors {
		if now.Sub(visitor.lastSeen) > inactiveThreshold {
			delete(rl.visitors, ip)
		}
	}
}

func (rl *RateLimiter) Stop() {
	close(rl.cleanupDone)
}
