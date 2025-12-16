package global

import (
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/email"
	zeebeeEmail "zeebee-backend/domains/zeebee/email"
)

var EmailClient *email.Email
var EmailAuthManager *zeebeeEmail.Manager

func initEmailClient(conf config.EmailConf) (err error) {
	EmailClient = email.NewEmail(email.SMTPConfig{
		Host:     conf.Host,
		Port:     conf.Port,
		User:     conf.Username,
		Password: conf.Password,
	})
	return nil
}
