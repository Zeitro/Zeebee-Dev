package llm

import (
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/assert"

	"zeebee-backend/domains/zeebee/config"
)

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
