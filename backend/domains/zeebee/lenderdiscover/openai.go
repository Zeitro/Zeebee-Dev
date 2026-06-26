package lenderdiscover

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const openAIChatURL = "https://api.openai.com/v1/chat/completions"

type openAIChatRequest struct {
	Model          string              `json:"model"`
	Messages       []openAIChatMessage `json:"messages"`
	ResponseFormat *struct {
		Type string `json:"type"`
	} `json:"response_format,omitempty"`
	Temperature float64 `json:"temperature"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// CallOpenAIJSON asks OpenAI for a single JSON object (caller defines schema in prompts).
func CallOpenAIJSON(ctx context.Context, apiKey, model, systemPrompt, userPrompt string) (content string, err error) {
	if len(apiKey) == 0 {
		return "", fmt.Errorf("OPENAI_API_KEY is not set")
	}
	body := openAIChatRequest{
		Model: model,
		Messages: []openAIChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		ResponseFormat: &struct {
			Type string `json:"type"`
		}{Type: "json_object"},
		Temperature: 0.2,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	// Detach from the caller's deadline so go-zero RPC timeout doesn't cancel the
	// long-running OpenAI call; the http.Client.Timeout governs the overall limit.
	reqCtx := context.WithoutCancel(ctx)
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, openAIChatURL, bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 8 * time.Minute}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var parsed openAIChatResponse
	if err := json.Unmarshal(b, &parsed); err != nil {
		return "", fmt.Errorf("openai response decode: %w; body=%s", err, string(b))
	}
	if parsed.Error != nil && len(parsed.Error.Message) > 0 {
		return "", fmt.Errorf("openai api error: %s", parsed.Error.Message)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("openai http %d: %s", res.StatusCode, string(b))
	}
	if len(parsed.Choices) == 0 || len(parsed.Choices[0].Message.Content) == 0 {
		return "", fmt.Errorf("openai empty choices")
	}
	return parsed.Choices[0].Message.Content, nil
}
