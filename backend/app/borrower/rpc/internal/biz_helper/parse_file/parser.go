// Package parse_file coordinates the staged async file-parse pipeline (escrow
// on_upload, rpa manual): dedup lock, Redis-backed progress/hold, and restart
// recovery. Per-parser workers live in package handler.
package parse_file

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/handler"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/llm"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// ErrParseInProgress is returned when a manual parse for the same
// (parser, document) is already in flight. Surfaced as gRPC AlreadyExists so
// the API layer maps it to HTTP 409.
var ErrParseInProgress = status.Error(codes.AlreadyExists, "parse in progress")

// RuleSig is the dedup-lock + reqmap key, scoped per parser+document. Exported
// so the progress-by-params query rebuilds the same key.
func RuleSig(documentId int64, parser constants.DocumentTemplateParser) string {
	return fmt.Sprintf("%s|%d", parser, documentId)
}

// Parser is the coordinator: it routes a request to the right Handler by
// (parser, source, trigger, rule), then owns the full run lifecycle — dedup
// lock, job/hold/reqmap, heartbeat, finalize, terminal Done, and progress
// persistence. Handlers are pure workers (extract + writeback) and report stage
// granularity back through a Reporter.
type Parser struct {
	repos    *repo.AllRepos
	progMgr  *progress.ProgressManager
	logger   logx.Logger
	nodeId   string
	handlers map[string]handler.Handler
}

// NewParser builds the coordinator and registers the built-in handlers.
func NewParser(repos *repo.AllRepos, llmClient llm.Client, progMgr *progress.ProgressManager, logger logx.Logger) *Parser {
	p := &Parser{
		repos:    repos,
		progMgr:  progMgr,
		logger:   logger,
		nodeId:   nodeName(),
		handlers: make(map[string]handler.Handler),
	}
	p.Register(constants.DocumentTemplateFileParseParserEscrowInfo, handler.NewEscrowHandler(repos, logger, llmClient))
	p.Register(constants.DocumentTemplateFileParseParserRPA, handler.NewRPAHandler(repos, logger, llmClient))
	return p
}

// Register maps a parser name to its handler.
func (p *Parser) Register(name constants.DocumentTemplateParser, h handler.Handler) {
	p.handlers[string(name)] = h
}

// Trigger is the dispatch origin.
type Trigger int

const (
	TriggerOnUpload Trigger = iota
	TriggerManual
)

// DispatchParam carries everything Dispatch needs to route + run a parse.
type DispatchParam struct {
	TransactionId  int64
	DocumentId     int64
	OrganizationId int64        // manual IDOR check; ignored for on_upload
	Files          []model.File // on_upload batch; nil for manual
	Parser         string       // explicit parser for manual; ignored for on_upload (rule-driven)
	Trigger        Trigger
}

// Dispatch is the single entry point for both on_upload and manual parses. The
// shared contract (FR-031): resolve the rule(s) by parser, then gate by the
// rule's declared TriggeredBy — a rule runs only via the trigger it declares
// (Trigger and Source are orthogonal; the gate looks only at TriggeredBy, the
// handler resolves file scope from Source). It then acquires the dedup lock,
// seeds job/hold/progress, and kicks the async runner. On_upload is fail-silent
// (a dispatch error is logged, not surfaced — the upload response must not be
// blocked by the side-channel parse); manual surfaces ErrParseInProgress on a
// duplicate.
func (p *Parser) Dispatch(ctx context.Context, param DispatchParam) (requestId string, err error) {
	if p == nil {
		return "", nil
	}
	defer func() {
		if err != nil {
			p.logger.WithContext(ctx).WithFields(logx.Field("param", param)).
				Errorf("dispatch parse: %s", err)
		}
	}()

	if param.Trigger == TriggerOnUpload {
		return p.dispatchOnUpload(ctx, param)
	}
	return p.dispatchManual(ctx, param)
}

// dispatchOnUpload reads the document's parse config by DocumentId (on_upload
// carries no parser name) and runs every Upload-triggered rule as its own async
// parse. Distinct parsers have distinct dedup sigs, so they don't collide.
// Fail-silent: dispatch errors are logged by the Dispatch defer, not surfaced to
// the upload response (the upload already succeeded).
func (p *Parser) dispatchOnUpload(ctx context.Context, param DispatchParam) (string, error) {
	doc, e := p.repos.DocumentRepo.FindOneWithFiles(ctx, param.DocumentId)
	if e != nil || doc == nil {
		return "", errorx.New("document not found")
	}
	rules := onUploadRules(docFileParseConfig(doc)) // gate: TriggeredBy == Upload
	var lastReq string
	for _, rule := range rules {
		h := p.handlers[string(rule.Parser)]
		if h == nil {
			continue
		}
		rid, e := p.start(ctx, param.TransactionId, param.DocumentId, rule, h, param.Files, true)
		if e != nil {
			return "", e
		}
		lastReq = rid
	}
	return lastReq, nil
}

