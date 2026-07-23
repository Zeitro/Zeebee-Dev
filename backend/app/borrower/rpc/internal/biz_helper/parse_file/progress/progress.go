package progress

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
)

// TTLs (seconds). lock/progress slide per file; done/reqmap/job outlive the
// run so late progress queries and restart recovery still resolve; hold is
// short so a crashed node's task is reclaimable quickly.
const (
	ttlLockSec     = 600
	ttlProgressSec = 600
	ttlDoneSec     = 3600
	ttlHoldSec     = 30
	ttlReqMapSec   = 3600
	ttlJobSec      = 3600
	HoldHeartbeat  = 10 * time.Second
)

// Job-level status.
const (
	StatusRunning  = "Running"
	StatusFinished = "Finished"
	StatusFailed   = "Failed"
)

// Stage names.
const (
	StageCollect   = "Collect"
	StageExtract   = "Extract"
	StageWriteback = "Writeback"
)

// Per-stage status.
const (
	StageStatusPending = "Pending"
	StageStatusRunning = "Running"
	StageStatusDone    = "Done"
)

// StageState is one stage of the staged pipeline: name, status, and processed/
// total counters the frontend composes into progress (no overall percent).
type StageState struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Processed int    `json:"processed"`
	Total     int    `json:"total"`
}

// Progress is the single Redis record holding all stages. No overall percent —
// the frontend composes it from stage + processed/total.
type Progress struct {
	Status string       `json:"status"` // Running | Finished | Failed
	Stage  string       `json:"stage"`  // current stage name
	Stages []StageState `json:"stages"`
	Error  string       `json:"error,omitempty"`
}

// DoneState is the long-TTL terminal record queried after the run key expires.
type DoneState struct {
	Status string `json:"status"` // Finished | Failed
	Error  string `json:"error,omitempty"`
}

// Job is the recovery descriptor scanned on startup/sweep. It carries identity
// + collect input params only; all stage I/O lives in Progress.
type Job struct {
	RequestId     string `json:"request_id"`
	TransactionId int64  `json:"transaction_id"`
	RuleSig       string `json:"rule_sig"`
	Parser        string `json:"parser"`
	Mode          string `json:"mode"`
	DocumentId    int64  `json:"document_id"`
}

// ProgressManager owns all bpt:parse:* Redis keys.
type ProgressManager struct {
	rdb    *redisx.Client
	logger logx.Logger
}

func NewProgressManager(rdb *redisx.Client, logger logx.Logger) *ProgressManager {
	return &ProgressManager{rdb: rdb, logger: logger}
}

// ---- dedup lock (per transaction+ruleSig) ----

func (m *ProgressManager) AcquireDedupLock(ctx context.Context, txId int64, ruleSig, requestId string) (bool, error) {
	return m.rdb.SetNX(ctx, lockKey(txId, ruleSig), requestId, time.Duration(ttlLockSec)*time.Second)
}

func (m *ProgressManager) RenewDedupLock(ctx context.Context, txId int64, ruleSig string) error {
	return m.rdb.Expire(ctx, lockKey(txId, ruleSig), ttlLockSec)
}

func (m *ProgressManager) ReleaseDedupLock(ctx context.Context, txId int64, ruleSig string) error {
	return m.rdb.Del(ctx, lockKey(txId, ruleSig))
}

// ---- hold lock (per request, heartbeat-driven, for crash detection/reclaim) ----

func (m *ProgressManager) AcquireHold(ctx context.Context, requestId, nodeId string) (bool, error) {
	return m.rdb.SetNX(ctx, holdKey(requestId), nodeId, time.Duration(ttlHoldSec)*time.Second)
}

func (m *ProgressManager) RenewHold(ctx context.Context, requestId string) error {
	return m.rdb.Expire(ctx, holdKey(requestId), ttlHoldSec)
}

func (m *ProgressManager) ReleaseHold(ctx context.Context, requestId string) error {
	return m.rdb.Del(ctx, holdKey(requestId))
}

// ---- progress / done ----

func (m *ProgressManager) InitProgress(ctx context.Context, requestId string, stages []StageState) error {
	p := &Progress{Status: StatusRunning, Stage: FirstStage(stages), Stages: stages}
	return m.saveProgress(ctx, requestId, p)
}

func (m *ProgressManager) SaveProgress(ctx context.Context, requestId string, p *Progress) error {
	return m.saveProgress(ctx, requestId, p)
}

func (m *ProgressManager) saveProgress(ctx context.Context, requestId string, p *Progress) error {
	b, err := json.Marshal(p)
	if err != nil {
		return errorx.Wrap(err, "marshal progress")
	}
	if err := m.rdb.SetExpire(ctx, progressKey(requestId), string(b), ttlProgressSec); err != nil {
		return err
	}
	return nil
}

