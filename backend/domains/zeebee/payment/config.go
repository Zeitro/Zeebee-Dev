package payment

import (
	"net/url"

	"zeebee-backend/common/errorx"
)

// Config represents the configuration for Stripe payment integration.
type Config struct {
	SecretKey     string
	WebhookSecret string
	SuccessUrl    string
	CancelUrl     string
	Currency      string
}

// Enabled reports whether Stripe integration is configured for this service.
func (c *Config) Enabled() bool {
	return c != nil && c.SecretKey != ""
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.SecretKey == "" {
		return errorx.New("stripe secret key is required")
	}
	if c.SuccessUrl == "" {
		return errorx.New("stripe checkout success url is required")
	}
	if _, err := url.Parse(c.SuccessUrl); err != nil {
		return errorx.Wrap(err, "stripe checkout success url is invalid")
	}
	if c.CancelUrl == "" {
		return errorx.New("stripe checkout cancel url is required")
	}
	if _, err := url.Parse(c.CancelUrl); err != nil {
		return errorx.Wrap(err, "stripe checkout cancel url is invalid")
	}
	if c.Currency == "" {
		c.Currency = "usd"
	}
	return nil
}
