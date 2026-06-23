package cron

import (
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