func (m *ProgressManager) GetProgress(ctx context.Context, requestId string) (*Progress, error) {
	raw, err := m.rdb.Get(ctx, progressKey(requestId))
	if err == nil && raw != "" {
		var p Progress
		if json.Unmarshal([]byte(raw), &p) == nil {
			return &p, nil
		}
	}
	return nil, nil
}

func (m *ProgressManager) WriteDone(ctx context.Context, requestId, status, errMsg string) error {
	d := DoneState{Status: status, Error: errMsg}
	b, err := json.Marshal(d)
	if err != nil {
		return errorx.Wrap(err, "marshal done")
	}
	if err := m.rdb.SetExpire(ctx, doneKey(requestId), string(b), ttlDoneSec); err != nil {
		return err
	}
	return m.rdb.Del(ctx, progressKey(requestId))
}

func (m *ProgressManager) GetDone(ctx context.Context, requestId string) (*DoneState, error) {
	raw, err := m.rdb.Get(ctx, doneKey(requestId))
	if err != nil || raw == "" {
		return nil, nil
	}
	var d DoneState
	if json.Unmarshal([]byte(raw), &d) != nil {
		return nil, nil
	}
	return &d, nil
}

// ---- ruleSig -> requestId map (progress-by-params mode b) ----

func (m *ProgressManager) SetReqMap(ctx context.Context, txId int64, ruleSig, requestId string) error {
	return m.rdb.SetExpire(ctx, reqMapKey(txId, ruleSig), requestId, ttlReqMapSec)
}

func (m *ProgressManager) GetReqMap(ctx context.Context, txId int64, ruleSig string) (string, error) {
	return m.rdb.Get(ctx, reqMapKey(txId, ruleSig))
}

// ---- stage checkpoint (extract resume) ----
// The checkpoint (marshaled map[int64] of already-extracted per-file results)
// lives in its OWN key at job TTL (ttlJobSec), not in the sliding progress
// record (ttlProgressSec). A crashed node's progress record can expire before
// the recovery sweeper reclaims the job; a job-TTL checkpoint survives that
// gap so a resumed rpa extract still skips finished files instead of re-calling
// the LLM.

func (m *ProgressManager) SetStageCheckpoint(ctx context.Context, requestId, stage string, output json.RawMessage) error {
	return m.rdb.SetExpire(ctx, ckptKey(requestId, stage), string(output), ttlJobSec)
}

func (m *ProgressManager) GetStageCheckpoint(ctx context.Context, requestId, stage string) (json.RawMessage, error) {
	raw, err := m.rdb.Get(ctx, ckptKey(requestId, stage))
	if err != nil || raw == "" {
		return nil, nil
	}
	return json.RawMessage(raw), nil
}

// ---- job descriptor (recovery scan) ----

func (m *ProgressManager) WriteJob(ctx context.Context, job *Job) error {
	b, err := json.Marshal(job)
	if err != nil {
		return errorx.Wrap(err, "marshal job")
	}
	return m.rdb.SetExpire(ctx, jobKey(job.RequestId), string(b), ttlJobSec)
}

func (m *ProgressManager) DeleteJob(ctx context.Context, requestId string) error {
	return m.rdb.Del(ctx, jobKey(requestId))
}

// ListJobs scans all bpt:parse:job:* descriptors and decodes them. Used by the
// recovery coordinator; not a hot-path call.
func (m *ProgressManager) ListJobs(ctx context.Context) ([]*Job, error) {
	keys, err := m.rdb.Scan(ctx, "bpt:parse:job:*")
	if err != nil {
		return nil, errorx.Wrap(err, "scan jobs")
	}
	jobs := make([]*Job, 0, len(keys))
	for _, key := range keys {
		raw, err := m.rdb.Get(ctx, key)
		if err != nil || raw == "" {
			continue
		}
		var job Job
		if json.Unmarshal([]byte(raw), &job) != nil {
			continue
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

func FirstStage(stages []StageState) string {
	for _, s := range stages {
		if s.Status != StageStatusDone {
			return s.Name
		}
	}
	if len(stages) > 0 {
		return stages[0].Name
	}
	return ""
}

// ---- key builders ----
func lockKey(txId int64, ruleSig string) string {
	return fmt.Sprintf("bpt:parse:lock:%d:%s", txId, ruleSig)
}
func progressKey(requestId string) string { return "bpt:parse:progress:" + requestId }
func doneKey(requestId string) string     { return "bpt:parse:done:" + requestId }
func holdKey(requestId string) string     { return "bpt:parse:hold:" + requestId }
func reqMapKey(txId int64, ruleSig string) string {
	return fmt.Sprintf("bpt:parse:reqmap:%d:%s", txId, ruleSig)
}
func jobKey(requestId string) string { return "bpt:parse:job:" + requestId }
func ckptKey(requestId, stage string) string {
	return fmt.Sprintf("bpt:parse:ckpt:%s:%s", requestId, stage)
}
