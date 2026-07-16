package biz_helper

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	anthopt "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
	"github.com/openai/openai-go/v3"
	openaiopt "github.com/openai/openai-go/v3/option"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/config"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

const (
	extractSystemPrompt = "You extract data from documents. Your entire response MUST be a single raw JSON object and nothing else: start with { and end with }. Do NOT wrap it in markdown code fences (no ```json, no ```). Do NOT add prose, explanations, backticks, or any text before or after the JSON."
	// maxPDFBytes caps a downloaded PDF at 10 MB.
	maxPDFBytes = 10 * 1024 * 1024
	// pdfPoolSize caps concurrent pdfium parses; extra callers block on GetInstance.
	// ponytail: raise if upload throughput demands it.
	pdfPoolSize = 5
)

// ============================================================
// LLM clients (openai / claude). Kept local until a second
// consumer justifies promoting to its own package.
// ============================================================

// chatClient abstracts the LLM used for file content extraction so openai and
// claude impls are interchangeable.
type chatClient interface {
	// ChatJSON sends a text prompt and returns raw model content (expected JSON).
	ChatJSON(ctx context.Context, system, user string) (string, error)
	// VisionJSON sends a text prompt plus one PNG image and returns raw content.
	VisionJSON(ctx context.Context, system, text string, imagePNG []byte) (string, error)
}

// NewChatClient resolves the provider named by biz.Provider from the configured
// provider list. Unknown providers fail fast.
func NewChatClient(providers []config.LlmProvider, biz config.LlmBizConf) (chatClient, error) {
	var prov *config.LlmProvider
	for i := range providers {
		if providers[i].Name == biz.Provider {
			prov = &providers[i]
			break
		}
	}
	if prov == nil {
		return nil, fmt.Errorf("llm provider %q not configured", biz.Provider)
	}
	switch prov.Name {
	case "openai":
		return newOpenAIClient(prov, biz.Model), nil
	case "claude":
		return newClaudeClient(prov, biz.Model), nil
	default:
		return nil, nil
	}
}

// --- openai ---

type openAIClient struct {
	client openai.Client
	model  openai.ChatModel
}

func newOpenAIClient(p *config.LlmProvider, model string) *openAIClient {
	opts := []openaiopt.RequestOption{
		openaiopt.WithAPIKey(p.ApiKey),
		openaiopt.WithRequestTimeout(60 * time.Second),
	}
	// openai-go appends "/chat/completions" to base_url, so base must carry "/v1".
	if base := strings.TrimRight(p.BaseUrl, "/"); base != "" {
		if !strings.HasSuffix(base, "/v1") {
			base += "/v1"
		}
		opts = append(opts, openaiopt.WithBaseURL(base))
	}
	return &openAIClient{client: openai.NewClient(opts...), model: model}
}

func (c *openAIClient) ChatJSON(ctx context.Context, system, user string) (string, error) {
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

func (c *openAIClient) VisionJSON(ctx context.Context, system, text string, imagePNG []byte) (string, error) {
	dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(imagePNG)
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(system),
			openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
				openai.TextContentPart(text),
				openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
					URL:    dataURL,
					Detail: "high",
				}),
			}),
		},
	})
	if err != nil {
		return "", err
	}
	return openAIFirstContent(resp)
}

func openAIFirstContent(resp *openai.ChatCompletion) (string, error) {
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}

// --- claude ---

type claudeClient struct {
	client    anthropic.Client
	model     anthropic.Model
	maxTokens int64
}

func newClaudeClient(p *config.LlmProvider, model string) *claudeClient {
	opts := []anthopt.RequestOption{
		anthopt.WithAPIKey(p.ApiKey),
		anthopt.WithRequestTimeout(60 * time.Second),
	}
	if base := strings.TrimRight(p.BaseUrl, "/"); base != "" {
		opts = append(opts, anthopt.WithBaseURL(base))
	}

	return &claudeClient{
		client:    anthropic.NewClient(opts...),
		model:     model,
		maxTokens: 4096,
	}
}

