package zeitro

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestConfigValidate tests the Config.Validate() method
func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				URL:      "https://api.zeitro.com",
				Username: "testuser",
				Password: "testpass",
			},
			wantErr: false,
		},
		{
			name: "empty URL",
			config: Config{
				URL:      "",
				Username: "testuser",
				Password: "testpass",
			},
			wantErr: true,
			errMsg:  "Zeitro URL is required",
		},
		{
			name: "invalid URL",
			config: Config{
				URL:      "://invalid-url",
				Username: "testuser",
				Password: "testpass",
			},
			wantErr: true,
			errMsg:  "Zeitro URL is invalid",
		},
		{
			name: "empty username",
			config: Config{
				URL:      "https://api.zeitro.com",
				Username: "",
				Password: "testpass",
			},
			wantErr: true,
			errMsg:  "Zeitro username is required",
		},
		{
			name: "empty password",
			config: Config{
				URL:      "https://api.zeitro.com",
				Username: "testuser",
				Password: "",
			},
			wantErr: true,
			errMsg:  "Zeitro password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Config.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

// TestTokenIsExpired tests the Token.IsExpired() method
func TestTokenIsExpired(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		token    Token
		expected bool
	}{
		{
			name: "valid token with plenty of time",
			token: Token{
				AccessToken: "test-token",
				ExpiresIn:   3600,
				CreatedAt:   now,
				ExpiresAt:   now.Add(3600 * time.Second),
			},
			expected: false,
		},
		{
			name: "expired token",
			token: Token{
				AccessToken: "test-token",
				ExpiresIn:   3600,
				CreatedAt:   now.Add(-3700 * time.Second),
				ExpiresAt:   now.Add(-100 * time.Second),
			},
			expected: true,
		},
		{
			name: "token in expiration window (< 10%)",
			token: Token{
				AccessToken: "test-token",
				ExpiresIn:   3600,
				CreatedAt:   now.Add(-3300 * time.Second),
				ExpiresAt:   now.Add(300 * time.Second),
			},
			expected: true,
		},
		{
			name:     "zero expiration time",
			token:    Token{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.token.IsExpired(); got != tt.expected {
				t.Errorf("Token.IsExpired() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestTokenIsValid tests the Token.IsValid() method
func TestTokenIsValid(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		token    Token
		expected bool
	}{
		{
			name: "valid token",
			token: Token{
				AccessToken: "test-token",
				ExpiresIn:   3600,
				CreatedAt:   now,
				ExpiresAt:   now.Add(3600 * time.Second),
			},
			expected: true,
		},
		{
			name: "empty access token",
			token: Token{
				AccessToken: "",
				ExpiresIn:   3600,
				CreatedAt:   now,
				ExpiresAt:   now.Add(3600 * time.Second),
			},
			expected: false,
		},
		{
			name: "expired token",
			token: Token{
				AccessToken: "test-token",
				ExpiresIn:   3600,
				CreatedAt:   now.Add(-3700 * time.Second),
				ExpiresAt:   now.Add(-100 * time.Second),
			},
			expected: false,
		},
		{
			name:     "zero expiration time",
			token:    Token{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.token.IsValid(); got != tt.expected {
				t.Errorf("Token.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestGetToken tests the GetToken() method with mock server
func TestGetToken(t *testing.T) {
	// Create mock server for token endpoint
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/auth/access-token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("unexpected method: %s", r.Method)
		}

		// Check basic auth
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("missing basic auth")
		}
		if username != "testuser" || password != "testpass" {
			t.Errorf("invalid credentials: %s:%s", username, password)
		}

		// Return mock token response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "mock-access-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`))
	}))
	defer mockServer.Close()

	// Create client with mock server URL
	config := Config{
		URL:      mockServer.URL,
		Username: "testuser",
		Password: "testpass",
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}
	client, err := NewClient(config, httpClient)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// Test GetToken
	ctx := context.Background()
	token, err := client.GetToken(ctx)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	// Verify token
	if token.AccessToken != "mock-access-token" {
		t.Errorf("AccessToken = %v, want %v", token.AccessToken, "mock-access-token")
	}
	if token.TokenType != "Bearer" {
		t.Errorf("TokenType = %v, want %v", token.TokenType, "Bearer")
	}
	if token.ExpiresIn != 3600 {
		t.Errorf("ExpiresIn = %v, want %v", token.ExpiresIn, 3600)
	}
	if !token.IsValid() {
		t.Error("Token should be valid")
	}
}

// TestGetTokenCaching tests that GetToken() caches tokens
func TestGetTokenCaching(t *testing.T) {
	requestCount := 0

	// Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "mock-access-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`))
	}))
	defer mockServer.Close()

	config := Config{
		URL:      mockServer.URL,
		Username: "testuser",
		Password: "testpass",
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}
	client, err := NewClient(config, httpClient)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.Background()

	// First call should fetch token
	_, err = client.GetToken(ctx)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	// Second call should use cached token
	_, err = client.GetToken(ctx)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	// Verify only one request was made (caching worked)
	if requestCount != 1 {
		t.Errorf("requestCount = %v, want 1 (token should be cached)", requestCount)
	}
}

// TestNewClientWithError tests NewClient with invalid configuration
func TestNewClientWithError(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		httpClient  *http.Client
		expectedErr string
	}{
		{
			name: "invalid config - empty URL",
			config: Config{
				URL:      "",
				Username: "testuser",
				Password: "testpass",
			},
			httpClient:  &http.Client{},
			expectedErr: "Zeitro URL is required",
		},
		{
			name: "nil HTTP client",
			config: Config{
				URL:      "https://api.zeitro.com",
				Username: "testuser",
				Password: "testpass",
			},
			httpClient:  nil,
			expectedErr: "HTTP client is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.config, tt.httpClient)
			if err == nil {
				t.Fatal("NewClient() expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("NewClient() error = %v, want error containing %v", err, tt.expectedErr)
			}
		})
	}
}

// TestSearchListings tests the SearchListings() method with mock server
func TestSearchListings(t *testing.T) {
	tokenReceived := false

	// Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/access-token":
			// Token endpoint
			tokenReceived = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"access_token": "mock-access-token",
				"token_type": "Bearer",
				"expires_in": 3600
			}`))
		case "/los/openhousesearch":
			// Search endpoint
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer mock-access-token" {
				t.Errorf("Authorization header = %v, want %v", authHeader, "Bearer mock-access-token")
			}

			// Verify request body
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type = %v, want application/json", r.Header.Get("Content-Type"))
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"FullAddress": "123 Main St",
					"State": "CA",
					"County": "San Francisco",
					"City": "San Francisco",
					"ZipCode": "94102",
					"PropertyType": "House",
					"Price": "1000000",
					"IncomeIndicator": "High",
					"MinorityPercent": "25",
					"TotalPopulation": "50000",
					"Hispanic": "15",
					"Black": "5",
					"Asian": "30",
					"Pacific": "2",
					"MsaCode": "41860",
					"DaysOnMarket": "10",
					"AgentName": "John Doe",
					"AgentPhone": "(415) 555-1234",
					"AgentEmail": "john@example.com",
					"Distance": "5.2"
				}
			]`))
		}
	}))
	defer mockServer.Close()

	config := Config{
		URL:      mockServer.URL,
		Username: "testuser",
		Password: "testpass",
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}
	client, err := NewClient(config, httpClient)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.Background()
	req := SearchRequest{
		Lat:  37.7749,
		Lon:  -122.4194,
		Dist: 10.0,
		Addr: "San Francisco, CA",
	}

	// Execute search
	listings, err := client.SearchListings(ctx, req)
	if err != nil {
		t.Fatalf("SearchListings() error = %v", err)
	}

	// Verify token was fetched
	if !tokenReceived {
		t.Error("token was not fetched")
	}

	// Verify listings
	if len(listings) != 1 {
		t.Fatalf("len(listings) = %v, want 1", len(listings))
	}

	listing := listings[0]
	if listing.FullAddress != "123 Main St" {
		t.Errorf("FullAddress = %v, want %v", listing.FullAddress, "123 Main St")
	}
	if listing.State != "CA" {
		t.Errorf("State = %v, want %v", listing.State, "CA")
	}
	if listing.City != "San Francisco" {
		t.Errorf("City = %v, want %v", listing.City, "San Francisco")
	}
	// Verify phone number (keeps original format)
	if listing.AgentPhone != "(415) 555-1234" {
		t.Errorf("AgentPhone = %v, want %v", listing.AgentPhone, "(415) 555-1234")
	}
}

