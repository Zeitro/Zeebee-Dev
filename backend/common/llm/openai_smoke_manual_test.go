//go:build manual

package llm

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"zeebee-backend/domains/zeebee/config"
)

// TestOpenAIFileSmoke verifies the real openai Files + Responses API end-to-end:
// UploadFile → Chat([fileId]) → DeleteFile. PDFs need the Responses API (Chat
// Completions can't carry files). Run manually:
//
//	go test -tags=manual -run TestOpenAIFileSmoke -v ./common/llm/
//
// ApiKey/BaseUrl/Model are read from app/borrower/etc/config-local.yaml. The
// configured openai endpoint MUST implement /v1/responses (a Chat-Completions-
// only proxy will 404 here).
func TestOpenAIFileSmoke(t *testing.T) {
	key, base, model := loadBorrowerOpenAI(t)
	providers := []config.LlmProvider{{Name: "openai", ApiKey: key, BaseUrl: base}}
	client, err := NewClient(providers, config.LlmBizConf{
		Provider: "openai", Model: model, MaxTokens: 2048, TimeoutSeconds: 120,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	pdf := minimalSmokePDF()
	provider, fileID, err := client.UploadFile(ctx, bytes.NewReader(pdf), "smoke.pdf", "application/pdf")
	if err != nil {
		t.Fatalf("UploadFile: %v", err)
	}
	t.Logf("uploaded: provider=%s fileID=%s", provider, fileID)
	defer func() {
		if err := client.DeleteFile(context.Background(), fileID); err != nil {
			t.Logf("DeleteFile (best-effort): %v", err)
		} else {
			t.Logf("deleted: %s", fileID)
		}
	}()

	resp, err := client.Chat(ctx,
		"You answer about a document. Reply with a single raw JSON object, nothing else.",
		"Reply with JSON: {\"pages_expected\":1,\"read_ok\":true} — confirm you can read the attached file.",
		WithAttachment(Attachment{FileID: fileID, ContentType: "application/pdf"}))
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	t.Logf("chat response: %s", resp)
	if !strings.HasPrefix(strings.TrimSpace(resp), "{") {
		t.Fatalf("response not JSON: %s", resp)
	}
}

// loadBorrowerOpenAI extracts the openai provider ApiKey/BaseUrl and the
// TransactionFileParseConf Model from the borrower local config. BaseUrl is
// optional (defaults to api.openai.com when empty).
func loadBorrowerOpenAI(t *testing.T) (key, base, model string) {
	t.Helper()
	b, err := os.ReadFile("/Users/wei/Projects/go/Zeebee-Dev/backend/app/borrower/etc/config-local.yaml")
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	inOpenAI, inParse := false, false
	for _, ln := range strings.Split(string(b), "\n") {
		s := strings.TrimSpace(ln)
		switch {
		case strings.HasPrefix(s, "- Name: openai"):
			inOpenAI = true
		case strings.HasPrefix(s, "- Name:"):
			inOpenAI = false
		case strings.HasPrefix(s, "TransactionFileParseConf:"):
			inParse = true
		}
		if inOpenAI {
			if v, ok := kv(s, "ApiKey:"); ok {
				key = v
			}
			if v, ok := kv(s, "BaseUrl:"); ok {
				base = v
			}
		}
		if inParse {
			if v, ok := kv(s, "Model:"); ok {
				model = v
			}
		}
	}
	if key == "" || model == "" {
		t.Fatalf("could not parse openai config (key=%q base=%q model=%q)", key, base, model)
	}
	return key, base, model
}
