package client

import (
	"context"
	"sync"

	"golang.org/x/oauth2"
	gmailapi "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	commonemail "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

func GenerateGmailAuthCodeUrl(config *oauth2.Config, email string) string {
	return config.AuthCodeURL(email, oauth2.AccessTypeOffline)
}

type GmailClient struct {
	mu        sync.Mutex
	config    *oauth2.Config
	token     *oauth2.Token
	emailType constants.EmailType
}

func NewGmailClient(config *oauth2.Config, code string) (client *GmailClient, err error) {
	token, err := config.Exchange(context.Background(), code, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, errorx.Wrap(err, "Exchange Token Failed")
	}
	return &GmailClient{
		mu:        sync.Mutex{},
		config:    config,
		token:     token,
		emailType: constants.EmailTypeGmail,
	}, nil
}

// GetEmail
func (cli *GmailClient) GetEmail(ctx context.Context) (email string, err error) {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return "", err
	}
	profile, err := srv.Users.GetProfile("me").Do()
	if err != nil {
		return "", err
	}
	return profile.EmailAddress, nil
}

// RefreshToken
func (cli *GmailClient) RefreshToken(ctx context.Context) (token *oauth2.Token, err error) {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	if cli.token == nil {
		return nil, errorx.New("Token Is Nil")
	}

	if cli.token.RefreshToken == "" {
		return nil, errorx.New("Refresh Token Is Nil")
	}
	if cli.token.Valid() {
		return cli.token, nil
	}
	ts := cli.config.TokenSource(ctx, cli.token)
	token, err = ts.Token()
	if err != nil {
		return nil, errorx.Wrap(err, "Refresh Token Failed")
	}
	if token == nil {
		return nil, errorx.New("New Token Is Nil")
	}

	cli.token = token
	return token, nil
}

// SendEmail
func (cli *GmailClient) SendEmail(ctx context.Context, msg commonemail.Message) (err error) {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return err
	}
	raw, err := buildRawMessage(msg)
	if err != nil {
		return errorx.Wrap(err, "Build Raw Message Failed")
	}
	_, err = srv.Users.Messages.Send("me", &gmailapi.Message{Raw: raw}).Do()
	return err
}

// getSrv
func (cli *GmailClient) getSrv(ctx context.Context) (srv *gmailapi.Service, err error) {
	httpClient := cli.config.Client(context.Background(), cli.token)
	srv, err = gmailapi.NewService(ctx,
		option.WithHTTPClient(httpClient),
	)
	return srv, errorx.WithStack(err)
}

// GetEmailType
func (cli *GmailClient) GetEmailType() constants.EmailType {
	return cli.emailType
}
