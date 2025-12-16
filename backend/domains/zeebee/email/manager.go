package email

import (
	"context"
	"sync"

	commonemail "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/email/client"
)

// var _manager *Manager
//
// func init() {
// 	_manager = NewManager()
// }
//
// // Register
// func Register(email string, emailClient client.Client) {
// 	_manager.Register(email, emailClient)
// }
//
// // SendMail
// func SendMail(ctx context.Context, email string, msg commonemail.Message) error {
// 	return _manager.SendMail(ctx, email, msg)
// }
//
// // GetClient
// func GetClient(email string) (emailClient client.Client, existed bool) {
// 	return _manager.GetClient(email)
// }

type Manager struct {
	mu      sync.RWMutex
	clients map[string]client.Client
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]client.Client),
	}
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

func (m *Manager) SendMail(ctx context.Context, email string, msg commonemail.Message) error {
	emailClient, existed := m.GetClient(email)
	if !existed {
		return errorx.New("Email Not Found")
	}
	return emailClient.SendEmail(ctx, msg)
}
