package llm

import (
	"bytes"
	"context"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/assert"

	"zeebee-backend/domains/zeebee/config"
)

func TestMockClientUploadReturnsProviderAndID(t *testing.T) {
	m := &MockClient{}
	provider, id, err := m.UploadFile(context.Background(), bytes.NewReader([]byte("pdf")), "a.pdf", "application/pdf")
	assert.NoError(t, err)
	assert.Equal(t, "mock", provider)
	assert.Equal(t, "file-1", id)

	// a second upload yields a distinct id (resume scenarios rely on uniqueness)
	_, id2, _ := m.UploadFile(context.Background(), bytes.NewReader([]byte("pdf")), "b.pdf", "application/pdf")
	assert.NotEqual(t, id, id2)
	assert.Equal(t, []string{"file-1", id2}, m.UploadedFiles)
}

func TestMockClientChatPassesAttachments(t *testing.T) {
	m := &MockClient{ChatContent: `{"k":"v"}`}
	got, err := m.Chat(context.Background(), "sys", "user",
		WithAttachment(Attachment{FileID: "file-1", ContentType: "application/pdf"}),
		WithAttachment(Attachment{Base64: "AAAA", ContentType: "image/png", Filename: "x.png"}),
	)
	assert.NoError(t, err)
	assert.Equal(t, `{"k":"v"}`, got)

	atts := m.Attachments()
	assert.Len(t, atts, 2)
	assert.Equal(t, "file-1", atts[0].FileID)
	assert.Equal(t, "AAAA", atts[1].Base64)

	// no-attachment Chat still works (plain text path) and records nothing.
	_, _ = m.Chat(context.Background(), "sys", "user")
	assert.Len(t, m.Attachments(), 2)
}

func TestMockClientDeleteRecorded(t *testing.T) {
	m := &MockClient{}
	assert.NoError(t, m.DeleteFile(context.Background(), "file-9"))
	assert.Equal(t, []string{"file-9"}, m.DeletedFiles)
}

func TestClaudeBetaBlockMapsByKind(t *testing.T) {
	// fileId document (pdf)
	b := claudeBetaBlock(Attachment{FileID: "f1", ContentType: "application/pdf"})
	assert.NotNil(t, b.OfDocument)
	// fileId image
	b = claudeBetaBlock(Attachment{FileID: "f2", ContentType: "image/png"})
	assert.NotNil(t, b.OfImage)
	// base64 pdf
	b = claudeBetaBlock(Attachment{Base64: "AA", ContentType: "application/pdf"})
	assert.NotNil(t, b.OfDocument)
	// base64 image
	b = claudeBetaBlock(Attachment{Base64: "AA", ContentType: "image/jpeg"})
	assert.NotNil(t, b.OfImage)
}

func TestOpenAIInputPartMapsByKind(t *testing.T) {
	// fileId document (pdf) -> input_file by id
	p := openaiInputPart(Attachment{FileID: "f1", ContentType: "application/pdf"})
	assert.NotNil(t, p.OfInputFile)
	// fileId image -> input_image by id
	p = openaiInputPart(Attachment{FileID: "f2", ContentType: "image/png"})
	assert.NotNil(t, p.OfInputImage)
	// base64 pdf -> input_file with inline data
	p = openaiInputPart(Attachment{Base64: "AA", ContentType: "application/pdf", Filename: "x.pdf"})
	assert.NotNil(t, p.OfInputFile)
	// base64 image -> input_image with data url
	p = openaiInputPart(Attachment{Base64: "AA", ContentType: "image/jpeg"})
	assert.NotNil(t, p.OfInputImage)
}

func TestNewClient(t *testing.T) {
	providers := []config.LlmProvider{
		{Name: "openai", ApiKey: "k", BaseUrl: "https://api.openai.com/v1"},
		{Name: "claude", ApiKey: "k", BaseUrl: "https://api.anthropic.com"},
	}

	c, err := NewClient(providers, config.LlmBizConf{Provider: "openai", Model: "gpt-4o"})
	assert.NoError(t, err)
	assert.NotNil(t, c)

	_, err = NewClient(providers, config.LlmBizConf{Provider: "claude", Model: "claude-3"})
	assert.NoError(t, err)

	// provider not in the configured list
	_, err = NewClient(providers, config.LlmBizConf{Provider: "gemini", Model: "x"})
	assert.Error(t, err)

	// empty provider list
	_, err = NewClient(nil, config.LlmBizConf{Provider: "openai", Model: "x"})
	assert.Error(t, err)
}

func TestOpenAIFirstContent(t *testing.T) {
	// no choices -> error
	_, err := openAIFirstContent(&openai.ChatCompletion{})
	assert.Error(t, err)

	resp := &openai.ChatCompletion{
		Choices: []openai.ChatCompletionChoice{{
			Message: openai.ChatCompletionMessage{Content: "hello"},
		}},
	}
	got, err := openAIFirstContent(resp)
	assert.NoError(t, err)
	assert.Equal(t, "hello", got)
}

func TestClaudeText_Empty(t *testing.T) {
	// no content blocks -> error
	_, err := claudeText(&anthropic.Message{})
	assert.Error(t, err)
}

func TestCleanJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain json", `{"a":1}`, `{"a":1}`},
		{"fenced json", "```json\n{\"a\":1}\n```", `{"a":1}`},
		{"fenced no lang", "```\n{\"a\":1}\n```", `{"a":1}`},
		{"trailing fence only", "```json\n{\"a\":1}", `{"a":1}`},
		{"trimmed whitespace", `  {"a":1}  `, `{"a":1}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, cleanJSON(c.in))
		})
	}
}
