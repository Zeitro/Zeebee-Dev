package parse_file

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/domains/zeebee/model"
)

// recoverSweepInterval is how often the background sweeper re-scans for
// abandoned jobs. Short enough that a crashed node's work is reclaimed quickly,
// long enough not to load Redis on an idle cluster.
const recoverSweepInterval = 30 * time.Second

// StartRecovery launches the background coordinator: an immediate scan on
// startup (so a restart resumes in-flight parses — including a half-drained
// escrow on_upload queue) plus a periodic sweep that reclaims work from crashed
// nodes. It runs for the life of the process.
func (p *Parser) StartRecovery() {
	if p == nil {
		return
	}
	threading.GoSafe(func() {
		ctx := context.Background()
		p.recoverJobs(ctx)
		ticker := time.NewTicker(recoverSweepInterval)
		defer ticker.Stop()
		for range ticker.C {
			p.recoverJobs(ctx)
		}
	})
}

// recoverJobs scans all job descriptors and tries to reclaim each that has no
// terminal done record and no live hold. SetNX on the hold guarantees a single
// owner across nodes.
func (p *Parser) recoverJobs(ctx context.Context) {
	jobs, err := p.progMgr.ListJobs(ctx)
	if err != nil {
		p.logger.WithContext(ctx).Errorf("recover: list jobs: %s", err)
		return
	}
	for _, job := range jobs {
		p.tryReclaim(ctx, job)
	}
}

// tryReclaim takes ownership of one abandoned job: skip if finished or still
// held, else steal the hold, re-derive the rule, and re-run through the same
// lifecycle as a fresh dispatch. The rpa extract stage resumes from its
// persisted checkpoint — it recalls already-extracted envelopes (fileId →
// content) and skips those files — so a half-done rpa run continues without
// re-calling the LLM for finished files or losing their results. Escrow re-runs
// outright (idempotent per-file upsert, typically a single file). Write-back is
// idempotent, so either path converges to the same row state.
func (p *Parser) tryReclaim(ctx context.Context, job *progress.Job) {
	if done, _ := p.progMgr.GetDone(ctx, job.RequestId); done != nil {
		return
	}
	ok, err := p.progMgr.AcquireHold(ctx, job.RequestId, p.nodeId)
	if err != nil {
		p.logger.WithContext(ctx).Errorf("recover: acquire hold %s: %s", job.RequestId, err)
		return
	}
	if !ok {
		return // another node owns it
	}
	doc, err := p.repos.DocumentRepo.FindOneWithFiles(ctx, job.DocumentId)
	if err != nil || doc == nil {
		_ = p.progMgr.ReleaseHold(ctx, job.RequestId)
		return
	}
	rule, ok := resolveRule(docFileParseConfig(doc), job.Parser)
	if !ok {
		p.finalize(ctx, job.TransactionId, job.RequestId, job.RuleSig)
		_ = p.progMgr.WriteDone(ctx, job.RequestId, progress.StatusFailed, "rule no longer resolvable")
		return
	}
	h := p.handlers[job.Parser]
	if h == nil {
		p.finalize(ctx, job.TransactionId, job.RequestId, job.RuleSig)
		return
	}
	// Hand the document's current files to the run: a passed-file scope
	// (SingleFile/Files) re-parses them; a collected scope (Directory/TagFiltered)
	// ignores them and re-gathers from the DB.
	files := make([]model.File, 0, len(doc.Files))
	for _, f := range doc.Files {
		files = append(files, *f)
	}
	threading.GoSafe(func() {
		p.run(job.TransactionId, job.DocumentId, job.RequestId, job.RuleSig, h, rule, files)
	})
}
