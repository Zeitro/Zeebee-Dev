package parse_file

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
	ttlQueueSec    = 3600
	ttlJobSec      = 3600
	holdHeartbeat  = 10 * time.Second
)

// Job-level status.
const (
	StatusRunning   = "Running"
	StatusFinished  = "Finished"
	StatusFailed    = "Failed"
)

// Stage names.
const (
	StageCollect    = "Collect"
	StageExtract    = "Extract"
	StageArbitrate  = "Arbitrate"
	StageWriteback  = "Writeback"
)

// Per-stage status.
const (
	StageStatusPending = "Pending"
	StageStatusRunning = "Running"
	StageStatusDone    = "Done"
	StageStatusSkipped = "Skipped"
)

// RetryState tracks per-stage retry. Extract counts per file (transient, not
// persisted across restart); other stages count per stage attempt.
type RetryState struct {
	Count    int   `json:"count"`
	Max      int   `json:"max"`
	BackoffMs int64 `json:"backoff_ms"`
}

// StageState is one stage of the staged pipeline. State holds intra-stage
// resume data (e.g. Extract's processedFileIds); Output holds the handoff
// payload consumed by the next stage (FR-028 pipeline handoff).
type StageState struct {
	Name      string          `json:"name"`
	Status    string          `json:"status"`
	Processed int             `json:"processed"`
	Total     int             `json:"total"`
	State     json.RawMessage `json:"state,omitempty"`
	Output    json.RawMessage `json:"output,omitempty"`
	Retry     RetryState      `json:"retry"`
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
	RequestId    string `json:"request_id"`
	TransactionId int64 `json:"transaction_id"`
	RuleSig      string `json:"rule_sig"`
	Parser       string `json:"parser"`
	Mode         string `json:"mode"`
	DocumentId   int64  `json:"document_id"`
}

// ProgressManager owns all bpt:parse:* Redis keys.
type ProgressManager struct {
	rdb    *redisx.Client
	logger logx.Logger
}

func NewProgressManager(rdb *redisx.Client, logger logx.Logger) *ProgressManager {
	return &ProgressManager{rdb: rdb, logger: logger}
}

// ---- key builders ----

func lockKey(txId int64, ruleSig string) string {
	return fmt.Sprintf("bpt:parse:lock:%d:%s", txId, ruleSig)
}
func progressKey(requestId string) string  { return "bpt:parse:progress:" + requestId }
func doneKey(requestId string) string      { return "bpt:parse:done:" + requestId }
func holdKey(requestId string) string      { return "bpt:parse:hold:" + requestId }
func reqMapKey(txId int64, ruleSig string) string {
	return fmt.Sprintf("bpt:parse:reqmap:%d:%s", txId, ruleSig)
}
func queueKey(txId, docId int64) string {
	return fmt.Sprintf("bpt:parse:queue:%d:%d", txId, docId)
}
func jobKey(requestId string) string { return "bpt:parse:job:" + requestId }

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
	p := &Progress{Status: StatusRunning, Stage: firstStage(stages), Stages: stages}
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

// ---- on_upload per-document serial queue ----

func (m *ProgressManager) EnqueueFile(ctx context.Context, txId, docId int64, fileId string) error {
	if err := m.rdb.PushBack(ctx, queueKey(txId, docId), fileId); err != nil {
		return err
	}
	return m.rdb.Expire(ctx, queueKey(txId, docId), ttlQueueSec)
}

func (m *ProgressManager) DequeueFile(ctx context.Context, txId, docId int64) (string, bool, error) {
	return m.rdb.PopFront(ctx, queueKey(txId, docId))
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

func firstStage(stages []StageState) string {
	for _, s := range stages {
		if s.Status != StageStatusDone && s.Status != StageStatusSkipped {
			return s.Name
		}
	}
	if len(stages) > 0 {
		return stages[0].Name
	}
	return ""
}
