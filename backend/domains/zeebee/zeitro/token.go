package zeitro

import "time"

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
