// Package llm provides a small provider-agnostic chat client (openai / claude)
// used for structured (JSON) text and multimodal file extraction. Shared by file
// parsing and lender discover so both run on the same provider/config plumbing.
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

// Attachment is a provider-agnostic file reference passed into a Chat call.
// Set exactly one of FileID (an already-uploaded provider file) or Base64
// (inline content); Filename + ContentType describe it for every provider. This
// lets claude and openai share one call surface — each maps an attachment to its
// own multimodal API. Both upload-based (fileId, for resume reuse + terminal
// delete) and upload-free (base64) flows go through the same type.
type Attachment struct {
	FileID      string // provider file id from a prior UploadFile
	Base64      string // inline base64-encoded content (no upload needed)
	Filename    string
	ContentType string // mime type, e.g. application/pdf, image/png
}

// chatConfig holds the optional attachments collected from ChatOptions.
type chatConfig struct {
	attachments []Attachment
}

// ChatOption configures optional file attachments on a Chat call.
type ChatOption func(*chatConfig)

// WithAttachment attaches one file (fileId or base64) to the next Chat call.
func WithAttachment(a Attachment) ChatOption {
	return func(c *chatConfig) { c.attachments = append(c.attachments, a) }
}

// Client abstracts the LLM so openai and claude are interchangeable. The single
// Chat entry point takes optional file attachments; no attachment = plain text.
type Client interface {
	// Provider returns this client's provider name (claude, openai). An uploaded
	// file id is provider-scoped, so callers match it on resume (research R11).
	Provider() string
	// Chat sends system+user with optional file attachments. No attachments is a
	// plain text chat; each provider maps attachments (fileId or base64 +
	// filename + contentType) to its own multimodal API.
	Chat(ctx context.Context, system, user string, opts ...ChatOption) (string, error)
	// UploadFile stores a binary at the provider and returns the provider name +
	// id to reference later (resume reuse, terminal delete). The id is NOT
	// portable across providers (research R11): a provider switch invalidates
	// prior ids, so callers persist the provider alongside the id and re-upload
	// on mismatch.
	UploadFile(ctx context.Context, reader io.Reader, filename, contentType string) (provider, fileID string, err error)
	// DeleteFile removes an uploaded provider file. Called at the terminal state
	// to stop billing; failures are logged, never blocking.
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

// isImage reports whether a content type is an image attachment.
func isImage(contentType string) bool { return strings.HasPrefix(contentType, "image/") }

// cleanJSON strips markdown code fences (```json ... ```) and surrounding prose
// that claude may emit since it has no native json_object response format.
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
