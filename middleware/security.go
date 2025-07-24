package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:;")
		
		// X-Frame-Options
		c.Header("X-Frame-Options", "DENY")
		
		// X-Content-Type-Options
		c.Header("X-Content-Type-Options", "nosniff")
		
		// X-XSS-Protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		// Strict Transport Security (only for HTTPS)
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		c.Next()
	}
}

// InputValidation validates and sanitizes input
func InputValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate Content-Type for POST/PUT requests
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// Allow specific origins or use * for development
		allowedOrigins := []string{
			"http://localhost:4000",
			"http://localhost:3000",
			"https://yourdomain.com", // Add your production domain
		}
		
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		
		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// RequestLogger logs incoming requests
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// URLValidation validates URL format
func URLValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware can be used to validate URL parameters
		// For example, validate short codes format
		shortCode := c.Param("shortCode")
		if shortCode != "" {
			// Validate short code format (alphanumeric, 3-50 characters)
			if len(shortCode) < 3 || len(shortCode) > 50 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid short code format",
				})
				c.Abort()
				return
			}
			
			// Check if short code contains only alphanumeric characters
			for _, char := range shortCode {
				if !((char >= 'a' && char <= 'z') || 
					 (char >= 'A' && char <= 'Z') || 
					 (char >= '0' && char <= '9') ||
					 char == '-' || char == '_') {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "Short code can only contain letters, numbers, hyphens, and underscores",
					})
					c.Abort()
					return
				}
			}
		}
		
		c.Next()
	}
} 