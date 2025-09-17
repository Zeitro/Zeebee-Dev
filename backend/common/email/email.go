package email

import (
	"io"

	"gopkg.in/gomail.v2"
)

// attachment
type EmailAttachment struct {
	Filename string
	Reader   io.Reader
	MimeType string // optional, e.g. "application/pdf"
}

type EmailMessage struct {
	Subject     string
	From        string
	To          []string
	Body        string
	Cc          []string
	Attachments []EmailAttachment
	IsHTML      bool
}

// SMTP config
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	// Proxy string
}

type Email struct {
	SMTP   SMTPConfig
	dialer *gomail.Dialer
}

func NewEmail(smtp SMTPConfig) *Email {
	d := gomail.NewDialer(smtp.Host, smtp.Port, smtp.User, smtp.Password)
	return &Email{SMTP: smtp, dialer: d}
}

// SendEmail sends an email using the provided SMTP configuration.
func (e *Email) SendEmail(msg EmailMessage) error {
	m := gomail.NewMessage()
	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	if len(msg.Cc) > 0 {
		m.SetHeader("Cc", msg.Cc...)
	}
	if msg.IsHTML {
		m.SetBody("text/html", msg.Body)
	} else {
		m.SetBody("text/plain", msg.Body)
	}
	for _, att := range msg.Attachments {
		if att.Reader != nil && att.Filename != "" {
			if att.MimeType != "" {
				m.Attach(att.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
					_, err := io.Copy(w, att.Reader)
					return err
				}), gomail.SetHeader(map[string][]string{"Content-Type": {att.MimeType}}))
			} else {
				m.Attach(att.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
					_, err := io.Copy(w, att.Reader)
					return err
				}))
			}
		}
	}

	err := e.dialer.DialAndSend(m)
	// if err != nil {
	// 	fmt.Printf("send email error: %v,config: %v,msg: %v\n", err, e.SMTP, msg.From)
	// }

	return err
}
