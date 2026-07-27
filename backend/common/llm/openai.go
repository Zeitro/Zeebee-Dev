package llm

import (
	"context"
	"fmt"
	"io"
	"strings"

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
}

func newOpenAIClient(p *config.LlmProvider, biz config.LlmBizConf) *openAIClient {
	opts := []openaiopt.RequestOption{
		openaiopt.WithAPIKey(p.ApiKey),
		openaiopt.WithRequestTimeout(timeoutOf(biz)),
	}
	// openai-go appends "/chat/completions" (and "/responses") to base_url, so
	// base must carry "/v1".
	if base := strings.TrimRight(p.BaseUrl, "/"); base != "" {
		if !strings.HasSuffix(base, "/v1") {
			base += "/v1"
		}
		opts = append(opts, openaiopt.WithBaseURL(base))
	}
	return &openAIClient{client: openai.NewClient(opts...), model: biz.Model, maxTokens: maxTokensOf(biz)}
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
	return c.chatWithAttachments(ctx, system, user, cfg.attachments)
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

// chatWithAttachments runs the multimodal call through the Responses API — the
// only openai surface that accepts uploaded file_ids (PDF via input_file,
// image via input_image). Chat Completions carries base64 images only.
func (c *openAIClient) chatWithAttachments(ctx context.Context, system, user string, atts []Attachment) (string, error) {
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
		return "", err
	}
	return cleanJSON(resp.OutputText()), nil
}

// openaiInputPart maps one attachment to a Responses input part by reference
// kind (fileId vs base64) and content type (file vs image).
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

// namedReader wraps an io.Reader with Filename() and ContentType() methods so
// the OpenAI SDK's multipart encoder picks up the real filename and content type
// instead of defaulting to "anonymous_file" / "application/octet-stream".
type namedReader struct {
	io.Reader
	filename    string
	contentType string
}

func (r *namedReader) Filename() string    { return r.filename }
func (r *namedReader) ContentType() string { return r.contentType }

// UploadFile stores a binary via the Files API and returns its id. PDFs upload
// as user_data (consumed by Responses input_file), images as vision
// (input_image). Chat references the id later.
func (c *openAIClient) UploadFile(ctx context.Context, reader io.Reader, filename, contentType string) (string, string, error) {
	purpose := openai.FilePurposeUserData
	if isImage(contentType) {
		purpose = openai.FilePurposeVision
	}
	file, err := c.client.Files.New(ctx, openai.FileNewParams{
		File:    &namedReader{Reader: reader, filename: filename, contentType: contentType},
		Purpose: purpose,
	})
	if err != nil {
		return "", "", err
	}
	return openAiProvider, file.ID, nil
}

// DeleteFile removes an uploaded file via the Files API.
func (c *openAIClient) DeleteFile(ctx context.Context, fileID string) error {
	_, err := c.client.Files.Delete(ctx, fileID)
	return err
}

func openAIFirstContent(resp *openai.ChatCompletion) (string, error) {
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}