func (c *claudeClient) ChatJSON(ctx context.Context, system, user string) (string, error) {
	resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		System:    []anthropic.TextBlockParam{{Text: system}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(user)),
		},
	})
	if err != nil {
		return "", err
	}
	return claudeText(resp)
}

func (c *claudeClient) VisionJSON(ctx context.Context, system, text string, imagePNG []byte) (string, error) {
	resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		System:    []anthropic.TextBlockParam{{Text: system}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(text),
				anthropic.NewImageBlockBase64("image/png", base64.StdEncoding.EncodeToString(imagePNG)),
			),
		},
	})
	if err != nil {
		return "", err
	}
	return claudeText(resp)
}

func claudeText(resp *anthropic.Message) (string, error) {
	var sb strings.Builder
	for _, block := range resp.Content {
		if t := block.AsText(); t.Text != "" {
			sb.WriteString(t.Text)
		}
	}
	if sb.Len() == 0 {
		return "", fmt.Errorf("claude returned no text")
	}
	return sb.String(), nil
}

// ============================================================
// PDF helpers: pdfium pool, download, render, text.
// ============================================================

var (
	pdfPoolOnce sync.Once
	pdfPool     pdfium.Pool
	pdfPoolErr  error
)

func getPdfPool() (pdfium.Pool, error) {
	pdfPoolOnce.Do(func() {
		// ponytail: wazero pure-Go runtime, no CGO, no pdfium binary install.
		pdfPool, pdfPoolErr = webassembly.Init(webassembly.Config{
			MinIdle:  0,
			MaxIdle:  2,
			MaxTotal: pdfPoolSize,
		})
	})
	return pdfPool, pdfPoolErr
}

func withPdfInstance(fn func(instance pdfium.Pdfium) error) error {
	pool, err := getPdfPool()
	if err != nil {
		return err
	}
	inst, err := pool.GetInstance(60 * time.Second)
	if err != nil {
		return err
	}
	defer func() { _ = inst.Close() }()
	return fn(inst)
}

// downloadPDF fetches a presigned file URL (≤10 MB) and returns a seekable
// reader for pdfium. The body is buffered into memory because a streaming http
// body is not seekable and pdfium seeks the xref ranges itself.
func downloadPDF(ctx context.Context, url string) (io.ReadSeeker, error) {
	ctx, _ = context.WithTimeout(ctx, 60*time.Second)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("download http %d", res.StatusCode)
	}
	if res.ContentLength > maxPDFBytes {
		return nil, fmt.Errorf("pdf too large: %d bytes", res.ContentLength)
	}
	// LimitReader +1 so an undersized Content-Length still trips the cap.
	b, err := io.ReadAll(io.LimitReader(res.Body, maxPDFBytes+1))
	if err != nil {
		return nil, err
	}
	if len(b) > maxPDFBytes {
		return nil, fmt.Errorf("pdf exceeds 10MB")
	}
	return bytes.NewReader(b), nil
}

// openPDF opens a seekable reader as a pdfium document. Size is derived via
// Seek since io.ReadSeeker has no Len(); the reader is reset to the start so
// pdfium reads from the beginning.
func openPDF(instance pdfium.Pdfium, r io.ReadSeeker) (references.FPDF_DOCUMENT, func(), error) {
	size, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return "", nil, err
	}
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return "", nil, err
	}
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		FileReader:     r,
		FileReaderSize: size,
	})
	if err != nil {
		return "", nil, err
	}
	return doc.Document, func() {
		_, _ = instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})
	}, nil
}

func readPageText(instance pdfium.Pdfium, doc references.FPDF_DOCUMENT, page int) (string, error) {
	tp, err := instance.FPDFText_LoadPage(&requests.FPDFText_LoadPage{
		Page: requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: page}},
	})
	if err != nil {
		return "", err
	}
	defer func() { _, _ = instance.FPDFText_ClosePage(&requests.FPDFText_ClosePage{TextPage: tp.TextPage}) }()

	count, err := instance.FPDFText_CountChars(&requests.FPDFText_CountChars{TextPage: tp.TextPage})
	if err != nil {
		return "", err
	}
	if count.Count <= 0 {
		return "", nil
	}
	res, err := instance.FPDFText_GetText(&requests.FPDFText_GetText{
		TextPage:   tp.TextPage,
		StartIndex: 0,
		Count:      count.Count,
	})
	if err != nil {
		return "", err
	}
	return res.Text, nil
}

