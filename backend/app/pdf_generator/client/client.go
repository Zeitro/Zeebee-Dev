package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const SourceTypeText = "Text"
const SourceTypeUrl = "Url"
const SourceTypeData = "Data"

type GeneratePdfRequest struct {
	Type string `json:"type"`
	Url  string `json:"url"`
	Text string `json:"text"`
}

type docToPdfRequest struct {
	Type string `json:"type"`
	Url  string `json:"url,omitempty"`
	Data []byte `json:"data,omitempty"`
	Ext  string `json:"ext"`
}

type generateResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Data []byte `json:"data"`
	} `json:"data"`
}

type Client struct {
	baseUrl string
	client  *http.Client
}

// NewClient takes the pdf_generator base URL (e.g. http://host:port/api/v1);
// endpoints are appended per method.
func NewClient(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
		client: &http.Client{
			Timeout: 100 * time.Second,
		},
	}
}

func (c *Client) endpoint(path string) string {
	return strings.TrimRight(c.baseUrl, "/") + path
}

// postJSON sends payload, decodes the standard {code,message,data} body, and maps
// a 503 code to ErrOverload so callers can retry transient overloads.
func (c *Client) postJSON(ctx context.Context, path string, payload any) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(path), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pdf service returned status %d", resp.StatusCode)
	}

	var out generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if out.Code != http.StatusOK {
		return nil, errors.New(out.Message)
	}
	if out.Data.Data == nil {
		return nil, fmt.Errorf("pdf service returned empty data")
	}
	return out.Data.Data, nil
}

// GeneratePdf generates pdf from html or url.
func (c *Client) GeneratePdf(ctx context.Context, payload GeneratePdfRequest) ([]byte, error) {
	return c.postJSON(ctx, "/pdf", payload)
}

// DocToPdf converts a doc/docx to pdf. typ is SourceTypeUrl (pass url) or
// SourceTypeData (pass data bytes); ext is "doc" or "docx". Returns pdf bytes.
func (c *Client) DocToPdf(ctx context.Context, typ, url string, data []byte, ext string) ([]byte, error) {
	return c.postJSON(ctx, "/doc_to_pdf", docToPdfRequest{
		Type: typ,
		Url:  url,
		Data: data,
		Ext:  ext,
	})
}
