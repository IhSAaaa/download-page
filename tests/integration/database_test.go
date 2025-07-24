package integration

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

// TestMain sets up the test database
func TestMain(m *testing.M) {
	// Setup test database connection
	var err error
	testDB, err = sql.Open("postgres", "postgres://postgres:password@localhost:5433/url_shortener_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	// Wait for database to be ready
	for i := 0; i < 30; i++ {
		err = testDB.Ping()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Exit(code)
}

// TestDatabaseConnection tests the database connection
func TestDatabaseConnection(t *testing.T) {
	err := testDB.Ping()
	require.NoError(t, err, "Database connection should be established")
}

// TestCreateTables tests table creation
func TestCreateTables(t *testing.T) {
	// Create URLs table
	createURLsTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		original_url TEXT NOT NULL,
		short_code VARCHAR(50) UNIQUE NOT NULL,
		user_id UUID,
		expires_at TIMESTAMP,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := testDB.Exec(createURLsTable)
	require.NoError(t, err, "Should create URLs table successfully")

	// Create clicks table
	createClicksTable := `
	CREATE TABLE IF NOT EXISTS clicks (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		url_id UUID REFERENCES urls(id) ON DELETE CASCADE,
		ip_address INET,
		user_agent TEXT,
		referer TEXT,
		country VARCHAR(2),
		city VARCHAR(100),
		clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = testDB.Exec(createClicksTable)
	require.NoError(t, err, "Should create clicks table successfully")
}

// TestURLOperations tests URL CRUD operations
func TestURLOperations(t *testing.T) {
	// Clean up before test
	_, err := testDB.Exec("DELETE FROM clicks")
	require.NoError(t, err)
	_, err = testDB.Exec("DELETE FROM urls")
	require.NoError(t, err)

	// Test 1: Insert URL
	t.Run("Insert URL", func(t *testing.T) {
		insertQuery := `
		INSERT INTO urls (original_url, short_code, is_active)
		VALUES ($1, $2, $3)
		RETURNING id`

		var id string
		err := testDB.QueryRow(insertQuery, "https://www.google.com", "test123", true).Scan(&id)
		require.NoError(t, err, "Should insert URL successfully")
		assert.NotEmpty(t, id, "Should return a valid UUID")
	})

	// Test 2: Query URL
	t.Run("Query URL", func(t *testing.T) {
		query := `
		SELECT original_url, short_code, is_active
		FROM urls
		WHERE short_code = $1`

		var originalURL, shortCode string
		var isActive bool
		err := testDB.QueryRow(query, "test123").Scan(&originalURL, &shortCode, &isActive)
		require.NoError(t, err, "Should query URL successfully")
		
		assert.Equal(t, "https://www.google.com", originalURL)
		assert.Equal(t, "test123", shortCode)
		assert.True(t, isActive)
	})

	// Test 3: Update URL
	t.Run("Update URL", func(t *testing.T) {
		updateQuery := `
		UPDATE urls
		SET is_active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE short_code = $2`

		result, err := testDB.Exec(updateQuery, false, "test123")
		require.NoError(t, err, "Should update URL successfully")
		
		rowsAffected, err := result.RowsAffected()
		require.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected, "Should update exactly one row")
	})

	// Test 4: Delete URL
	t.Run("Delete URL", func(t *testing.T) {
		deleteQuery := `DELETE FROM urls WHERE short_code = $1`
		
		result, err := testDB.Exec(deleteQuery, "test123")
		require.NoError(t, err, "Should delete URL successfully")
		
		rowsAffected, err := result.RowsAffected()
		require.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected, "Should delete exactly one row")
	})
}

// TestClickOperations tests click tracking operations
func TestClickOperations(t *testing.T) {
	// Clean up before test
	_, err := testDB.Exec("DELETE FROM clicks")
	require.NoError(t, err)
	_, err = testDB.Exec("DELETE FROM urls")
	require.NoError(t, err)

	// Insert test URL first
	var urlID string
	err = testDB.QueryRow(`
		INSERT INTO urls (original_url, short_code, is_active)
		VALUES ($1, $2, $3)
		RETURNING id`, "https://www.google.com", "test123", true).Scan(&urlID)
	require.NoError(t, err)

	// Test 1: Insert click
	t.Run("Insert Click", func(t *testing.T) {
		insertQuery := `
		INSERT INTO clicks (url_id, ip_address, user_agent, referer, country, city)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

		var clickID string
		err := testDB.QueryRow(insertQuery, 
			urlID, 
			"192.168.1.1", 
			"Mozilla/5.0", 
			"https://example.com", 
			"US", 
			"New York").Scan(&clickID)
		
		require.NoError(t, err, "Should insert click successfully")
		assert.NotEmpty(t, clickID, "Should return a valid UUID")
	})

	// Test 2: Query clicks count
	t.Run("Query Clicks Count", func(t *testing.T) {
		query := `SELECT COUNT(*) FROM clicks WHERE url_id = $1`
		
		var count int
		err := testDB.QueryRow(query, urlID).Scan(&count)
		require.NoError(t, err, "Should query clicks count successfully")
		assert.Equal(t, 1, count, "Should have exactly one click")
	})

	// Test 3: Query clicks with analytics
	t.Run("Query Clicks Analytics", func(t *testing.T) {
		query := `
		SELECT 
			COUNT(*) as total_clicks,
			COUNT(DISTINCT ip_address) as unique_clicks,
			COUNT(DISTINCT country) as countries
		FROM clicks 
		WHERE url_id = $1`
		
		var totalClicks, uniqueClicks, countries int
		err := testDB.QueryRow(query, urlID).Scan(&totalClicks, &uniqueClicks, &countries)
		require.NoError(t, err, "Should query analytics successfully")
		
		assert.Equal(t, 1, totalClicks, "Should have 1 total click")
		assert.Equal(t, 1, uniqueClicks, "Should have 1 unique click")
		assert.Equal(t, 1, countries, "Should have 1 country")
	})
}

// TestDatabaseConstraints tests database constraints
func TestDatabaseConstraints(t *testing.T) {
	// Clean up before test
	_, err := testDB.Exec("DELETE FROM clicks")
	require.NoError(t, err)
	_, err = testDB.Exec("DELETE FROM urls")
	require.NoError(t, err)

	// Test 1: Unique constraint on short_code
	t.Run("Unique Short Code Constraint", func(t *testing.T) {
		// Insert first URL
		_, err := testDB.Exec(`
			INSERT INTO urls (original_url, short_code, is_active)
			VALUES ($1, $2, $3)`, "https://www.google.com", "test123", true)
		require.NoError(t, err)

		// Try to insert second URL with same short_code
		_, err = testDB.Exec(`
			INSERT INTO urls (original_url, short_code, is_active)
			VALUES ($1, $2, $3)`, "https://www.github.com", "test123", true)
		
		assert.Error(t, err, "Should fail due to unique constraint violation")
	})

	// Test 2: Foreign key constraint
	t.Run("Foreign Key Constraint", func(t *testing.T) {
		// Try to insert click with non-existent URL ID
		_, err := testDB.Exec(`
			INSERT INTO clicks (url_id, ip_address)
			VALUES ($1, $2)`, "00000000-0000-0000-0000-000000000000", "192.168.1.1")
		
		assert.Error(t, err, "Should fail due to foreign key constraint violation")
	})
} 