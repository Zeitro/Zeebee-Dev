package llm

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"zeebee-backend/domains/zeebee/config"
)

const claudeProvider = "claude"

type claudeClient struct {
	client    anthropic.Client
	model     anthropic.Model
	maxTokens int64
}

func newClaudeClient(p *config.LlmProvider, biz config.LlmBizConf) *claudeClient {
	opts := []option.RequestOption{
		option.WithAPIKey(p.ApiKey),
		option.WithRequestTimeout(timeoutOf(biz)),
	}
	if base := strings.TrimRight(p.BaseUrl, "/"); base != "" {
		opts = append(opts, option.WithBaseURL(base))
	}
	return &claudeClient{
		client:    anthropic.NewClient(opts...),
		model:     biz.Model,
		maxTokens: maxTokensOf(biz),
	}
}
func (c *claudeClient) Provider() string { return claudeProvider }

func (c *claudeClient) Chat(ctx context.Context, system, user string, opts ...ChatOption) (string, error) {
	cfg := &chatConfig{}
	for _, o := range opts {
		o(cfg)
	}
	if len(cfg.attachments) == 0 {
		return c.chatText(ctx, system, user)
	}
	return c.chatWithAttachments(ctx, system, user, cfg.attachments)
}

func (c *claudeClient) chatText(ctx context.Context, system, user string) (string, error) {
	resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		System:    []anthropic.TextBlockParam{{Text: system}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(user)),
		},
	})
	if err != nil {
		return "", err
	}
	return claudeText(resp)
}

// chatWithAttachments is the consolidated multimodal call: each attachment
// becomes a beta content block (fileId or base64, document or image), followed
// by the text prompt. Uses the beta Messages API since file-id sources are
// beta-only.
func (c *claudeClient) chatWithAttachments(ctx context.Context, system, user string, atts []Attachment) (string, error) {
	blocks := make([]anthropic.BetaContentBlockParamUnion, 0, len(atts)+1)
	for _, a := range atts {
		blocks = append(blocks, claudeBetaBlock(a))
	}
	blocks = append(blocks, anthropic.NewBetaTextBlock(user))
	resp, err := c.client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		System:    []anthropic.BetaTextBlockParam{{Text: system}},
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(blocks...),
		},
	},
		// file-id document/image sources need the files-api beta header on the
		// messages call; Beta.Files.* add it automatically, Beta.Messages does not.
		option.WithHeader("anthropic-beta", "files-api-2025-04-14"),
	)
	if err != nil {
		return "", err
	}
	return betaClaudeText(resp)
}

// claudeBetaBlock maps one attachment to a beta content block by reference kind
// (fileId vs base64) and content type (document vs image).
func claudeBetaBlock(a Attachment) anthropic.BetaContentBlockParamUnion {
	if a.FileID != "" {
		if isImage(a.ContentType) {
			return anthropic.NewBetaImageBlock(anthropic.BetaFileImageSourceParam{FileID: a.FileID})
		}
		return anthropic.NewBetaDocumentBlock(anthropic.BetaFileDocumentSourceParam{FileID: a.FileID})
	}
	if isImage(a.ContentType) {
		return anthropic.NewBetaImageBlock(anthropic.BetaBase64ImageSourceParam{
			Data:      a.Base64,
			MediaType: anthropic.BetaBase64ImageSourceMediaType(a.ContentType),
		})
	}
	return anthropic.NewBetaDocumentBlock(anthropic.BetaBase64PDFSourceParam{Data: a.Base64})
}

// UploadFile uses the beta Files API to store a binary and returns its id. The
// id is later referenced by Chat via a beta document/image source.
func (c *claudeClient) UploadFile(ctx context.Context, reader io.Reader, filename, contentType string) (string, string, error) {
	_ = contentType // claude infers type from the uploaded filename/bytes
	resp, err := c.client.Beta.Files.Upload(ctx, anthropic.BetaFileUploadParams{File: reader})
	if err != nil {
		return "", "", err
	}
	return claudeProvider, resp.ID, nil
}

// DeleteFile removes an uploaded file via the beta Files API.
func (c *claudeClient) DeleteFile(ctx context.Context, fileID string) error {
	_, err := c.client.Beta.Files.Delete(ctx, fileID, anthropic.BetaFileDeleteParams{})
	return err
}

func claudeText(resp *anthropic.Message) (string, error) {
	var sb strings.Builder
	for _, block := range resp.Content {
		if t := block.AsText(); t.Text != "" {
			sb.WriteString(t.Text)
		}
	}
	if sb.Len() == 0 {
		return "", fmt.Errorf("claude returned no text")
	}
	return cleanJSON(sb.String()), nil
}

func betaClaudeText(resp *anthropic.BetaMessage) (string, error) {
	var sb strings.Builder
	for _, block := range resp.Content {
		if t := block.AsText(); t.Text != "" {
			sb.WriteString(t.Text)
		}
	}
	if sb.Len() == 0 {
		return "", fmt.Errorf("claude returned no text")
	}
	return cleanJSON(sb.String()), nil
}
