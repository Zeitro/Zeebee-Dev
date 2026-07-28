package llm

import (
	"context"
	"fmt"
	"io"
	"sync"
)

// MockClient is a fake Client for tests: it returns canned content and records
// every call so tests assert attachment passthrough, upload ids, and deletes.
// All fields are safe for concurrent use (a parse run uploads many files).
type MockClient struct {
	// ChatContent is returned by Chat (default "{}" when empty).
	ChatContent string
	ChatErr     error

	// UploadProvider/UploadFileID are returned by UploadFile; a monotonic suffix
	// makes each upload distinct when UploadFileID is empty.
	UploadProvider string
	UploadFileID   string
	UploadErr      error

	DeleteErr error

	mu            sync.Mutex
	seen          []Attachment // attachments seen across all Chat calls
	chatCalls     int          // Chat invocations — proves single vs per-file call
	UploadCount   int
	UploadedFiles []string // reader-driven uploads (content read for fidelity)
	DeletedFiles  []string
}

// Provider returns the mock provider name (default "mock").
func (m *MockClient) Provider() string {
	if m.UploadProvider != "" {
		return m.UploadProvider
	}
	return "mock"
}

// Chat records the attachments from opts and returns the canned content.
func (m *MockClient) Chat(_ context.Context, _, _ string, opts ...ChatOption) (string, Usage, error) {
	cfg := &chatConfig{}
	for _, o := range opts {
		o(cfg)
	}
	m.mu.Lock()
	m.seen = append(m.seen, cfg.attachments...)
	m.chatCalls++
	m.mu.Unlock()
	if m.ChatErr != nil {
		return "", Usage{}, m.ChatErr
	}
	if m.ChatContent != "" {
		return m.ChatContent, Usage{}, nil
	}
	return "{}", Usage{}, nil
}

// UploadFile returns the configured provider + a distinct file id per call.
func (m *MockClient) UploadFile(_ context.Context, _ io.Reader, _, _ string) (string, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.UploadCount++
	if m.UploadErr != nil {
		return "", "", m.UploadErr
	}
	provider := m.UploadProvider
	if provider == "" {
		provider = "mock"
	}
	id := m.UploadFileID
	if id == "" {
		id = fmt.Sprintf("file-%d", m.UploadCount)
	}
	m.UploadedFiles = append(m.UploadedFiles, id)
	return provider, id, nil
}

// DeleteFile records the id and returns the configured error.
func (m *MockClient) DeleteFile(_ context.Context, fileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DeletedFiles = append(m.DeletedFiles, fileID)
	return m.DeleteErr
}

// ListFiles returns one page of UploadedFiles, mirroring provider cursor paging:
// afterID selects the start position (empty = first page), limit caps the page
// size (<=0 = all remaining), hasMore is set when items remain beyond the page.
// An afterID not present in UploadedFiles yields an empty page (cursor exhausted).
func (m *MockClient) ListFiles(_ context.Context, limit int, afterID string) ([]FileInfo, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	start := 0
	if afterID != "" {
		idx := -1
		for i, id := range m.UploadedFiles {
			if id == afterID {
				idx = i
				break
			}
		}
		if idx < 0 {
			return nil, false, nil
		}
		start = idx + 1
	}
	rest := m.UploadedFiles[start:]
	hasMore := false
	if limit > 0 && len(rest) > limit {
		rest = rest[:limit]
		hasMore = true
	}
	out := make([]FileInfo, len(rest))
	for i, id := range rest {
		out[i] = FileInfo{ID: id}
	}
	return out, hasMore, nil
}

// Attachments returns a copy of every attachment seen by Chat.
func (m *MockClient) Attachments() []Attachment {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Attachment, len(m.seen))
	copy(out, m.seen)
	return out
}

// ChatCalls returns the number of Chat invocations (single vs per-file call proof).
func (m *MockClient) ChatCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.chatCalls
}
