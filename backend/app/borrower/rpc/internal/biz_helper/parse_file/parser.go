package parse_file

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// Parser orchestrates async file extraction. It holds shared deps and
// dispatches to per-parser handlers (escrow.go, rpa.go) so this file stays the
// thin entry point and each parser's logic lives in its own file.
//
// US1 wires the escrow OnUpload path; ManualTree/ManualSelect (rpa 5-field +
// arbitrate) arrive with US2/US3.
type Parser struct {
	repos  *repo.AllRepos
	llm    llm.Client
	prog   *ProgressManager
	logger logx.Logger
	nodeId string
}

func NewParser(repos *repo.AllRepos, llmClient llm.Client, prog *ProgressManager, logger logx.Logger) *Parser {
	return &Parser{
		repos:  repos,
		llm:    llmClient,
		prog:   prog,
		logger: logger,
		nodeId: nodeName(),
	}
}

// DispatchOnUpload is fire-and-forget: for each OnUpload rule on the document
// it routes the uploaded PDFs to the matching parser handler. It never blocks
// the upload response and never returns an error to the caller (fail-silent,
// FR-019).
func (p *Parser) DispatchOnUpload(ctx context.Context, txId int64, doc *model.Document, files []model.File) {
	if p == nil || p.llm == nil {
		return // parse disabled
	}
	rules := onUploadRules(docFileParseConfig(doc))
	if len(rules) == 0 {
		return
	}
	var pdfFiles []model.File
	for _, f := range files {
		if isPDF(f.FileName) {
			pdfFiles = append(pdfFiles, f)
		}
	}
	if len(pdfFiles) == 0 {
		return
	}
	for _, rule := range rules {
		switch rule.Parser {
		case model.FileParseParserEscrowInfo:
			p.dispatchEscrowOnUpload(ctx, txId, doc.Id, pdfFiles)
		// case model.FileParseParserRpa: rpa is ManualTree only (US2) — never OnUpload.
		case model.FileParseParserRpa:
			// todo handle rpa
		default:
			continue
		}
	}
}

// runBackground kicks fn in a recovered goroutine (go-zero threading.GoSafe)
// so a panic in one parse path never takes down the rpc process.
func (p *Parser) runBackground(_ context.Context, fn func()) {
	threading.GoSafe(fn)
}

var nodeNameOnce sync.Once
var nodeNameCache string

// nodeName identifies this process for hold-lock ownership; lazily computed.
func nodeName() string {
	nodeNameOnce.Do(func() {
		host, _ := os.Hostname()
		nodeNameCache = host + ":" + strconv.Itoa(os.Getpid())
	})
	return nodeNameCache
}
