package global

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	driveapi "google.golang.org/api/drive/v3"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/domains/zeebee/drive"
)

// DriveOAuthConfig is the shared Google Drive OAuth2 config (nil when unconfigured).
var DriveOAuthConfig *oauth2.Config

// DriveAuthUrlGenerator builds Drive consent URLs (nil when unconfigured).
var DriveAuthUrlGenerator *drive.DriveAuthCodeUrlGenerator

// initDrive sets up the Google Drive OAuth2 config and consent URL generator.
// When no client id is configured the integration is left disabled (configs nil).
func initDrive(conf config.RpcConf) (err error) {
	if conf.DriveAuthConf.Google.ClientId == "" {
		return nil
	}
	DriveOAuthConfig = &oauth2.Config{
		ClientID:     conf.DriveAuthConf.Google.ClientId,
		ClientSecret: conf.DriveAuthConf.Google.ClientSecret,
		RedirectURL:  conf.DriveAuthConf.Google.CallbackUrl,
		Scopes:       []string{driveapi.DriveScope},
		Endpoint:     google.Endpoint,
	}
	DriveAuthUrlGenerator = drive.NewDriveAuthCodeUrlGenerator(DriveOAuthConfig)
	return nil
}
