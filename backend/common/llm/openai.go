package llm

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
	openaiopt "github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"

	"zeebee-backend/domains/zeebee/config"
)

const openAiProvider = "openai"

type openAIClient struct {
	client    openai.Client
	model     openai.ChatModel
	maxTokens int64
	env       string // tags uploaded filenames so shared-account files are distinguishable
}

func newOpenAIClient(p *config.LlmProvider, conf config.LlmBizConf, env string) *openAIClient {
	opts := []openaiopt.RequestOption{
		openaiopt.WithAPIKey(p.ApiKey),
		openaiopt.WithRequestTimeout(timeoutOf(conf)),
	}
	// openai-go appends /chat/completions and /responses to base_url, so base needs /v1.
	if base := strings.TrimRight(p.BaseUrl, "/"); base != "" {
		if !strings.HasSuffix(base, "/v1") {
			base += "/v1"
		}
		opts = append(opts, openaiopt.WithBaseURL(base))
	}
	return &openAIClient{client: openai.NewClient(opts...), model: conf.Model, maxTokens: maxTokensOf(conf), env: env}
}

func (c *openAIClient) Provider() string { return openAiProvider }

func (c *openAIClient) Chat(ctx context.Context, system, user string, opts ...ChatOption) (content string, usage Usage, err error) {
	cfg := &chatConfig{}
	for _, o := range opts {
		o(cfg)
	}
	if len(cfg.attachments) == 0 {
		return c.chatText(ctx, system, user)
	}
	return c.chatWithAttachments(ctx, system, user, cfg.attachments)
}

func (c *openAIClient) chatText(ctx context.Context, system, user string) (content string, usage Usage, err error) {
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(system),
			openai.UserMessage(user),
		},
	})
	if err != nil {
		return
	}
	if content, err = openAIFirstContent(resp); err != nil {
		return
	}
	usage = Usage{
		InputTokens:       resp.Usage.PromptTokens,
		OutputTokens:      resp.Usage.CompletionTokens,
		InputCachedTokens: resp.Usage.PromptTokensDetails.CachedTokens,
	}
	return
}

// chatWithAttachments runs the multimodal call via the Responses API — the only
// openai surface accepting uploaded file_ids.
func (c *openAIClient) chatWithAttachments(ctx context.Context, system, user string, atts []Attachment) (content string, usage Usage, err error) {
	parts := make([]responses.ResponseInputContentUnionParam, 0, len(atts)+1)
	for _, a := range atts {
		parts = append(parts, openaiInputPart(a))
	}
	parts = append(parts, responses.ResponseInputContentParamOfInputText(user))
	resp, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model:           string(c.model),
		Instructions:    param.NewOpt(system),
		MaxOutputTokens: param.NewOpt(c.maxTokens),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemParamOfMessage(
					responses.ResponseInputMessageContentListParam(parts),
					responses.EasyInputMessageRoleUser,
				),
			},
		},
	})
	if err != nil {
		return
	}
	content = cleanJSON(resp.OutputText())
	usage = Usage{
		InputTokens:       resp.Usage.InputTokens,
		OutputTokens:      resp.Usage.OutputTokens,
		InputCachedTokens: resp.Usage.InputTokensDetails.CachedTokens,
	}
	return
}

func openaiInputPart(a Attachment) responses.ResponseInputContentUnionParam {
	if a.FileID != "" {
		if isImage(a.ContentType) {
			return responses.ResponseInputContentUnionParam{OfInputImage: &responses.ResponseInputImageParam{FileID: param.NewOpt(a.FileID)}}
		}
		return responses.ResponseInputContentUnionParam{OfInputFile: &responses.ResponseInputFileParam{
			FileID: param.NewOpt(a.FileID),
		}}
	}
	if isImage(a.ContentType) {
		return responses.ResponseInputContentUnionParam{OfInputImage: &responses.ResponseInputImageParam{ImageURL: param.NewOpt("data:" + a.ContentType + ";base64," + a.Base64)}}
	}
	return responses.ResponseInputContentUnionParam{OfInputFile: &responses.ResponseInputFileParam{
		FileData: param.NewOpt(a.Base64),
		Filename: param.NewOpt(a.Filename),
	}}
}

// namedReader exposes Filename/ContentType so the SDK's multipart encoder uses
// the real values instead of anonymous_file / application/octet-stream.
type namedReader struct {
	io.Reader
	filename    string
	contentType string
}

func (r *namedReader) Filename() string    { return r.filename }
func (r *namedReader) ContentType() string { return r.contentType }

// UploadFile stores a binary via the Files API; PDFs as user_data (Responses
// input_file), images as vision (input_image).
func (c *openAIClient) UploadFile(ctx context.Context, reader io.Reader, filename, contentType string) (string, string, error) {
	purpose := openai.FilePurposeUserData
	if isImage(contentType) {
		purpose = openai.FilePurposeVision
	}
	file, err := c.client.Files.New(ctx, openai.FileNewParams{
		File:    &namedReader{Reader: reader, filename: EnvFilename(c.env, filename), contentType: contentType},
		Purpose: purpose,
	})
	if err != nil {
		return "", "", err
	}
	return openAiProvider, file.ID, nil
}

func (c *openAIClient) DeleteFile(ctx context.Context, fileID string) error {
	_, err := c.client.Files.Delete(ctx, fileID)
	return err
}

// ListFiles returns one page of uploaded files. OpenAI has no scope filter; the
// caller filters client-side if it tagged filenames (see EnvFilename).
func (c *openAIClient) ListFiles(ctx context.Context, limit int, afterID string) ([]FileInfo, bool, error) {
	page, err := c.client.Files.List(ctx, openAIListParams(limit, afterID))
	if err != nil {
		return nil, false, err
	}
	out := make([]FileInfo, len(page.Data))
	for i, f := range page.Data {
		out[i] = FileInfo{
			ID:        f.ID,
			Filename:  f.Filename,
			Bytes:     f.Bytes,
			CreatedAt: time.Unix(f.CreatedAt, 0),
		}
	}
	return out, page.HasMore, nil
}

func openAIFirstContent(resp *openai.ChatCompletion) (string, error) {
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}

// openAIListParams builds the Files list query: limit<=0 and empty afterID
// fall back to the provider default (no opt set).
func openAIListParams(limit int, afterID string) openai.FileListParams {
	params := openai.FileListParams{}
	if limit > 0 {
		params.Limit = param.NewOpt(int64(limit))
	}
	if afterID != "" {
		params.After = param.NewOpt(afterID)
	}
	return params
}
