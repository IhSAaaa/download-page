package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"url-shortener/database"
	"url-shortener/handlers"
	"url-shortener/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create tables
	if err := database.CreateTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := gin.New()

	// Use custom logger instead of gin.Default()
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())

	// Security middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.InputValidation())
	r.Use(middleware.CORS())

	// Rate limiting middleware (100 requests per minute per IP)
	r.Use(middleware.RateLimitMiddleware(100, time.Minute))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "URL Shortener Service is running",
			"version": "1.0.0",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// URL shortening endpoints
		api.POST("/shorten", handlers.CreateShortURL)
		api.GET("/urls", handlers.GetAllURLs)
		api.GET("/urls/:id", handlers.GetURLByID)
		api.DELETE("/urls/:id", handlers.DeleteURL)
		
		// Analytics endpoints
		api.GET("/analytics/:id", handlers.GetURLAnalytics)
		api.GET("/analytics", handlers.GetAllAnalytics)
	}

	// URL validation middleware for short code routes
	r.Use(middleware.URLValidation())

	// Serve static files for React frontend (before the short URL route)
	r.Static("/static", "./frontend/dist/static")
	r.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")
	r.StaticFile("/manifest.json", "./frontend/dist/manifest.json")
	
	// Redirect endpoint (for short URLs) - must be after static files
	r.GET("/:shortCode", handlers.RedirectToOriginal)
	
	// Fallback for React Router - serve index.html for all non-API routes
	r.NoRoute(func(c *gin.Context) {
		// Don't serve index.html for API routes
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}
		
		// Serve React app for all other routes
		c.File("./frontend/dist/index.html")
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 