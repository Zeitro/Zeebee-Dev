package zeitro

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"zeebee-backend/common/errorx"
)

// Config represents the configuration for Zeitro API integration
type Config struct {
	URL      string
	Username string
	Password string
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.URL == "" {
		return errorx.New("Zeitro URL is required")
	}
	if _, err := url.Parse(c.URL); err != nil {
		return errorx.Wrap(err, "Zeitro URL is invalid")
	}
	if c.Username == "" {
		return errorx.New("Zeitro username is required")
	}
	if c.Password == "" {
		return errorx.New("Zeitro password is required")
	}
	// if c.MaxRetries <= 0 {
	// 	c.MaxRetries = 3 // default 3 retries
	// }
	return nil
}

// TokenResponse represents the token response from Zeitro API
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Token represents the authentication token for Zeitro API
type Token struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

// IsExpired checks if the token is expired or will expire soon
func (t *Token) IsExpired() bool {
	if t.ExpiresAt.IsZero() {
		return true
	}
	// Consider token expired if less than 10% of lifetime remains
	threshold := time.Duration(t.ExpiresIn/10) * time.Second
	return time.Until(t.ExpiresAt) < threshold
}

// IsValid checks if the token is valid and not expired
func (t *Token) IsValid() bool {
	return t.AccessToken != "" && !t.ExpiresAt.IsZero() && !t.IsExpired()
}

// ZeitroProperty represents a property from Zeitro API response
type ZeitroProperty struct {
	FullAddress     string `json:"FullAddress"`
	State           string `json:"State"`
	County          string `json:"County"`
	City            string `json:"City"`
	ZipCode         string `json:"ZipCode"`
	PropertyType    string `json:"PropertyType"`
	Price           string `json:"Price"`
	IncomeIndicator string `json:"IncomeIndicator"`
	MinorityPercent string `json:"MinorityPercent"`
	TotalPopulation string `json:"TotalPopulation"`
	Hispanic        string `json:"Hispanic"`
	Black           string `json:"Black"`
	Asian           string `json:"Asian"`
	Pacific         string `json:"Pacific"`
	MsaCode         string `json:"MsaCode"`
	DaysOnMarket    string `json:"DaysOnMarket"`
	AgentName       string `json:"AgentName"`
	AgentPhone      string `json:"AgentPhone"`
	AgentEmail      string `json:"AgentEmail"`
	Distance        string `json:"Distance"`
}

// PropertyListing represents a property listing after conversion
type PropertyListing struct {
	FullAddress     string
	State           string
	County          string
	City            string
	ZipCode         string
	PropertyType    string
	Price           string
	IncomeIndicator string
	MinorityPercent string
	TotalPopulation string
	Hispanic        string
	Black           string
	Asian           string
	Pacific         string
	MsaCode         string
	DaysOnMarket    string
	AgentName       string
	AgentPhone      string
	AgentEmail      string
	Distance        string
}

// Client represents the Zeitro API client
type Client struct {
	config  *Config
	client  *http.Client // Reusable HTTP client injected from outside
	token   *Token
	tokenMu sync.RWMutex
	sf      singleflight.Group
}

// NewClient creates a new Zeitro client with injected HTTP client
func NewClient(cfg Config, httpClient *http.Client) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	if httpClient == nil {
		return nil, errorx.New("HTTP client is required")
	}

	return &Client{
		config: &cfg,
		client: httpClient,
	}, nil
}

// GetToken returns a valid access token, refreshing if necessary
func (c *Client) GetToken(ctx context.Context) (*Token, error) {
	// Fast path: read with read lock
	c.tokenMu.RLock()
	if c.token != nil && c.token.IsValid() {
		c.tokenMu.RUnlock()
		return c.token, nil
	}
	c.tokenMu.RUnlock()

	// Slow path: need to refresh token
	result, err, _ := c.sf.Do("fetch-token", func() (interface{}, error) {
		return c.fetchToken(ctx)
	})

	if err != nil {
		// If refresh fails, try to use old token if still valid
		c.tokenMu.RLock()
		if c.token != nil && !c.token.ExpiresAt.IsZero() {
			c.tokenMu.RUnlock()
			return c.token, nil
		}
		c.tokenMu.RUnlock()
		return nil, err
	}

	token, ok := result.(*Token)
	if !ok {
		return nil, errorx.New("Failed to get token")
	}

	// Update cached token
	c.tokenMu.Lock()
	c.token = token
	c.tokenMu.Unlock()

	return token, nil
}

