package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter represents a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimitMiddleware creates a middleware that limits requests per IP
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)
	
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()
		
		// Check if client has exceeded rate limit
		if !limiter.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
				"retry_after": int(window.Seconds()),
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	windowStart := now.Add(-rl.window)
	
	// Get existing requests for this IP
	requests, exists := rl.requests[ip]
	if !exists {
		requests = []time.Time{}
	}
	
	// Filter out old requests outside the window
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}
	
	// Check if limit is exceeded
	if len(validRequests) >= rl.limit {
		return false
	}
	
	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests
	
	return true
}

// Cleanup removes old entries to prevent memory leaks
func (rl *RateLimiter) Cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	windowStart := now.Add(-rl.window)
	
	for ip, requests := range rl.requests {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}
		
		if len(validRequests) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = validRequests
		}
	}
}

// StartCleanup starts a background goroutine to clean up old entries
func (rl *RateLimiter) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for range ticker.C {
			rl.Cleanup()
		}
	}()
} 