// dispatchManual resolves the single rule for the explicit parser, then gates it
// by TriggeredBy == Manual. A parser whose only rule is Upload-triggered cannot
// be fired manually — the trigger a rule declares is the only trigger that runs
// it.
func (p *Parser) dispatchManual(ctx context.Context, param DispatchParam) (string, error) {
	doc, e := p.assertOwnedBranch(ctx, param.TransactionId, param.DocumentId, param.OrganizationId)
	if e != nil {
		return "", e
	}
	rule, ok := resolveRule(docFileParseConfig(doc), param.Parser) // by parser, trigger-agnostic
	if !ok {
		return "", errorx.New("no matching parse rule for document")
	}
	if !ruleMatchesTrigger(rule, TriggerManual) {
		return "", errorx.New("parse rule does not apply to manual trigger")
	}
	h := p.handlers[string(rule.Parser)]
	if h == nil {
		return "", errorx.New("unknown parser")
	}
	return p.start(ctx, param.TransactionId, param.DocumentId, rule, h, param.Files, false)
}

// start acquires the dedup lock (→ ErrParseInProgress on manual / join ongoing
// on on_upload), seeds job/reqmap/hold/progress, and kicks the runner.
func (p *Parser) start(ctx context.Context, txId, docId int64, rule dto.DocumentTemplateFileParseRule, h handler.Handler, files []model.File, onUpload bool) (string, error) {
	sig := RuleSig(docId, rule.Parser)
	requestId := uuid.NewString()
	ok, err := p.progMgr.AcquireDedupLock(ctx, txId, sig, requestId)
	if err != nil {
		return "", errorx.Wrap(err, "acquire parse lock")
	}
	if !ok {
		if onUpload {
			reqId, _ := p.progMgr.GetReqMap(ctx, txId, sig) // an active drainer owns this run
			return reqId, nil
		}
		return "", ErrParseInProgress
	}
	job := &progress.Job{
		RequestId: requestId, TransactionId: txId, RuleSig: sig,
		Parser: string(rule.Parser), Mode: string(rule.Source), DocumentId: docId,
	}
	if err := p.seed(ctx, requestId, job, sig, h.Stages(files)); err != nil {
		_ = p.progMgr.ReleaseDedupLock(ctx, txId, sig)
		return "", err
	}
	threading.GoSafe(func() {
		p.run(txId, docId, requestId, sig, h, rule, files)
	})
	return requestId, nil
}

// run drives one parse through the shared lifecycle (heartbeat + finalize +
// terminal Done). It hands the rule + in-scope files straight to the handler;
// the handler decides whether to use the passed files or collect its own scope
// from the DB (per rule.Source). Trigger and scope stay orthogonal here — this
// function does not branch on either.
func (p *Parser) run(txId, docId int64, requestId, sig string, h handler.Handler, rule dto.DocumentTemplateFileParseRule, files []model.File) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stopHeart := p.startHeartbeat(ctx, txId, requestId, sig)
	defer stopHeart()
	defer p.finalize(ctx, txId, requestId, sig)

	rep := p.reporter(requestId)
	// A recovered run may start after the sliding progress record (ttlProgressSec)
	// has expired: re-seed it so the handler's stage writes are visible to polls
	// again. Fresh runs skip this (progress already seeded by start()).
	if prog, _ := p.progMgr.GetProgress(ctx, requestId); prog == nil {
		_ = p.progMgr.InitProgress(ctx, requestId, h.Stages(files))
	}
	if err := h.Parse(ctx, txId, docId, files, rule, rep); err != nil {
		p.logger.WithContext(ctx).WithFields(
			logx.Field("request_id", requestId),
			logx.Field("stacks", errorx.GetErrorStack(err)),
		).Errorf("parse run: %s", err)
		_ = p.progMgr.WriteDone(ctx, requestId, progress.StatusFailed, err.Error())
		return
	}
	_ = p.progMgr.WriteDone(ctx, requestId, progress.StatusFinished, "")
}

// assertOwnedBranch verifies the transaction belongs to orgId, the document
// exists, is a Branch attached to that transaction, and returns the document
// for rule resolution.
func (p *Parser) assertOwnedBranch(ctx context.Context, txId, documentId, orgId int64) (*model.Document, error) {
	tx, err := p.repos.BorrowerPropertyTransactionRepo.GetOrganizationTransaction(ctx, txId, orgId)
	if err != nil || tx == nil {
		return nil, errorx.New("transaction not found")
	}
	doc, err := p.repos.DocumentRepo.FindOneWithFiles(ctx, documentId)
	if err != nil || doc == nil {
		return nil, errorx.New("document not found")
	}
	if doc.Type != constants.DocumentTypeBranch {
		return nil, errorx.New("document must be a branch folder")
	}
	entDoc, err := p.repos.EntityDocumentRepo.GetByEntity(ctx, constants.EntityDocumentTypeBorrowerPropertyTransaction, txId)
	if err != nil || entDoc == nil || entDoc.Id != doc.EntityDocumentId {
		return nil, errorx.New("document does not belong to transaction")
	}
	return doc, nil
}

