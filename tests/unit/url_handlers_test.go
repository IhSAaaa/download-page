package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock database for testing
type MockDB struct {
	urls map[string]interface{}
}

func (m *MockDB) CreateURL(url interface{}) error {
	// Mock implementation
	return nil
}

func (m *MockDB) GetURLByShortCode(shortCode string) (interface{}, error) {
	// Mock implementation
	if url, exists := m.urls[shortCode]; exists {
		return url, nil
	}
	return nil, nil
}

func (m *MockDB) GetURLs() ([]interface{}, error) {
	// Mock implementation
	urls := make([]interface{}, 0)
	for _, url := range m.urls {
		urls = append(urls, url)
	}
	return urls, nil
}

func (m *MockDB) RecordClick(urlID string, clickData interface{}) error {
	// Mock implementation
	return nil
}

func (m *MockDB) GetAnalytics(urlID string) (interface{}, error) {
	// Mock implementation
	return map[string]interface{}{
		"total_clicks": 100,
		"unique_clicks": 50,
	}, nil
}

// TestCreateShortURL tests the URL creation endpoint
func TestCreateShortURL(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	
	// Mock database
	mockDB := &MockDB{
		urls: make(map[string]interface{}),
	}

	// Setup routes (you'll need to import your actual handlers)
	// router.POST("/api/shorten", handlers.CreateShortURL)

	// Test case 1: Valid URL creation
	t.Run("Valid URL Creation", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"original_url": "https://www.google.com",
			"custom_code":  "test123",
		}
		
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.Contains(t, response, "short_url")
		assert.Contains(t, response, "original_url")
	})

	// Test case 2: Invalid URL
	t.Run("Invalid URL", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"original_url": "invalid-url",
			"custom_code":  "test456",
		}
		
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case 3: Missing original URL
	t.Run("Missing Original URL", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"custom_code": "test789",
		}
		
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestRedirectToOriginal tests the URL redirection endpoint
func TestRedirectToOriginal(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	
	// Mock database with test data
	mockDB := &MockDB{
		urls: map[string]interface{}{
			"test123": map[string]interface{}{
				"original_url": "https://www.google.com",
				"is_active":    true,
			},
		},
	}

	// Setup routes (you'll need to import your actual handlers)
	// router.GET("/:shortCode", handlers.RedirectToOriginal)

	// Test case 1: Valid short code
	t.Run("Valid Short Code", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test123", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusMovedPermanently, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "google.com")
	})

	// Test case 2: Invalid short code
	t.Run("Invalid Short Code", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/invalid", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestGetURLs tests the URL listing endpoint
func TestGetURLs(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	
	// Mock database
	mockDB := &MockDB{
		urls: map[string]interface{}{
			"test1": map[string]interface{}{
				"original_url": "https://www.google.com",
				"short_code":   "test1",
			},
			"test2": map[string]interface{}{
				"original_url": "https://www.github.com",
				"short_code":   "test2",
			},
		},
	}

	// Setup routes (you'll need to import your actual handlers)
	// router.GET("/api/urls", handlers.GetURLs)

	// Test case: Get all URLs
	t.Run("Get All URLs", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/urls", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.Len(t, response, 2)
	})
}

// TestGetAnalytics tests the analytics endpoint
func TestGetAnalytics(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	
	// Mock database
	mockDB := &MockDB{}

	// Setup routes (you'll need to import your actual handlers)
	// router.GET("/api/analytics/:urlID", handlers.GetAnalytics)

	// Test case: Get analytics for valid URL
	t.Run("Get Analytics", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/analytics/test123", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.Contains(t, response, "total_clicks")
		assert.Contains(t, response, "unique_clicks")
	})
} 