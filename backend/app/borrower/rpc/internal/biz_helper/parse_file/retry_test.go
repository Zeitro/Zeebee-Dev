package parse_file

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryDoSucceedsFirstAttempt(t *testing.T) {
	calls := 0
	err := retryDo(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil || calls != 1 {
		t.Fatalf("err=%v calls=%d, want nil/1", err, calls)
	}
}

func TestRetryDoSucceedsAfterRetry(t *testing.T) {
	// Second attempt succeeds — exercises one exponential backoff (~2s).
	calls := 0
	start := time.Now()
	err := retryDo(context.Background(), func() error {
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
	// Always-fail fn + a ctx cancelled mid-first-backoff → retryDo returns
	// ctx.Err() quickly instead of waiting out all retries (~6s).
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	calls := 0
	start := time.Now()
	err := retryDo(ctx, func() error {
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
