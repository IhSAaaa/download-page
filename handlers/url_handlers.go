package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"url-shortener/database"
	"url-shortener/models"
)

// CreateShortURL creates a new shortened URL
func CreateShortURL(c *gin.Context) {
	var req models.CreateURLRequest
	
	// Bind the raw JSON first to handle empty strings properly
	var rawData map[string]interface{}
	if err := c.ShouldBindJSON(&rawData); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}
	
	// Extract and validate original URL
	if originalURL, ok := rawData["original_url"].(string); ok && originalURL != "" {
		// Basic URL validation
		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "URL must start with http:// or https://",
			})
			return
		}
		req.OriginalURL = originalURL
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Original URL is required",
		})
		return
	}
	
	// Extract optional fields
	if customCode, ok := rawData["custom_code"].(string); ok && customCode != "" {
		req.CustomCode = &customCode
	}
	
	if title, ok := rawData["title"].(string); ok && title != "" {
		req.Title = &title
	}
	
	if description, ok := rawData["description"].(string); ok && description != "" {
		req.Description = &description
	}
	
	// Handle expires_at field
	if expiresAtStr, ok := rawData["expires_at"].(string); ok && expiresAtStr != "" {
		expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid expires_at format",
				"details": "Expected ISO 8601 format (e.g., 2024-01-01T12:00:00Z)",
			})
			return
		}
		req.ExpiresAt = &expiresAt
	}

	// Validate custom code if provided
	if req.CustomCode != nil {
		customCode := *req.CustomCode
		if len(customCode) < 3 || len(customCode) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Custom code must be between 3 and 50 characters",
			})
			return
		}

		// Validate custom code format (alphanumeric and hyphens only)
		for _, char := range customCode {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Custom code can only contain letters, numbers, and hyphens",
				})
				return
			}
		}

		// Check if custom code already exists
		var exists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE custom_code = $1)", customCode).Scan(&exists)
		if err != nil {
			log.Printf("Database error checking custom code: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Custom code already exists"})
			return
		}
	}

	// Create new URL
	url := models.NewURL(req.OriginalURL, req.CustomCode)
	url.Title = req.Title
	url.Description = req.Description
	url.ExpiresAt = req.ExpiresAt

	// Insert into database
	_, err := database.DB.Exec(`
		INSERT INTO urls (id, original_url, short_code, custom_code, title, description, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, url.ID, url.OriginalURL, url.ShortCode, url.CustomCode, url.Title, url.Description, url.ExpiresAt, url.CreatedAt, url.UpdatedAt)

	if err != nil {
		log.Printf("Database error creating URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create URL"})
		return
	}

	// Generate QR code
	baseURL := getBaseURL(c)
	shortURL := baseURL + "/" + url.ShortCode
	if url.CustomCode != nil {
		shortURL = baseURL + "/" + *url.CustomCode
	}

	qrCode, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
	if err != nil {
		// QR code generation failed, but URL was created successfully
		log.Printf("QR code generation failed: %v", err)
		qrCode = nil
	}

	response := url.ToResponse(baseURL)
	if qrCode != nil {
		response.QRCode = fmt.Sprintf("data:image/png;base64,%s", qrCode)
	}

	log.Printf("URL shortened successfully: %s -> %s", url.OriginalURL, shortURL)

	c.JSON(http.StatusCreated, gin.H{
		"message": "URL shortened successfully",
		"data":    response,
	})
}

// RedirectToOriginal redirects short URL to original URL
func RedirectToOriginal(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
		return
	}

	// Get URL from database
	var url models.URL
	err := database.DB.QueryRow(`
		SELECT id, original_url, short_code, custom_code, is_active, expires_at, click_count
		FROM urls 
		WHERE (short_code = $1 OR custom_code = $1) AND is_active = true
	`, shortCode).Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.CustomCode, &url.IsActive, &url.ExpiresAt, &url.ClickCount)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if URL is expired
	if url.IsExpired() {
		c.JSON(http.StatusGone, gin.H{"error": "URL has expired"})
		return
	}

	// Record click
	go recordClick(url.ID, c)

	// Increment click count
	_, err = database.DB.Exec("UPDATE urls SET click_count = click_count + 1, updated_at = $1 WHERE id = $2", time.Now(), url.ID)
	if err != nil {
		// Log error but don't fail the redirect
		fmt.Printf("Failed to update click count: %v\n", err)
	}

	// Redirect to original URL
	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

// GetAllURLs gets all URLs with pagination
func GetAllURLs(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	limit := getIntQuery(c, "limit", 10)
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM urls").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get URLs
	rows, err := database.DB.Query(`
		SELECT id, original_url, short_code, custom_code, title, description, is_active, expires_at, click_count, created_at, updated_at
		FROM urls 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var urls []models.URLResponse
	baseURL := getBaseURL(c)

	for rows.Next() {
		var url models.URL
		err := rows.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.CustomCode, &url.Title, &url.Description, &url.IsActive, &url.ExpiresAt, &url.ClickCount, &url.CreatedAt, &url.UpdatedAt)
		if err != nil {
			continue
		}
		urls = append(urls, url.ToResponse(baseURL))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": urls,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// GetURLByID gets a specific URL by ID
func GetURLByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL ID is required"})
		return
	}

	var url models.URL
	err := database.DB.QueryRow(`
		SELECT id, original_url, short_code, custom_code, title, description, is_active, expires_at, click_count, created_at, updated_at
		FROM urls WHERE id = $1
	`, id).Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.CustomCode, &url.Title, &url.Description, &url.IsActive, &url.ExpiresAt, &url.ClickCount, &url.CreatedAt, &url.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	baseURL := getBaseURL(c)
	response := url.ToResponse(baseURL)

	// Generate QR code
	shortURL := baseURL + "/" + url.ShortCode
	if url.CustomCode != nil {
		shortURL = baseURL + "/" + *url.CustomCode
	}

	qrCode, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
	if err == nil {
		response.QRCode = fmt.Sprintf("data:image/png;base64,%s", qrCode)
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// DeleteURL deletes a URL
func DeleteURL(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL ID is required"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM urls WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}

// GetURLAnalytics gets analytics for a specific URL
func GetURLAnalytics(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL ID is required"})
		return
	}

	// Get basic URL info
	var url models.URL
	err := database.DB.QueryRow("SELECT id, click_count FROM urls WHERE id = $1", id).Scan(&url.ID, &url.ClickCount)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get unique clicks count
	var uniqueClicks int64
	err = database.DB.QueryRow("SELECT COUNT(DISTINCT ip_address) FROM clicks WHERE url_id = $1", id).Scan(&uniqueClicks)
	if err != nil {
		uniqueClicks = 0
	}

	// Get top countries
	rows, err := database.DB.Query(`
		SELECT country, COUNT(*) as clicks 
		FROM clicks 
		WHERE url_id = $1 AND country IS NOT NULL 
		GROUP BY country 
		ORDER BY clicks DESC 
		LIMIT 5
	`, id)
	if err != nil {
		rows = nil
	}

	var topCountries []models.Country
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var country models.Country
			rows.Scan(&country.Country, &country.Clicks)
			topCountries = append(topCountries, country)
		}
	}

	// Get top devices
	rows, err = database.DB.Query(`
		SELECT device, COUNT(*) as clicks 
		FROM clicks 
		WHERE url_id = $1 AND device IS NOT NULL 
		GROUP BY device 
		ORDER BY clicks DESC 
		LIMIT 5
	`, id)
	if err != nil {
		rows = nil
	}

	var topDevices []models.Device
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var device models.Device
			rows.Scan(&device.Device, &device.Clicks)
			topDevices = append(topDevices, device)
		}
	}

	// Get top browsers
	rows, err = database.DB.Query(`
		SELECT browser, COUNT(*) as clicks 
		FROM clicks 
		WHERE url_id = $1 AND browser IS NOT NULL 
		GROUP BY browser 
		ORDER BY clicks DESC 
		LIMIT 5
	`, id)
	if err != nil {
		rows = nil
	}

	var topBrowsers []models.Browser
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var browser models.Browser
			rows.Scan(&browser.Browser, &browser.Clicks)
			topBrowsers = append(topBrowsers, browser)
		}
	}

	// Get click timeline (last 30 days)
	rows, err = database.DB.Query(`
		SELECT DATE(clicked_at) as date, COUNT(*) as clicks 
		FROM clicks 
		WHERE url_id = $1 AND clicked_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(clicked_at) 
		ORDER BY date DESC
	`, id)
	if err != nil {
		rows = nil
	}

	var timeline []models.Timeline
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var t models.Timeline
			rows.Scan(&t.Date, &t.Clicks)
			timeline = append(timeline, t)
		}
	}

	// Get last clicked at
	var lastClickedAt *time.Time
	err = database.DB.QueryRow("SELECT MAX(clicked_at) FROM clicks WHERE url_id = $1", id).Scan(&lastClickedAt)
	if err != nil {
		lastClickedAt = nil
	}

	analytics := models.Analytics{
		URLID:         url.ID,
		TotalClicks:   url.ClickCount,
		UniqueClicks:  uniqueClicks,
		TopCountries:  topCountries,
		TopDevices:    topDevices,
		TopBrowsers:   topBrowsers,
		ClickTimeline: timeline,
		LastClickedAt: lastClickedAt,
	}

	c.JSON(http.StatusOK, gin.H{"data": analytics})
}