// TestSearchListingsWithError tests error handling in SearchListings
func TestSearchListingsWithError(t *testing.T) {
	// Mock server that returns error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/access-token":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"access_token": "mock-token",
				"token_type": "Bearer",
				"expires_in": 3600
			}`))
		case "/los/openhousesearch":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`Internal Server Error`))
		}
	}))
	defer mockServer.Close()

	config := Config{
		URL:      mockServer.URL,
		Username: "testuser",
		Password: "testpass",
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}
	client, err := NewClient(config, httpClient)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.Background()
	req := SearchRequest{
		Lat:  37.7749,
		Lon:  -122.4194,
		Dist: 10.0,
	}

	_, err = client.SearchListings(ctx, req)
	if err == nil {
		t.Error("SearchListings() expected error, got nil")
	}
}

// TestFormatPhone tests the formatPhone helper function
func TestFormatPhone(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "phone with dashes",
			input:    "(415) 555-1234",
			expected: "4155551234",
		},
		{
			name:     "phone with spaces",
			input:    "415 555 1234",
			expected: "4155551234",
		},
		{
			name:     "phone with dots",
			input:    "415.555.1234",
			expected: "4155551234",
		},
		{
			name:     "clean phone",
			input:    "4155551234",
			expected: "4155551234",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "phone with extensions",
			input:    "415-555-1234 x567",
			expected: "4155551234567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPhone(tt.input)
			if result != tt.expected {
				t.Errorf("formatPhone(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
