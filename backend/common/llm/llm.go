// Package llm provides a small provider-agnostic chat client (openai / claude)
// used for structured (JSON) text and vision extraction. Shared by file parsing
// and lender discover so both run on the same provider/config plumbing.
package llm

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	anthopt "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/openai/openai-go/v3"
	openaiopt "github.com/openai/openai-go/v3/option"

	"zeebee-backend/domains/zeebee/config"
)

const (
	defaultMaxTokens  = int64(4096)
	defaultTimeoutSec = 60
)

// Client abstracts the LLM used for structured extraction so openai and claude
// implementations are interchangeable.
type Client interface {
	// ChatJSON sends a text prompt and returns raw model content (expected JSON).
	ChatJSON(ctx context.Context, system, user string) (string, error)
	// VisionJSON sends a text prompt plus one PNG image and returns raw content.
	VisionJSON(ctx context.Context, system, text string, imagePNG []byte) (string, error)
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
	case "openai":
		return newOpenAIClient(prov, biz), nil
	case "claude":
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

// --- openai ---

type openAIClient struct {
	client openai.Client
	model  openai.ChatModel
}

func newOpenAIClient(p *config.LlmProvider, biz config.LlmBizConf) *openAIClient {
	opts := []openaiopt.RequestOption{
		openaiopt.WithAPIKey(p.ApiKey),
		openaiopt.WithRequestTimeout(timeoutOf(biz)),
	}
	// openai-go appends "/chat/completions" to base_url, so base must carry "/v1".
	if base := strings.TrimRight(p.BaseUrl, "/"); base != "" {
		if !strings.HasSuffix(base, "/v1") {
			base += "/v1"
		}
		opts = append(opts, openaiopt.WithBaseURL(base))
	}
	return &openAIClient{client: openai.NewClient(opts...), model: openai.ChatModel(biz.Model)}
}

func (c *openAIClient) ChatJSON(ctx context.Context, system, user string) (string, error) {
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(system),
			openai.UserMessage(user),
		},
	})
	if err != nil {
		return "", err
	}
	return openAIFirstContent(resp)
}

func (c *openAIClient) VisionJSON(ctx context.Context, system, text string, imagePNG []byte) (string, error) {
	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(imagePNG)
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(system),
			openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
				openai.TextContentPart(text),
				openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
					URL:    dataURL,
					Detail: "high",
				}),
			}),
		},
	})
	if err != nil {
		return "", err
	}
	return openAIFirstContent(resp)
}

func openAIFirstContent(resp *openai.ChatCompletion) (string, error) {
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}

// --- claude ---

type claudeClient struct {
	client    anthropic.Client
	model     anthropic.Model
	maxTokens int64
}

func newClaudeClient(p *config.LlmProvider, biz config.LlmBizConf) *claudeClient {
	opts := []anthopt.RequestOption{
		anthopt.WithAPIKey(p.ApiKey),
		anthopt.WithRequestTimeout(timeoutOf(biz)),
	}
	if base := strings.TrimRight(p.BaseUrl, "/"); base != "" {
		opts = append(opts, anthopt.WithBaseURL(base))
	}
	return &claudeClient{
		client:    anthropic.NewClient(opts...),
		model:     anthropic.Model(biz.Model),
		maxTokens: maxTokensOf(biz),
	}
}

func (c *claudeClient) ChatJSON(ctx context.Context, system, user string) (string, error) {
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

func (c *claudeClient) VisionJSON(ctx context.Context, system, text string, imagePNG []byte) (string, error) {
	resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		System:    []anthropic.TextBlockParam{{Text: system}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(text),
				anthropic.NewImageBlockBase64("image/png", base64.StdEncoding.EncodeToString(imagePNG)),
			),
		},
	})
	if err != nil {
		return "", err
	}
	return claudeText(resp)
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