// GetAllAnalytics gets analytics for all URLs
func GetAllAnalytics(c *gin.Context) {
	// Get total clicks
	var totalClicks int64
	err := database.DB.QueryRow("SELECT SUM(click_count) FROM urls").Scan(&totalClicks)
	if err != nil {
		totalClicks = 0
	}

	// Get total URLs
	var totalURLs int64
	err = database.DB.QueryRow("SELECT COUNT(*) FROM urls").Scan(&totalURLs)
	if err != nil {
		totalURLs = 0
	}

	// Get top performing URLs
	rows, err := database.DB.Query(`
		SELECT id, original_url, short_code, custom_code, click_count 
		FROM urls 
		ORDER BY click_count DESC 
		LIMIT 10
	`)
	if err != nil {
		rows = nil
	}

	var topURLs []gin.H
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var id, originalURL, shortCode string
			var customCode *string
			var clickCount int64
			rows.Scan(&id, &originalURL, &shortCode, &customCode, &clickCount)
			
			displayCode := shortCode
			if customCode != nil {
				displayCode = *customCode
			}

			topURLs = append(topURLs, gin.H{
				"id":          id,
				"original_url": originalURL,
				"short_code":   displayCode,
				"click_count":  clickCount,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_clicks": totalClicks,
			"total_urls":   totalURLs,
			"top_urls":     topURLs,
		},
	})
}

// Helper functions

func recordClick(urlID string, c *gin.Context) {
	click := models.Click{
		ID:        uuid.New().String(),
		URLID:     urlID,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Referer:   getStringPtr(c.GetHeader("Referer")),
		ClickedAt: time.Now(),
	}

	// TODO: Add geolocation and device detection
	// For now, we'll store basic information

	_, err := database.DB.Exec(`
		INSERT INTO clicks (id, url_id, ip_address, user_agent, referer, clicked_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, click.ID, click.URLID, click.IPAddress, click.UserAgent, click.Referer, click.ClickedAt)

	if err != nil {
		fmt.Printf("Failed to record click: %v\n", err)
	}
}

func getBaseURL(c *gin.Context) string {
	// Check if APP_URL is set in environment
	if appURL := os.Getenv("APP_URL"); appURL != "" {
		return appURL
	}
	
	// Fallback to request host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

func getIntQuery(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultValue
}

func getStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
} 