func renderPagePNG(instance pdfium.Pdfium, doc references.FPDF_DOCUMENT, page int) ([]byte, error) {
	res, err := instance.RenderToFile(&requests.RenderToFile{
		RenderPageInDPI: &requests.RenderPageInDPI{
			Page:       requests.Page{ByIndex: &requests.PageByIndex{Document: doc, Index: page}},
			DPI:        150,
			RenderForm: true,
		},
		OutputFormat: requests.RenderToFileOutputFormatPNG,
		OutputTarget: requests.RenderToFileOutputTargetBytes,
	})
	if err != nil {
		return nil, err
	}
	if res.ImageBytes == nil {
		return nil, fmt.Errorf("render returned no image bytes")
	}
	return *res.ImageBytes, nil
}

// ============================================================
// Extractors + orchestrator
// ============================================================

type escrowFields struct {
	EscrowCompany string `json:"escrow_company"`
	OfficerName   string `json:"officer_name"`
	OfficerEmail  string `json:"officer_email"`
	OfficerPhone  string `json:"officer_phone"`
	EscrowNumber  string `json:"escrow_number"`
}

func (e escrowFields) complete() bool {
	return e.EscrowCompany != "" && e.OfficerName != "" && e.OfficerEmail != "" &&
		e.OfficerPhone != "" && e.EscrowNumber != ""
}

// parsePurchasePrice extracts the float from the model's JSON response.
func parsePurchasePrice(content string) (float64, error) {
	var probe struct {
		PurchasePrice float64 `json:"purchase_price"`
	}
	if err := json.Unmarshal([]byte(content), &probe); err != nil {
		return 0, fmt.Errorf("rpa response not json: %s", content)
	}
	return probe.PurchasePrice, nil
}

// parseEscrowFields extracts escrow fields from the model's JSON response.
func parseEscrowFields(content string) (escrowFields, error) {
	var f escrowFields
	if err := json.Unmarshal([]byte(content), &f); err != nil {
		return f, fmt.Errorf("escrow response not json: %s", content)
	}
	return f, nil
}

// extractPurchasePrice pulls purchase_price (float) from page-1 text. 0 = absent.
func extractPurchasePrice(ctx context.Context, llm chatClient, r io.ReadSeeker) (price float64, err error) {
	err = withPdfInstance(func(inst pdfium.Pdfium) error {
		doc, cleanup, e := openPDF(inst, r)
		if e != nil {
			return e
		}
		defer cleanup()
		text, e := readPageText(inst, doc, 0)
		if e != nil {
			return e
		}
		content, e := llm.ChatJSON(ctx, extractSystemPrompt, rpaPrompt(text))
		if e != nil {
			return e
		}
		price, e = parsePurchasePrice(content)
		return e
	})
	return price, err
}

// extractEscrowFields vision-reads escrow fields off the rendered first page.
func extractEscrowFields(ctx context.Context, llm chatClient, r io.ReadSeeker) (fields escrowFields, err error) {
	err = withPdfInstance(func(inst pdfium.Pdfium) error {
		doc, cleanup, e := openPDF(inst, r)
		if e != nil {
			return e
		}
		defer cleanup()
		png, e := renderPagePNG(inst, doc, 0)
		if e != nil {
			return e
		}
		content, e := llm.VisionJSON(ctx, extractSystemPrompt, escrowPrompt(), png)
		if e != nil {
			return e
		}
		fields, e = parseEscrowFields(content)
		return e
	})
	return fields, err
}

