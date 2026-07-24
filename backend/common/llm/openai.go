package llm

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/openai/openai-go/v3"
	openaiopt "github.com/openai/openai-go/v3/option"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/config"
)

const openAiProvider = "openai"

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
	return &openAIClient{client: openai.NewClient(opts...), model: biz.Model}
}

func (c *openAIClient) Provider() string { return openAiProvider }

func (c *openAIClient) Chat(ctx context.Context, system, user string, opts ...ChatOption) (string, error) {
	cfg := &chatConfig{}
	for _, o := range opts {
		o(cfg)
	}
	if len(cfg.attachments) == 0 {
		return c.chatText(ctx, system, user)
	}
	// openai chat completions support base64 image attachments; pdf/file-id
	// multimodal needs the responses API and is out of scope (file-parse runs on
	// the claude provider). Be explicit rather than silently dropping a file.
	parts := []openai.ChatCompletionContentPartUnionParam{openai.TextContentPart(user)}
	for _, a := range cfg.attachments {
		if a.Base64 == "" || !isImage(a.ContentType) {
			return "", errorx.New("openai chat supports base64 image attachments only")
		}
		dataURL := "data:" + a.ContentType + ";base64," + a.Base64
		parts = append(parts, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL:    dataURL,
			Detail: "high",
		}))
	}
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(system),
			openai.UserMessage(parts),
		},
	})
	if err != nil {
		return "", err
	}
	return openAIFirstContent(resp)
}

func (c *openAIClient) chatText(ctx context.Context, system, user string) (string, error) {
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

// UploadFile/DeleteFile are unsupported on openai: file upload + multimodal
// file-id reference run on the claude provider only.
func (c *openAIClient) UploadFile(context.Context, io.Reader, string, string) (string, string, error) {
	return "", "", errorx.New("openai provider does not support file upload")
}

func (c *openAIClient) DeleteFile(context.Context, string) error {
	return errorx.New("openai provider does not support file delete")
}

func openAIFirstContent(resp *openai.ChatCompletion) (string, error) {
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}
