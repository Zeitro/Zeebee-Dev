package cron

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewEngine(t *testing.T) {
	e := NewEngine()
	if e == nil {
		t.Fatal("NewEngine returned nil")
	}
}

func TestEngine_Register_DuplicateName(t *testing.T) {
	e := NewEngine()
	if err := e.Register("task1", "* * * * * *", func() {}); err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	err := e.Register("task1", "* * * * * *", func() {})
	if err == nil {
		t.Fatal("expected error on duplicate name, got nil")
	}
	if err.Error() != "task task1 already exists" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEngine_Register_InvalidSpec(t *testing.T) {
	e := NewEngine()

	err := e.Register("task1", "invalid-spec", func() {})
	if err == nil {
		t.Fatal("expected error on invalid spec, got nil")
	}
}

func TestEngine_Register_Validation(t *testing.T) {
	cases := []struct {
		name string
		task string
		spec string
		cmd  func()
	}{
		{"empty name", "", "* * * * * *", func() {}},
		{"empty spec", "task1", "", func() {}},
		{"nil cmd", "task1", "* * * * * *", nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := NewEngine()
			err := e.Register(c.task, c.spec, c.cmd)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

func TestEngine_Unregister(t *testing.T) {
	e := NewEngine()
	if err := e.Register("task1", "* * * * * *", func() {}); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if err := e.Unregister("task1"); err != nil {
		t.Fatalf("unregister failed: %v", err)
	}
	if e.IsExisted("task1") {
		t.Error("task1 should be removed after Unregister")
	}
}

func TestEngine_Unregister_NotExist(t *testing.T) {
	e := NewEngine()

	err := e.Unregister("nonexistent")
	if err == nil {
		t.Fatal("expected error on unregister non-existent task, got nil")
	}
	if err.Error() != "task nonexistent does not exist" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEngine_IsExisted(t *testing.T) {
	e := NewEngine()

	if e.IsExisted("task1") {
		t.Error("expected IsExisted to return false before register")
	}

	if err := e.Register("task1", "* * * * * *", func() {}); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if !e.IsExisted("task1") {
		t.Error("expected IsExisted to return true after register")
	}
}

func TestEngine_StartStop(t *testing.T) {
	e := NewEngine()
	e.Start()
	e.Stop()
}

// TestEngine_TaskExecutes verifies a registered task actually fires.
func TestEngine_TaskExecutes(t *testing.T) {
	e := NewEngine()
	e.Start()
	defer e.Stop()

	var fired atomic.Int32
	done := make(chan struct{})

	// Spec: run every second, with seconds field enabled.
	if err := e.Register("every-second", "* * * * * *", func() {
		if fired.Add(1) >= 1 {
			select {
			case done <- struct{}{}:
			default:
			}
		}
	}); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("task did not fire within 5 seconds")
	}
}

// TestEngine_UnregisterStopsExecution verifies removing a task halts it.
func TestEngine_UnregisterStopsExecution(t *testing.T) {
	e := NewEngine()
	e.Start()
	defer e.Stop()

	var count atomic.Int32
	if err := e.Register("counting", "* * * * * *", func() {
		count.Add(1)
	}); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Let it run once or twice.
	time.Sleep(2500 * time.Millisecond)
	if err := e.Unregister("counting"); err != nil {
		t.Fatalf("unregister failed: %v", err)
	}

	before := count.Load()
	// Wait and ensure no further firings.
	time.Sleep(2500 * time.Millisecond)
	after := count.Load()

	if after > before {
		t.Errorf("task fired after unregister: before=%d after=%d", before, after)
	}
}

// TestEngine_ConcurrentRegister confirms concurrent registration does not
// panic and only one task ends up registered (Register is atomic under e.mu,
// so the first writer wins; the rest see "already exists").
func TestEngine_ConcurrentRegister(t *testing.T) {
	e := NewEngine()
	e.Start()
	defer e.Stop()

	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_ = e.Register("shared", "* * * * * *", func() {})
		}()
	}
	wg.Wait()

	if !e.IsExisted("shared") {
		t.Error("expected task 'shared' to be registered")
	}
}

// --- Replace tests ---

func TestEngine_Replace_NewTask(t *testing.T) {
	e := NewEngine()
	if err := e.Replace("task1", "* * * * * *", func() {}); err != nil {
		t.Fatalf("replace on non-existent failed: %v", err)
	}
	if !e.IsExisted("task1") {
		t.Error("expected task1 registered after replace on non-existent")
	}
}

func TestEngine_Replace_Existing(t *testing.T) {
	e := NewEngine()
	e.Start()
	defer e.Stop()

	var count atomic.Int32
	if err := e.Register("task1", "* * * * * *", func() { count.Add(1) }); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	// Baseline: every-second spec fires ~2 ticks in 2.2s.
	time.Sleep(2200 * time.Millisecond)
	oldRate := count.Load()

	// Switch to every-5s; in <5s it fires at most once.
	if err := e.Replace("task1", "*/5 * * * * *", func() { count.Add(1) }); err != nil {
		t.Fatalf("replace existing failed: %v", err)
	}
	if !e.IsExisted("task1") {
		t.Fatal("expected task1 still registered after replace")
	}
	time.Sleep(2200 * time.Millisecond)
	newRate := count.Load() - oldRate

	if newRate > 1 {
		t.Errorf("new schedule should be slower than old: oldRate=%d newRate=%d", oldRate, newRate)
	}
}

func TestEngine_Replace_Validation(t *testing.T) {
	cases := []struct {
		name string
		task string
		spec string
		cmd  func()
	}{
		{"empty name", "", "* * * * * *", func() {}},
		{"empty spec", "task1", "", func() {}},
		{"nil cmd", "task1", "* * * * * *", nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := NewEngine()
			if err := e.Replace(c.task, c.spec, c.cmd); err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

// TestEngine_Replace_InvalidSpec_NewTask: invalid spec on non-existent task
// returns error and registers nothing.
func TestEngine_Replace_InvalidSpec_NewTask(t *testing.T) {
	e := NewEngine()
	if err := e.Replace("task1", "invalid-spec", func() {}); err == nil {
		t.Fatal("expected error on invalid spec, got nil")
	}
	if e.IsExisted("task1") {
		t.Error("task1 must not be registered when spec invalid")
	}
}

// TestEngine_Replace_InvalidSpec_PreservesExisting verifies a bad new spec
// does NOT destroy the currently-registered task: Replace adds the new entry
// before removing the old one, so AddFunc failure leaves the old task intact.
func TestEngine_Replace_InvalidSpec_PreservesExisting(t *testing.T) {
	e := NewEngine()
	if err := e.Register("task1", "* * * * * *", func() {}); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if err := e.Replace("task1", "invalid-spec", func() {}); err == nil {
		t.Fatal("expected error on invalid spec, got nil")
	}
	if !e.IsExisted("task1") {
		t.Error("task1 should survive a failed replace")
	}
}

// TestEngine_Replace_ExecutesUpdatedCmd verifies the old cmd stops firing and
// the new cmd starts firing after Replace.
func TestEngine_Replace_ExecutesUpdatedCmd(t *testing.T) {
	e := NewEngine()
	e.Start()
	defer e.Stop()

	var old, neu atomic.Int32
	if err := e.Register("task1", "* * * * * *", func() { old.Add(1) }); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	time.Sleep(1500 * time.Millisecond)
	if old.Load() == 0 {
		t.Fatal("old cmd never fired before replace")
	}

	if err := e.Replace("task1", "* * * * * *", func() { neu.Add(1) }); err != nil {
		t.Fatalf("replace failed: %v", err)
	}
	oldAtReplace := old.Load()
	time.Sleep(2500 * time.Millisecond)

	if neu.Load() == 0 {
		t.Error("new cmd never fired after replace")
	}
	// Tolerate one in-flight tick already dispatched before Remove took effect.
	if old.Load() > oldAtReplace+1 {
		t.Errorf("old cmd kept firing after replace: before=%d after=%d", oldAtReplace, old.Load())
	}
}

// TestEngine_ConcurrentReplace_SameName confirms concurrent Replace of the same
// name does not panic under -race and leaves the map consistent (the name ends
// registered exactly once). Replace is atomic under e.mu, so every goroutine
// either fully replaces or observes the final replaced state — no errors.
func TestEngine_ConcurrentReplace_SameName(t *testing.T) {
	e := NewEngine()
	e.Start()
	defer e.Stop()

	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	var errCount atomic.Int32
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if err := e.Replace("shared", "* * * * * *", func() {}); err != nil {
				errCount.Add(1)
			}
		}()
	}
	wg.Wait()

	if !e.IsExisted("shared") {
		t.Fatal("expected 'shared' registered after concurrent replace")
	}
	// map consistent: re-registering the same name must report duplicate.
	if err := e.Register("shared", "* * * * * *", func() {}); err == nil {
		t.Error("re-register of 'shared' should report duplicate after concurrent replace")
	}
	if errCount.Load() != 0 {
		t.Errorf("expected no errors under atomic replace, got %d/%d", errCount.Load(), n)
	}
}

// TestEngine_ConcurrentReplace_DistinctNames confirms concurrent Replace of
// distinct names all succeed (no cross-name interference).
func TestEngine_ConcurrentReplace_DistinctNames(t *testing.T) {
	e := NewEngine()
	e.Start()
	defer e.Stop()

	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	var errCount atomic.Int32
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			if err := e.Replace(fmt.Sprintf("task-%d", i), "* * * * * *", func() {}); err != nil {
				errCount.Add(1)
			}
		}(i)
	}
	wg.Wait()

	if errCount.Load() > 0 {
		t.Errorf("expected all distinct-name replaces to succeed, got %d errors", errCount.Load())
	}
	for i := 0; i < n; i++ {
		if !e.IsExisted(fmt.Sprintf("task-%d", i)) {
			t.Errorf("task-%d missing after concurrent replace", i)
		}
	}
}

// TestEngine_ConcurrentReplace_MixedOps races Replace against Unregister and
// Register on the same name. No panic, final state usable.
func TestEngine_ConcurrentReplace_MixedOps(t *testing.T) {
	e := NewEngine()
	e.Start()
	defer e.Stop()

	if err := e.Register("shared", "* * * * * *", func() {}); err != nil {
		t.Fatalf("seed register failed: %v", err)
	}

	const n = 60
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			switch i % 3 {
			case 0:
				_ = e.Replace("shared", "* * * * * *", func() {})
			case 1:
				_ = e.Unregister("shared")
			case 2:
				_ = e.Register("shared", "* * * * * *", func() {})
			}
		}(i)
	}
	wg.Wait()

	// engine still functional regardless of final state.
	_ = e.Replace("shared", "* * * * * *", func() {})
	if !e.IsExisted("shared") {
		t.Fatal("engine unusable after mixed concurrent ops")
	}
}