// PropertyTransactionFileParser parses uploaded PDFs and writes extracted
// content back to the transaction/escrow. All work happens AFTER the upload tx
// commits, in a separate tx; failures are returned for the caller to log but do
// not fail the upload. Only documents carrying a file_parse_config AND pdf files
// are parsed; multi-file: first non-empty value per field wins, stop once all
// fields are filled.
type PropertyTransactionFileParser struct {
	repos  *repo.AllRepos
	llm    chatClient
	logger logx.Logger
}

func NewPropertyTransactionFileParser(repos *repo.AllRepos, llm chatClient, logger logx.Logger) *PropertyTransactionFileParser {
	return &PropertyTransactionFileParser{
		repos:  repos,
		llm:    llm,
		logger: logger,
	}
}

// ParseAndUpdateTransactionFile dispatches on the document's file_parse_config.
// transactionId is the borrower_property_transaction id for write-back.
func (p *PropertyTransactionFileParser) ParseAndUpdateTransactionFile(
	ctx context.Context,
	transactionId int64,
	doc *model.Document,
	files []model.File,
) error {
	parserName := readParserName(doc)
	switch parserName {
	case "rpa":
		return p.parseRPA(ctx, transactionId, files)
	case "escrow_info":
		return p.parseEscrow(ctx, transactionId, files)
	default:
		return nil
	}
}

func (p *PropertyTransactionFileParser) parseRPA(ctx context.Context, transactionId int64, files []model.File) error {
	var price float64
	for _, f := range files {
		if price > 0 {
			break // first valid price wins
		}
		if !isPDF(f.FileName) {
			p.logger.WithContext(ctx).Infof("not a PDF: %s", f.FileName)
			continue
		}
		data, err := downloadPDF(ctx, f.FileUrl)
		if err != nil {
			p.logger.WithContext(ctx).
				WithFields(
					logx.Field("file", f.FileName),
					logx.Field("stacks", errorx.GetErrorStack(err)),
				).
				Errorf("downloadPDF: %s", err)
			continue
		}
		got, err := extractPurchasePrice(ctx, p.llm, data)
		if err != nil {
			p.logger.WithContext(ctx).
				WithFields(
					logx.Field("file", f.FileName),
					logx.Field("stacks", errorx.GetErrorStack(err)),
				).
				Errorf("extractPurchasePrice: %s", err)
			continue
		}
		if got > 0 {
			price = got
		}
	}
	if price <= 0 {
		return nil // nothing extracted; keep existing
	}
	return p.repos.BorrowerPropertyTransactionRepo.Transaction(ctx, func(tx *gorm.DB) error {
		return p.repos.BorrowerPropertyTransactionRepo.UpdateColumnsByIdWithTx(ctx, tx, transactionId, map[string]interface{}{
			"purchase_amount": price,
		})
	})
}

