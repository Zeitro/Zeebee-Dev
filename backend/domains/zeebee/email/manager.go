package email

import (
	"context"
	"sync"

	commonemail "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/email/client"
	"zeebee-backend/domains/zeebee/limiter"
)

type Manager struct {
	mu         sync.RWMutex
	clients    map[string]client.Client
	limiterMgr *limiter.LimiterManager
}

// EmailType constants for different email types
type EmailType string

func NewManager(mgr *limiter.LimiterManager) *Manager {
	return &Manager{
		clients:    make(map[string]client.Client),
		limiterMgr: mgr,
	}
}

// BuildKey builds a key for the limiter
func (m *Manager) BuildKey(strategy limiter.LimiterStrategy, email string, recipient string, organizationId int64, organizationType constants.UserType) string {
	return stringx.Format("email:%s:%s:%s:%d:%s:limiter", strategy, email, recipient, organizationId, organizationType)
}

// SetLimiter sets the email limiter for the email manager
func (m *Manager) SetLimiter(mgr *limiter.LimiterManager) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limiterMgr = mgr
}

func (m *Manager) Register(email string, emailClient client.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[email] = emailClient
}

func (m *Manager) GetClient(email string) (emailClient client.Client, existed bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.clients[email]
	if !ok {
		return nil, false
	}
	return c, true
}

// SendMail sends an email with rate limiting based on the provided email type and rate limits
func (m *Manager) SendMail(ctx context.Context, strategy limiter.LimiterStrategy, limiterKey string, email string, msg commonemail.Message) error {
	if m.limiterMgr != nil {
		allowed, err := m.limiterMgr.Limit(ctx, strategy, limiterKey)
		if err != nil {
			return errorx.WithStack(err)
		}
		if !allowed {
			return errorx.New("Rate Limit Exceeded")
		}
	}

	emailClient, existed := m.GetClient(email)
	if !existed {
		return errorx.New("Email Not Found")
	}
	return emailClient.SendEmail(ctx, msg)
}

// SendMailWithoutLimit sends an email without rate limiting - for backward compatibility
func (m *Manager) SendMailWithoutLimit(ctx context.Context, email string, msg commonemail.Message) error {
	emailClient, existed := m.GetClient(email)
	if !existed {
		return errorx.New("Email Not Found")
	}
	return emailClient.SendEmail(ctx, msg)
}
