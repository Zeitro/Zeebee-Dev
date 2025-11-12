package email

import (
	"testing"
)

func TestEmail_SendEmail(t *testing.T) {
	c := SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		User:     "info@zeebeeloans.com",
		Password: "fyxxgqlqqxzitzef",
	}
	msg := Message{
		From:    "info@zeebeeloans.com",
		To:      []string{"wei@yeelee-capital.com"},
		Subject: "test email",
		Body:    body,
		IsHTML:  true,
	}
	err := NewEmail(c).SendEmail(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send successfully")
}

var body = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>ZeeBee User Email</title>
  </head>
  <body>
    <div
      style="
        max-width: 600px;
        margin: 0 auto;
        padding: 32px 24px 0;
        text-align: center;
        color: rgba(0, 0, 0, 0.65);
        font-size: 14px;
        font-family: Roboto, sans-serif;
      "
    >
      <img
        src="https://zeebee-files-staging.s3.us-west-1.amazonaws.com/files-dev/logo.png"
        alt="logo"
        style="height: 28px"
      />
      <div
        style="
          font-size: 20px;
          font-weight: 500;
          line-height: 28px;
          color: #000;
          font-family: Roboto, sans-serif;
          padding: 28px 0;
          text-align: left;
        "
      >
        Hello {User first name},
      </div>

      <div
        style="
          color: #000;
          font-family: Roboto, sans-serif;
          font-size: 14px;
        "
      >
        We received a sign-up request for ZeeBee using this email address.
        However, we found that you already have an existing account with us.
      </div>

      <div
        style="
          color: rgba(0, 0, 0, 0.65);
          font-family: Roboto, sans-serif;
          font-size: 14px;
          padding: 48px 0 28px;
        "
      >
        To access your account, simply sign in with your existing credentials.
      </div>
      <a
        href="#"
        style="
          display: inline-block;
          padding: 9px 20px;
          background-color: #1677ff;
          color: white;
          text-decoration: none;
          border-radius: 8px;
          transition: background-color 0.3s ease;
          width: max-content;
          font-size: 16px;
          box-shadow: 0 2px 0 0 rgba(5, 145, 255, 0.1);
          font-family: Roboto, sans-serif;
        "
        >Sign in to your account</a
      >

      <div
        style="
          width: 100%;
          padding: 24px;
          font-size: 12px;
          background-color: #f5f5f5;
          border-top: 1px solid rgba(5, 145, 255, 0.1);
          margin-top: 60px;
          font-family: Roboto, sans-serif;
          color: rgba(0, 0, 0, 0.65);
        "
      >
        ZeeBee | 123 Business St | San Francisco, CA, 94105
      </div>
    </div>
  </body>
</html>
`
