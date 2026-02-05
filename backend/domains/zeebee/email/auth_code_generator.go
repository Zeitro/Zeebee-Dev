package email

import (
	"fmt"

	"golang.org/x/oauth2"

	"zeebee-backend/domains/zeebee/constants"
)

type AuthCodeUrlGenerator interface {
	GenerateAuthCodeUrl(emailScope constants.EmailPurpose) (url string)
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

func (g *GmailAuthCodeUrlGenerator) GenerateAuthCodeUrl(emailScope constants.EmailPurpose) (url string) {
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
