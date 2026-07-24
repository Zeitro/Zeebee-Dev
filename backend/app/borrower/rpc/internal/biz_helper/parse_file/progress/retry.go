package progress

import (
	"context"
	"time"
)

const (
	retryMax    = 3
	retryBaseMs = 2000 // 2s, 4s exponential backoff
)

// RetryStore reads/persists a stage's retry counter (job TTL). Recovered runs continue counting.
type RetryStore interface {
	StageRetry(ctx context.Context, stage string) RetryState
	SaveStageRetry(ctx context.Context, stage string, retry RetryState)
}

// RetryDo retries fn for one stage, continuing the persisted count across recovery.
// Extract wraps the single consolidated multimodal call; Collect/Writeback wrap whole-stage redo.
func RetryDo(ctx context.Context, store RetryStore, stage string, fn func() error) error {
	rs := store.StageRetry(ctx, stage)
	if rs.Max == 0 {
		rs.Max = retryMax
	}
	var err error
	for attempt := rs.Count; attempt < rs.Max; attempt++ {
		err = fn()
		if err == nil {
			store.SaveStageRetry(ctx, stage, RetryState{Max: rs.Max}) // reset on success
			return nil
		}
		rs.Count = attempt + 1
		rs.BackoffMs = retryBaseMs << attempt
		store.SaveStageRetry(ctx, stage, rs)
		if attempt == rs.Max-1 {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(retryBaseMs<<attempt) * time.Millisecond):
		}
	}
	return err
}
