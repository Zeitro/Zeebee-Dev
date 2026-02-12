package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const SourceTypeText = "Text"
const SourceTypeUrl = "Url"

type GeneratePdfRequest struct {
	Type string `json:"type"`
	Url  string `json:"url"`
	Text string `json:"text"`
}

type GeneratePdfResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Data []byte `json:"data"`
	} `json:"data"`
}

type Client struct {
	url    string
	client *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
		client: &http.Client{
			Timeout: 100 * time.Second,
		},
	}
}

// GeneratePdf generates pdf from html or url
func (c *Client) GeneratePdf(ctx context.Context, payload GeneratePdfRequest) ([]byte, error) {

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.url,
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PDF service returned status %d", resp.StatusCode)
	}

	var generatePdfResponse GeneratePdfResponse
	err = json.NewDecoder(resp.Body).Decode(&generatePdfResponse)
	if err != nil {
		return nil, err
	}
	if generatePdfResponse.Code != http.StatusOK {
		return nil, fmt.Errorf(generatePdfResponse.Message)
	}
	if generatePdfResponse.Data.Data == nil {
		return nil, fmt.Errorf("PDF service returned empty data")
	}
	return generatePdfResponse.Data.Data, nil
}
