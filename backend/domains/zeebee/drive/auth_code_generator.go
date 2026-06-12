package drive

import (
	"golang.org/x/oauth2"
)

// DriveAuthCodeUrlGenerator builds Google OAuth consent URLs for Drive authorization.
type DriveAuthCodeUrlGenerator struct {
	config *oauth2.Config
}

// NewDriveAuthCodeUrlGenerator creates a generator bound to the given OAuth config.
func NewDriveAuthCodeUrlGenerator(config *oauth2.Config) *DriveAuthCodeUrlGenerator {
	return &DriveAuthCodeUrlGenerator{config: config}
}

// GenerateAuthCodeUrl returns the consent URL for the given state.
// AccessTypeOffline + ApprovalForce guarantee a refresh token is always issued.
func (g *DriveAuthCodeUrlGenerator) GenerateAuthCodeUrl(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// GetConfig exposes the underlying OAuth config.
func (g *DriveAuthCodeUrlGenerator) GetConfig() *oauth2.Config {
	return g.config
}