// seed writes the recovery job, params→requestId map, hold lock, and initial
// staged progress so the first poll sees Running and a restart can reclaim.
func (p *Parser) seed(ctx context.Context, requestId string, job *progress.Job, sig string, stages []progress.StageState) error {
	if err := p.progMgr.WriteJob(ctx, job); err != nil {
		return errorx.Wrap(err, "write job")
	}
	if err := p.progMgr.SetReqMap(ctx, job.TransactionId, sig, requestId); err != nil {
		return errorx.Wrap(err, "write reqmap")
	}
	if _, err := p.progMgr.AcquireHold(ctx, requestId, p.nodeId); err != nil {
		return errorx.Wrap(err, "acquire hold")
	}
	return p.progMgr.InitProgress(ctx, requestId, stages)
}

// startHeartbeat renews hold + dedup lock until the returned stop func runs.
func (p *Parser) startHeartbeat(ctx context.Context, txId int64, requestId, sig string) (stop func()) {
	ticker := time.NewTicker(progress.HoldHeartbeat / 2)
	done := make(chan struct{})
	threading.GoSafe(func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			case <-ticker.C:
				_ = p.progMgr.RenewHold(ctx, requestId)
				_ = p.progMgr.RenewDedupLock(ctx, txId, sig)
			}
		}
	})
	return func() { close(done); ticker.Stop() }
}

func (p *Parser) finalize(ctx context.Context, txId int64, requestId, sig string) {
	_ = p.progMgr.ReleaseHold(ctx, requestId)
	_ = p.progMgr.ReleaseDedupLock(ctx, txId, sig)
	_ = p.progMgr.DeleteJob(ctx, requestId)
}

// reporter wraps the ProgressManager into the narrow handler.Reporter interface
// so handlers report stage granularity without touching Redis directly.
func (p *Parser) reporter(requestId string) handler.Reporter {
	return redisReporter{progMgr: p.progMgr, requestId: requestId}
}

type redisReporter struct {
	progMgr   *progress.ProgressManager
	requestId string
}

func (r redisReporter) MarkStage(ctx context.Context, stage, status string) {
	prog, err := r.progMgr.GetProgress(ctx, r.requestId)
	if err != nil || prog == nil {
		return
	}
	if stage != "" {
		prog.Stage = stage
	}
	for i := range prog.Stages {
		if status != "" && prog.Stages[i].Name == stage {
			prog.Stages[i].Status = status
		}
	}
	_ = r.progMgr.SaveProgress(ctx, r.requestId, prog)
}

func (r redisReporter) SetStageTotal(ctx context.Context, stage string, total int) {
	r.mut(ctx, stage, func(s *progress.StageState) { s.Total = total })
}

func (r redisReporter) SetStageProcessed(ctx context.Context, stage string, processed int) {
	r.mut(ctx, stage, func(s *progress.StageState) { s.Processed = processed })
}

// StageOutput returns the stage's marshaled resume checkpoint (extract envelopes)
// from its own job-TTL key so a recovered run recalls already-extracted results
// even after the sliding progress record has expired.
func (r redisReporter) StageOutput(ctx context.Context, stage string) json.RawMessage {
	out, _ := r.progMgr.GetStageCheckpoint(ctx, r.requestId, stage)
	return out
}

// SaveStageOutput writes the stage checkpoint to its own job-TTL key.
func (r redisReporter) SaveStageOutput(ctx context.Context, stage string, output json.RawMessage) {
	_ = r.progMgr.SetStageCheckpoint(ctx, r.requestId, stage, output)
}

func (r redisReporter) mut(ctx context.Context, stage string, fn func(*progress.StageState)) {
	prog, err := r.progMgr.GetProgress(ctx, r.requestId)
	if err != nil || prog == nil {
		return
	}
	for i := range prog.Stages {
		if prog.Stages[i].Name == stage {
			fn(&prog.Stages[i])
		}
	}
	_ = r.progMgr.SaveProgress(ctx, r.requestId, prog)
}

var nodeNameOnce sync.Once
var nodeNameCache string

// nodeName identifies this process for hold-lock ownership; lazily computed.
func nodeName() string {
	nodeNameOnce.Do(func() {
		nodeNameCache = uuid.New().String()
	})
	return nodeNameCache
}
