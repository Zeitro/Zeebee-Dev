package progress

import (
	"context"
	"time"
)

// Per-stage retry policy. Extract retries per file (a bad file must not abort
// the whole parse); Collect/Arbitrate/Writeback retry the whole stage. Each
// stage counts its own attempts independently (FR-030).
const (
	retryMax    = 3
	retryBaseMs = 2000 // 2s, 4s, 8s exponential backoff
)

// retryDo calls fn up to retryMax times with exponential backoff (2s, 4s, 8s).
// It returns the last error once attempts are exhausted. The backoff sleeps
// are ctx-aware so a cancelled parse does not block.
func RetryDo(ctx context.Context, fn func() error) error {
	var err error
	for attempt := 0; attempt < retryMax; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		if attempt == retryMax-1 {
			break
		}
		backoff := time.Duration(retryBaseMs<<(attempt)) * time.Millisecond
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}
	return err
}

// retryStateFor returns the RetryState snapshot for persistence in progress.
func RetryStateFor(attempt int) RetryState {
	return RetryState{Count: attempt, Max: retryMax, BackoffMs: int64(retryBaseMs)}
}
