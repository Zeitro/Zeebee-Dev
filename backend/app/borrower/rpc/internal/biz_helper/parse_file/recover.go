package parse_file

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/domains/zeebee/model"
)

// recoverSweepInterval: how often the sweeper re-scans for abandoned jobs —
// short enough to reclaim a crashed node fast, long enough to idle quietly.
const recoverSweepInterval = 30 * time.Second

// StartRecovery runs the background coordinator: an immediate scan on startup
// (resume in-flight parses, incl. a half-drained escrow queue) + periodic sweep
// reclaiming crashed-node work. Runs for the life of the process.
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

// recoverJobs scans all jobs and reclaims each with no terminal Done and no live
// hold. SetNX on the hold guarantees one owner across nodes.
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

// tryReclaim steals one abandoned job: skip if finished or still held, else take
// the hold, re-derive the rule, re-run via the same lifecycle. RPA extract
// resumes from its checkpoint (recalls extracted envelopes, skips those files);
// escrow re-runs outright (idempotent per-file upsert). Write-back idempotent →
// both converge.
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
	// Pass the doc's current files: passed-file scope re-parses them; collected
	// scope ignores them and re-gathers from the DB.
	files := make([]model.File, 0, len(doc.Files))
	for _, f := range doc.Files {
		files = append(files, *f)
	}
	threading.GoSafe(func() {
		p.run(ctx, job.TransactionId, job.DocumentId, job.RequestId, job.RuleSig, h, rule, files)
	})
}
