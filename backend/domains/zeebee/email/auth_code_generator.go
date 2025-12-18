package email

import (
	"fmt"

	"golang.org/x/oauth2"

	"zeebee-backend/domains/zeebee/constants"
)

const (
	_oauthOrganizationIdKey   = "organization_id"
	_oauthEmailScopeKey       = "email_scope"
	_oauthEmailKey            = "auth_email"
	_oauthOrganizationTypeKey = "organization_type"
)

type AuthCodeUrlGenerator interface {
	GenerateAuthCodeUrl(emailScope constants.EmailScope) (url string)
	GetConfig() *oauth2.Config
}

type GmailAuthCodeUrlGenerator struct {
	config *oauth2.Config
}

func NewGmailAuthCodeUrlGenerator(config *oauth2.Config) *GmailAuthCodeUrlGenerator {
	return &GmailAuthCodeUrlGenerator{
		config: config,
	}
}

func (g *GmailAuthCodeUrlGenerator) GenerateAuthCodeUrl(emailScope constants.EmailScope) (url string) {
	// randomState := random.RandString(10)
	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	}
	url = g.config.AuthCodeURL(fmt.Sprintf("%v", emailScope), opts...)
	return
}
func (g *GmailAuthCodeUrlGenerator) GetConfig() *oauth2.Config {
	return g.config
}
