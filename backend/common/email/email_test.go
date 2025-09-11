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
	msg := EmailMessage{
		From:    "info@zeebeeloans.com",
		To:      []string{"wei@yeelee-capital.com"},
		Subject: "测试邮件",
		Body:    body,
		IsHTML:  true,
	}
	err := NewEmail(c).SendEmail(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("发送成功")
}

var body = `<html lang="en">
    
    <head>
        <META http-equiv="Content-Type" content="text/html; charset=UTF-8">
            
        <meta charset="UTF-8">
            
        <meta content="width=device-width, initial-scale=1.0" name="viewport">
            
        <title>Password Reset Email</title>
            
        <style>
        * {
            margin: 0;
            padding: 0;
        }
        body {
            color: rgba(0, 0, 0, 0.65);
            font-size: 14px;
            line-height: 22px;
            font-family: Roboto, sans-serif;
        }
        .content {
            max-width: 600px;
            margin: 0 auto;
            padding: 32px 24px 0;
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 28px;
            text-align: center;
        }

        .content img {
            height: 28px;
        }

        .content h2 {
            font-size: 20px;
            font-weight: 500;
            line-height: 28px;
            align-self: flex-start;
            color: #000;
        }
        .content > p {
            padding-bottom: 20px;
            color: #000;
        }
        .button {
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
        }
        .button:hover {
            background-color: #1f88ff;
        }
        .tips {
            color: rgba(0, 0, 0, 0.65);
        }
        .footer {
            width: 100%;
            padding: 24px;
            font-size: 12px;
            background-color: #f5f5f5;
            border-top: 1px solid rgba(5, 145, 255, 0.1);
            margin-top: 32px;
        }
    </style>
        
    </head>
    
    <body>
        
        <div class="content">
                <img alt="logo" src="https://zeebee-files-staging.s3.us-west-1.amazonaws.com/files-dev/logo.png">
    
            <h2>Hello zeng1,</h2>
                
            <p>
        We received a sign-up request for ZeeBee using this email address.
        However, we found that you already have an existing account with us.
    </p>
                
            <div class="tips">
        To access your account, simply sign in with your existing credentials.
    </div>
                <a class="button" href="http://localhost:3000/login">Sign in to your account</a>

    
            <div class="footer">
        ZeeBee | 123 Business St | San Francisco, CA, 94105
    </div>
            
        </div>
        
    </body>
    
</html>
`
