package global

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmailapi "google.golang.org/api/gmail/v1"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/email"
	"zeebee-backend/domains/zeebee/constants"
	zeebeeEmail "zeebee-backend/domains/zeebee/email"
)

var EmailClient *email.Email
var EmailAuthManager *zeebeeEmail.Manager
var EmailAuthUrlGenerator map[constants.EmailProvider]zeebeeEmail.AuthCodeUrlGenerator

// initEmailClient initializes the email client with the given RPC configuration.
// It sets up both SMTP configuration for basic email sending and OAuth2 configuration
// for Gmail authentication. The function returns any error encountered during
// initialization.
func initEmailClient(conf config.RpcConf) (err error) {
	// Setup email client
	EmailClient = email.NewEmail(email.SMTPConfig{
		Host:     conf.EmailConf.Host,
		Port:     conf.EmailConf.Port,
		User:     conf.EmailConf.Username,
		Password: conf.EmailConf.Password,
	})

	// Setup Gmail OAuth2 configuration
	gmailAuthConfig := &oauth2.Config{
		ClientID:     conf.EmailAuthConf.Gmail.ClientId,
		ClientSecret: conf.EmailAuthConf.Gmail.ClientSecret,
		RedirectURL:  conf.EmailAuthConf.Gmail.CallbackUrl,
		Scopes: []string{
			gmailapi.GmailSendScope,
			gmailapi.GmailReadonlyScope,
		},
		Endpoint: google.Endpoint,
	}
	EmailAuthUrlGenerator = map[constants.EmailProvider]zeebeeEmail.AuthCodeUrlGenerator{
		constants.EmailProviderGmail: zeebeeEmail.NewGmailAuthCodeUrlGenerator(gmailAuthConfig),
	}

	// Setup email authorization manager
	EmailAuthManager = zeebeeEmail.NewManager(LimiterMgr)

	return nil
}
