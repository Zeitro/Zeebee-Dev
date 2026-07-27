// Package llm is a provider-agnostic chat client (openai/claude) for structured
// text and multimodal file extraction. Shared by file parsing and lender discover.
package llm

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"zeebee-backend/domains/zeebee/config"
)

const (
	defaultMaxTokens  = int64(4096)
	defaultTimeoutSec = 60
)

// Attachment is a provider-agnostic file reference for a Chat call. Set exactly
// one of FileID (prior UploadFile) or Base64 (inline); claude and openai each
// map it to their own multimodal API.
type Attachment struct {
	FileID      string // provider file id from a prior UploadFile
	Base64      string // inline base64-encoded content (no upload needed)
	Filename    string
	ContentType string // mime type, e.g. application/pdf, image/png
}

type chatConfig struct {
	attachments []Attachment
}

// Usage is the provider-agnostic token accounting returned by Chat.
type Usage struct {
	InputTokens       int64 `json:"input_tokens"`
	OutputTokens      int64 `json:"output_tokens"`
	InputCachedTokens int64 `json:"input_cached_tokens"`
}

// ChatOption configures optional file attachments on a Chat call.
type ChatOption func(*chatConfig)

// WithAttachment attaches one file (fileId or base64) to the next Chat call.
func WithAttachment(a Attachment) ChatOption {
	return func(c *chatConfig) { c.attachments = append(c.attachments, a) }
}

// Client abstracts the LLM so providers are interchangeable.
type Client interface {
	// Provider returns the provider name (claude, openai); file ids are
	// provider-scoped, so callers match it on resume.
	Provider() string
	// Chat sends system+user with optional file attachments; no attachments = plain text.
	Chat(ctx context.Context, system, user string, opts ...ChatOption) (string, Usage, error)
	// UploadFile stores a binary and returns provider + id to reference later
	// (resume reuse, terminal delete). The id is NOT portable across providers —
	// persist the provider alongside and re-upload on mismatch.
	UploadFile(ctx context.Context, reader io.Reader, filename, contentType string) (provider, fileID string, err error)
	// DeleteFile removes an uploaded file; called at terminal state to stop
	// billing. Failures logged, never blocking.
	DeleteFile(ctx context.Context, fileID string) error
}

// NewClient resolves the provider named by biz.Provider from the configured
// provider list. Unknown providers fail fast.
func NewClient(providers []config.LlmProvider, biz config.LlmBizConf) (Client, error) {
	var prov *config.LlmProvider
	for i := range providers {
		if providers[i].Name == biz.Provider {
			prov = &providers[i]
			break
		}
	}
	if prov == nil {
		return nil, fmt.Errorf("llm provider %q not configured", biz.Provider)
	}
	switch prov.Name {
	case openAiProvider:
		return newOpenAIClient(prov, biz), nil
	case claudeProvider:
		return newClaudeClient(prov, biz), nil
	default:
		return nil, fmt.Errorf("llm provider %q not supported", prov.Name)
	}
}

func maxTokensOf(biz config.LlmBizConf) int64 {
	if biz.MaxTokens > 0 {
		return biz.MaxTokens
	}
	return defaultMaxTokens
}

func timeoutOf(biz config.LlmBizConf) time.Duration {
	if biz.TimeoutSeconds > 0 {
		return time.Duration(biz.TimeoutSeconds) * time.Second
	}
	return time.Duration(defaultTimeoutSec) * time.Second
}

func isImage(contentType string) bool { return strings.HasPrefix(contentType, "image/") }

// cleanJSON strips markdown fences and prose: claude has no json_object format.
func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	// drop opening fence (with optional language tag)
	if nl := strings.IndexByte(s, '\n'); nl >= 0 {
		s = s[nl+1:]
	}
	s = strings.TrimSuffix(strings.TrimSpace(s), "```")
	return strings.TrimSpace(s)
}
