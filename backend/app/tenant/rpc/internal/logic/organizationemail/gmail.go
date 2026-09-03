package organizationemaillogic

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmailapi "google.golang.org/api/gmail/v1"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	zeebeeemail "zeebee-backend/domains/zeebee/email"
)

func gmailAuthCodeGenerator(baseLogic *logic.BaseLogic, userType constants.UserType) (*zeebeeemail.GmailAuthCodeUrlGenerator, error) {
	platform, err := baseLogic.PlatformFor(userType)
	if err != nil {
		return nil, err
	}
	gmail := platform.EmailConf.Gmail
	if gmail.ClientId == "" || gmail.ClientSecret == "" || gmail.CallbackUrl == "" {
		return nil, errorx.New("gmail oauth is not configured")
	}
	return zeebeeemail.NewGmailAuthCodeUrlGenerator(&oauth2.Config{
		ClientID:     gmail.ClientId,
		ClientSecret: gmail.ClientSecret,
		RedirectURL:  gmail.CallbackUrl,
		Scopes:       []string{gmailapi.GmailSendScope, gmailapi.GmailReadonlyScope},
		Endpoint:     google.Endpoint,
	}), nil
}
