package models

import (
	"time"
	"github.com/google/uuid"
)

// URL represents a shortened URL
type URL struct {
	ID          string    `json:"id" db:"id"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	CustomCode  *string   `json:"custom_code,omitempty" db:"custom_code"`
	Title       *string   `json:"title,omitempty" db:"title"`
	Description *string   `json:"description,omitempty" db:"description"`
	UserID      *string   `json:"user_id,omitempty" db:"user_id"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	ClickCount  int64     `json:"click_count" db:"click_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateURLRequest represents the request to create a new short URL
type CreateURLRequest struct {
	OriginalURL string     `json:"original_url" binding:"required,url"`
	CustomCode  *string    `json:"custom_code,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// URLResponse represents the response for URL operations
type URLResponse struct {
	ID          string     `json:"id"`
	OriginalURL string     `json:"original_url"`
	ShortURL    string     `json:"short_url"`
	CustomCode  *string    `json:"custom_code,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	QRCode      string     `json:"qr_code,omitempty"`
	ClickCount  int64      `json:"click_count"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Click represents a click on a shortened URL
type Click struct {
	ID        string    `json:"id" db:"id"`
	URLID     string    `json:"url_id" db:"url_id"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Referer   *string   `json:"referer,omitempty" db:"referer"`
	Country   *string   `json:"country,omitempty" db:"country"`
	City      *string   `json:"city,omitempty" db:"city"`
	Device    *string   `json:"device,omitempty" db:"device"`
	Browser   *string   `json:"browser,omitempty" db:"browser"`
	OS        *string   `json:"os,omitempty" db:"os"`
	ClickedAt time.Time `json:"clicked_at" db:"clicked_at"`
}

// Analytics represents analytics data for a URL
type Analytics struct {
	URLID           string    `json:"url_id"`
	TotalClicks     int64     `json:"total_clicks"`
	UniqueClicks    int64     `json:"unique_clicks"`
	TopCountries    []Country `json:"top_countries"`
	TopDevices      []Device  `json:"top_devices"`
	TopBrowsers     []Browser `json:"top_browsers"`
	ClickTimeline   []Timeline `json:"click_timeline"`
	LastClickedAt   *time.Time `json:"last_clicked_at"`
}

// Country represents country analytics
type Country struct {
	Country string `json:"country"`
	Clicks  int64  `json:"clicks"`
}

// Device represents device analytics
type Device struct {
	Device string `json:"device"`
	Clicks int64  `json:"clicks"`
}

// Browser represents browser analytics
type Browser struct {
	Browser string `json:"browser"`
	Clicks  int64  `json:"clicks"`
}

// Timeline represents click timeline
type Timeline struct {
	Date  string `json:"date"`
	Clicks int64  `json:"clicks"`
}

// NewURL creates a new URL instance
func NewURL(originalURL string, customCode *string) *URL {
	return &URL{
		ID:          uuid.New().String(),
		OriginalURL: originalURL,
		ShortCode:   generateShortCode(),
		CustomCode:  customCode,
		IsActive:    true,
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// generateShortCode generates a random 6-character short code
func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// IsExpired checks if the URL has expired
func (u *URL) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}

// IncrementClickCount increments the click count
func (u *URL) IncrementClickCount() {
	u.ClickCount++
	u.UpdatedAt = time.Now()
}

// ToResponse converts URL to URLResponse
func (u *URL) ToResponse(baseURL string) URLResponse {
	shortCode := u.ShortCode
	if u.CustomCode != nil {
		shortCode = *u.CustomCode
	}

	return URLResponse{
		ID:          u.ID,
		OriginalURL: u.OriginalURL,
		ShortURL:    baseURL + "/" + shortCode,
		CustomCode:  u.CustomCode,
		Title:       u.Title,
		Description: u.Description,
		ClickCount:  u.ClickCount,
		ExpiresAt:   u.ExpiresAt,
		CreatedAt:   u.CreatedAt,
	}
} 