// fetchToken fetches a new access token from Zeitro API
func (c *Client) fetchToken(ctx context.Context) (*Token, error) {
	// Build request
	reqURL := fmt.Sprintf("%s/auth/access-token", strings.TrimSuffix(c.config.URL, "/"))
	body := strings.NewReader("grant_type=client_credentials")

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, body)
	if err != nil {
		return nil, errorx.New("Failed to create token request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.config.Username, c.config.Password)

	// Execute request using injected HTTP client
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errorx.New("Failed to fetch token: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errorx.New(fmt.Sprintf("Token request failed with status %d: %s", resp.StatusCode, string(bodyBytes)))
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, errorx.New("Failed to decode token response: " + err.Error())
	}

	// Create token with expiration times
	now := time.Now()
	token := &Token{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		ExpiresIn:   tokenResp.ExpiresIn,
		CreatedAt:   now,
		ExpiresAt:   now.Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	return token, nil
}

// SearchRequest represents a search request
type SearchRequest struct {
	Lat  float64 `json:"Lat"`
	Lon  float64 `json:"Lon"`
	Dist float64 `json:"Dist"`
	Addr string  `json:"Addr"`
}

// SearchListings searches for property listings from Zeitro API
func (c *Client) SearchListings(ctx context.Context, req SearchRequest) ([]PropertyListing, error) {
	// Get valid token
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errorx.New("Failed to marshal request: " + err.Error())
	}

	// Build HTTP request
	reqURL := fmt.Sprintf("%s/los/openhousesearch", strings.TrimSuffix(c.config.URL, "/"))
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, errorx.New("Failed to create request: " + err.Error())
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	// Execute request with retry using injected HTTP client
	var resp *http.Response
	resp, err = c.client.Do(httpReq)

	if err != nil {
		return nil, errorx.New("Failed to search listings: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errorx.New(fmt.Sprintf("Search request failed with status %d: %s", resp.StatusCode, string(bodyBytes)))
	}

	// Parse response
	var zeitroResp []ZeitroProperty
	if err := json.NewDecoder(resp.Body).Decode(&zeitroResp); err != nil {
		return nil, errorx.New("Failed to decode response: " + err.Error())
	}

	// Convert to internal model
	return c.convertProperties(zeitroResp), nil
}

// convertProperties converts Zeitro properties to internal model
func (c *Client) convertProperties(properties []ZeitroProperty) []PropertyListing {
	result := make([]PropertyListing, 0, len(properties))
	for _, prop := range properties {
		entity := PropertyListing{
			FullAddress:     strings.TrimSpace(prop.FullAddress),
			State:           strings.TrimSpace(prop.State),
			County:          strings.TrimSpace(prop.County),
			City:            strings.TrimSpace(prop.City),
			ZipCode:         strings.TrimSpace(prop.ZipCode),
			PropertyType:    strings.TrimSpace(prop.PropertyType),
			Price:           strings.TrimSpace(prop.Price),
			IncomeIndicator: strings.TrimSpace(prop.IncomeIndicator),
			MinorityPercent: strings.TrimSpace(prop.MinorityPercent),
			TotalPopulation: strings.TrimSpace(prop.TotalPopulation),
			Hispanic:        strings.TrimSpace(prop.Hispanic),
			Black:           strings.TrimSpace(prop.Black),
			Asian:           strings.TrimSpace(prop.Asian),
			Pacific:         strings.TrimSpace(prop.Pacific),
			MsaCode:         strings.TrimSpace(prop.MsaCode),
			DaysOnMarket:    strings.TrimSpace(prop.DaysOnMarket),
			AgentName:       strings.TrimSpace(prop.AgentName),
			AgentPhone:      strings.TrimSpace(prop.AgentPhone),
			AgentEmail:      strings.TrimSpace(prop.AgentEmail),
			Distance:        strings.TrimSpace(prop.Distance),
		}
		result = append(result, entity)
	}
	return result
}

// formatPhone removes non-numeric characters from phone number
func formatPhone(phone string) string {
	var result strings.Builder
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			result.WriteRune(r)
		}
	}
	return result.String()
}