func (p *PropertyTransactionFileParser) parseEscrow(ctx context.Context, transactionId int64, files []model.File) error {
	var merged escrowFields
	for _, f := range files {
		if merged.complete() {
			break
		}
		if !isPDF(f.FileName) {
			continue
		}
		data, err := downloadPDF(ctx, f.FileUrl)
		if err != nil {
			p.logger.WithContext(ctx).
				WithFields(
					logx.Field("file", f.FileName),
					logx.Field("stacks", errorx.GetErrorStack(err)),
				).
				Errorf("downloadPDF: %s", err)
			continue
		}
		got, err := extractEscrowFields(ctx, p.llm, data)
		if err != nil {
			p.logger.WithContext(ctx).
				WithFields(
					logx.Field("file", f.FileName),
					logx.Field("stacks", errorx.GetErrorStack(err)),
				).
				Errorf("extractEscrowFields: %s", err)
			continue
		}
		// first non-empty value per field wins
		if merged.EscrowCompany == "" {
			merged.EscrowCompany = got.EscrowCompany
		}
		if merged.OfficerName == "" {
			merged.OfficerName = got.OfficerName
		}
		if merged.OfficerEmail == "" {
			merged.OfficerEmail = got.OfficerEmail
		}
		if merged.OfficerPhone == "" {
			merged.OfficerPhone = got.OfficerPhone
		}
		if merged.EscrowNumber == "" {
			merged.EscrowNumber = got.EscrowNumber
		}
	}

	if merged == (escrowFields{}) {
		return nil // nothing extracted
	}

	return p.repos.PropertyTransactionEscrowRepo.Transaction(ctx, func(tx *gorm.DB) error {
		rows, err := p.repos.PropertyTransactionEscrowRepo.FindAllWithTx(ctx, tx, map[string]interface{}{
			"borrower_property_transaction_id": transactionId,
		})
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			esc := &model.PropertyTransactionEscrow{
				BorrowerPropertyTransactionId: transactionId,
				Checklist:                     dto.PropertyTransactionEscrowChecklist{},
			}
			applyEscrowFields(esc, merged)
			_, err = p.repos.PropertyTransactionEscrowRepo.CreateWithTx(ctx, tx, esc, false)
			return err
		}
		// update the oldest row (min id) as the primary escrow
		esc := rows[0]
		for _, r := range rows[1:] {
			if r.Id < esc.Id {
				esc = r
			}
		}
		applyEscrowFields(esc, merged)
		_, err = p.repos.PropertyTransactionEscrowRepo.UpdateOneByIdWithTx(ctx, tx, esc.Id, esc)
		return err
	})
}

// applyEscrowFields copies only non-empty extracted fields (skip-empty merge).
func applyEscrowFields(esc *model.PropertyTransactionEscrow, f escrowFields) {
	if f.EscrowCompany != "" {
		esc.EscrowCompany = f.EscrowCompany
	}
	if f.OfficerName != "" {
		esc.OfficerName = f.OfficerName
	}
	if f.OfficerEmail != "" {
		esc.OfficerEmail = f.OfficerEmail
	}
	if f.OfficerPhone != "" {
		esc.OfficerPhone = f.OfficerPhone
	}
	if f.EscrowNumber != "" {
		esc.No = f.EscrowNumber
	}
}

func readParserName(doc *model.Document) string {
	if len(doc.Metadata) == 0 {
		return ""
	}
	var meta struct {
		FileParseConfig *struct {
			Parser string `json:"parser"`
		} `json:"file_parse_config"`
	}
	if err := json.Unmarshal(doc.Metadata, &meta); err != nil {
		return ""
	}
	if meta.FileParseConfig == nil {
		return ""
	}
	return meta.FileParseConfig.Parser
}

// isPDF gates parsing to PDF files only; non-PDF uploads to a parse-enabled doc
// are skipped.
func isPDF(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".pdf")
}

func rpaPrompt(text string) string {
	return fmt.Sprintf(`Extract the purchase price from the following document text.
"purchase_price" MUST be a JSON number (float), not a string. Strip thousands separators (commas), currency symbols, and whitespace. Examples: "1,750,000" -> 1750000, "$1,750,000.00" -> 1750000.
Respond as JSON: {"purchase_price": <number>}
If not found, respond: {"purchase_price": 0}

Document text:
%s`, text)
}

func escrowPrompt() string {
	return `Extract these fields from the attached first page of a wire/escrow instruction document.
- escrow_company: the escrow/title company name in standard Title Case (e.g. "Mission Escrow, Inc.").
- officer_name: the escrow officer's full name only (e.g. "Renee Lo").
- officer_email: the escrow officer's email address, if printed on the page.
- officer_phone: the escrow officer's phone number as DIGITS ONLY — strip spaces, parentheses, dashes, dots, plus signs (e.g. "(408) 513-6005" becomes "4085136005").
- escrow_number: the value next to the "Escrow No." label (e.g. "202656-RL").
Only extract what is actually printed on the page. If a field is absent, unreadable, or uncertain, use null — never guess.
Respond as JSON: {"escrow_company": <string|null>, "officer_name": <string|null>, "officer_email": <string|null>, "officer_phone": <string|null>, "escrow_number": <string|null>}`
}
