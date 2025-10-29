package global

import (
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/email"
)

var EmailClient *email.Email

func initEmailClient(conf config.EmailConf) (err error) {
	EmailClient = email.NewEmail(email.SMTPConfig{
		Host:     conf.Host,
		Port:     conf.Port,
		User:     conf.Username,
		Password: conf.Password,
	})
	return nil
}
