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

// TTLs (seconds). lock/progress slide per file; done/reqmap/job outlive the run.
const (
	ttlLockSec     = 600
	ttlProgressSec = 600
	ttlDoneSec     = 3600
	ttlHoldSec     = 30
	ttlReqMapSec   = 3600
	ttlJobSec      = 3600
	HoldHeartbeat  = 10 * time.Second
)

const (
	StatusRunning  = "Running"
	StatusFinished = "Finished"
	StatusFailed   = "Failed"
)

const (
	StageCollect   = "Collect"
	StageExtract   = "Extract"
	StageWriteback = "Writeback"
)

const (
	StageStatusPending = "Pending"
	StageStatusRunning = "Running"
	StageStatusDone    = "Done"
)

// RetryState is the persisted per-stage retry counter (job TTL); recovery continues counting (T009b).
type RetryState struct {
	Count     int `json:"count"`
	Max       int `json:"max"`
	BackoffMs int `json:"backoff_ms"`
}

// StageState is one pipeline stage; State holds a stage-specific blob (Extract {uploaded,total}).
type StageState struct {
	Name      string          `json:"name"`
	Status    string          `json:"status"`
	Processed int             `json:"processed"`
	Total     int             `json:"total"`
	State     json.RawMessage `json:"state,omitempty"`
	Retry     RetryState      `json:"retry"`
}

// Progress holds all stages; no overall percent — frontend composes from processed/total.
type Progress struct {
	Status string       `json:"status"`
	Stage  string       `json:"stage"`
	Stages []StageState `json:"stages"`
	Error  string       `json:"error,omitempty"`
}

type DoneState struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Job is the recovery descriptor: identity + collect params only; stage I/O lives in Progress + the files checkpoint.
type Job struct {
	RequestId     string `json:"request_id"`
	TransactionId int64  `json:"transaction_id"`
	RuleSig       string `json:"rule_sig"`
	Parser        string `json:"parser"`
	Mode          string `json:"mode"`
	DocumentId    int64  `json:"document_id"`
}

// FileCheckpoint is one uploaded file's provider reference; the id is provider-scoped, not portable (R11).
type FileCheckpoint struct {
	Provider       string `json:"provider"`
	ProviderFileID string `json:"provider_file_id"`
	// MaxPages: slice cap (0 = whole file); a change invalidates the cached id.
	MaxPages int `json:"max_pages"`
}

type extractState struct {
	Uploaded int `json:"uploaded"`
	Total    int `json:"total"`
}

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

// ---- hold lock (per request, heartbeat-driven, crash reclaim) ----

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
	return m.rdb.SetExpire(ctx, progressKey(requestId), string(b), ttlProgressSec)
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

// ---- Extract files checkpoint (resume: fileId → provider file id) ----
// Own key at job TTL (survives expired progress). One run = one hold → single writer, RMW safe.

func (m *ProgressManager) SaveFileCheckpoint(ctx context.Context, requestId string, fileId int64, fc FileCheckpoint) error {
	files, _ := m.GetFilesCheckpoint(ctx, requestId)
	if files == nil {
		files = map[int64]FileCheckpoint{}
	}
	files[fileId] = fc
	b, err := json.Marshal(files)
	if err != nil {
		return errorx.Wrap(err, "marshal files checkpoint")
	}
	return m.rdb.SetExpire(ctx, filesKey(requestId), string(b), ttlJobSec)
}

func (m *ProgressManager) GetFilesCheckpoint(ctx context.Context, requestId string) (map[int64]FileCheckpoint, error) {
	raw, err := m.rdb.Get(ctx, filesKey(requestId))
	if err != nil || raw == "" {
		return nil, nil
	}
	var files map[int64]FileCheckpoint
	if json.Unmarshal([]byte(raw), &files) != nil {
		return nil, nil
	}
	return files, nil
}

func (m *ProgressManager) DeleteFilesCheckpoint(ctx context.Context, requestId string) error {
	return m.rdb.Del(ctx, filesKey(requestId))
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

// ListJobs scans all bpt:parse:job:* descriptors. Recovery-only, not hot path.
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

// ExtractStateJSON marshals the Extract upload-progress State blob.
func ExtractStateJSON(uploaded, total int) json.RawMessage {
	b, _ := json.Marshal(extractState{Uploaded: uploaded, Total: total})
	return b
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
func jobKey(requestId string) string   { return "bpt:parse:job:" + requestId }
func filesKey(requestId string) string { return "bpt:parse:files:" + requestId }
