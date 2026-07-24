package progress

import (
	"context"
	"errors"
	"testing"
	"time"
)

type memRetryStore struct{ retries map[string]RetryState }

func newRetryStore() *memRetryStore { return &memRetryStore{retries: map[string]RetryState{}} }

func (m *memRetryStore) StageRetry(_ context.Context, stage string) RetryState  { return m.retries[stage] }
func (m *memRetryStore) SaveStageRetry(_ context.Context, stage string, r RetryState) {
	m.retries[stage] = r
}

func TestRetryDoSucceedsFirstAttempt(t *testing.T) {
	calls := 0
	err := RetryDo(context.Background(), newRetryStore(), StageExtract, func() error {
		calls++
		return nil
	})
	if err != nil || calls != 1 {
		t.Fatalf("err=%v calls=%d, want nil/1", err, calls)
	}
}

func TestRetryDoSucceedsAfterRetry(t *testing.T) {
	calls := 0
	start := time.Now()
	err := RetryDo(context.Background(), newRetryStore(), StageExtract, func() error {
		calls++
		if calls < 2 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("err=%v, want nil", err)
	}
	if calls != 2 {
		t.Fatalf("calls=%d, want 2", calls)
	}
	if elapsed := time.Since(start); elapsed < time.Second {
		t.Fatalf("no backoff applied: elapsed=%v", elapsed)
	}
}

func TestRetryDoContextCancelAbortsBackoff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	calls := 0
	start := time.Now()
	err := RetryDo(ctx, newRetryStore(), StageExtract, func() error {
		calls++
		return errors.New("always")
	})
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if elapsed > 2*time.Second {
		t.Fatalf("did not abort on ctx cancel: elapsed=%v", elapsed)
	}
}

// TestRetryDoContinuesPersistedCount: a pre-seeded count=2 means only one more
// attempt runs (T009b) — recovery continues counting instead of restarting.
func TestRetryDoContinuesPersistedCount(t *testing.T) {
	store := newRetryStore()
	store.retries[StageExtract] = RetryState{Count: 2, Max: retryMax}
	calls := 0
	_ = RetryDo(context.Background(), store, StageExtract, func() error {
		calls++
		return errors.New("still failing")
	})
	if calls != 1 {
		t.Fatalf("calls=%d, want 1 (only the remaining attempt from count=2)", calls)
	}
}
