//go:build manual

package llm

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"zeebee-backend/domains/zeebee/config"
)

// TestClaudeFileSmoke verifies the real claude beta Files API end-to-end:
// UploadFile → Chat([fileId]) → DeleteFile. This is the riskiest, least-testable
// path (beta file-id message reference). Run manually:
//
//	go test -tags=manual -run TestClaudeFileSmoke -v ./common/llm/
//
// ApiKey/BaseUrl/Model are read from app/borrower/etc/config-local.yaml.
func TestClaudeFileSmoke(t *testing.T) {
	key, base, model := loadBorrowerClaude(t)
	providers := []config.LlmProvider{{Name: "claude", ApiKey: key, BaseUrl: base}}
	client, err := NewClient(providers, config.LlmBizConf{
		Provider: "claude", Model: model, MaxTokens: 2048, TimeoutSeconds: 120,
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

// loadBorrowerClaude extracts the claude provider ApiKey/BaseUrl and the
// TransactionFileParseConf Model from the borrower local config.
func loadBorrowerClaude(t *testing.T) (key, base, model string) {
	t.Helper()
	b, err := os.ReadFile("/Users/wei/Projects/go/Zeebee-Dev/backend/app/borrower/etc/config-local.yaml")
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	inClaude, inParse := false, false
	for _, ln := range strings.Split(string(b), "\n") {
		s := strings.TrimSpace(ln)
		switch {
		case strings.HasPrefix(s, "- Name: claude"):
			inClaude = true
		case strings.HasPrefix(s, "- Name:"):
			inClaude = false
		case strings.HasPrefix(s, "TransactionFileParseConf:"):
			inParse = true
		}
		if inClaude {
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
	if key == "" || base == "" || model == "" {
		t.Fatalf("could not parse claude config (key=%q base=%q model=%q)", key, base, model)
	}
	return key, base, model
}

func kv(line, prefix string) (string, bool) {
	s := strings.TrimSpace(line)
	if !strings.HasPrefix(s, prefix) {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(s, prefix)), true
}

// minimalSmokePDF builds a valid minimal 1-page PDF the Files API can ingest.
func minimalSmokePDF() []byte {
	var b bytes.Buffer
	b.WriteString("%PDF-1.4\n")
	o1 := b.Len()
	b.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	o2 := b.Len()
	b.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")
	o3 := b.Len()
	b.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>\nendobj\n")
	xref := b.Len()
	b.WriteString("xref\n0 4\n0000000000 65535 f \n")
	fmt.Fprintf(&b, "%010d 00000 n \n", o1)
	fmt.Fprintf(&b, "%010d 00000 n \n", o2)
	fmt.Fprintf(&b, "%010d 00000 n \n", o3)
	fmt.Fprintf(&b, "trailer\n<< /Size 4 /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", xref)
	return b.Bytes()
}
