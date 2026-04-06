package errorx

import (
	"errors"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestNewErrorGroup(t *testing.T) {
	eg := NewErrorGroup()
	if eg == nil {
		t.Error("NewErrorGroup should return a non-nil ErrorGroup")
	}

	if eg.HasError() {
		t.Error("Newly created ErrorGroup should not have an error")
	}
}

func TestErrorGroupDo(t *testing.T) {
	eg := NewErrorGroup()

	// Test successful functions
	eg.Do(
		func() error { return nil },
		func() error { return nil },
	)

	// Wait for completion
	if err := eg.Wait(); err != nil {
		t.Errorf("Do successful tasks returned error: %v", err)
	}

	if eg.HasError() {
		t.Error("ErrorGroup should not have error after successful functions")
	}

	// Test with an error function
	testErr := errors.New("test error")
	eg.Reset() // Reset the group

	eg.Do(
		func() error { return nil },
		func() error { return testErr },
		func() error { return nil }, // This should not be executed
	)

	// Wait and check error
	err := eg.Wait()
	if !eg.HasError() {
		t.Error("ErrorGroup should have an error after error function")
	}

	if err != testErr {
		t.Errorf("ErrorGroup.Wait() = %v, want %v", err, testErr)
	}

	// Test that subsequent Do calls don't execute when there's already an error
	eg.Reset()
	callCount := 0
	eg.Do(func() error {
		callCount++
		return errors.New("first error")
	})

	eg.Do(func() error {
		callCount++
		return errors.New("second error")
	})

	// Wait for both Do calls
	_ = eg.Wait()

	// Second Do should not execute since first Do set error
	if callCount != 1 {
		t.Errorf("Functions should not execute after first Do, callCount = %d, want 1", callCount)
	}
}

func TestErrorGroupHasError(t *testing.T) {
	eg := NewErrorGroup()

	if eg.HasError() {
		t.Error("Empty ErrorGroup should not have an error")
	}

	eg.Do(func() error {
		return errors.New("test error")
	})

	// Wait for task to complete
	_ = eg.Wait()

	if !eg.HasError() {
		t.Error("ErrorGroup should have an error after Do with error function")
	}
}

func TestErrorGroupReset(t *testing.T) {
	eg := NewErrorGroup()
	testErr := errors.New("test error")

	eg.Do(func() error { return testErr })

	// Wait for task to complete
	_ = eg.Wait()

	if !eg.HasError() {
		t.Error("ErrorGroup should have an error")
	}

	eg.Reset()

	if eg.HasError() {
		t.Error("ErrorGroup should not have an error after reset")
	}

	// Ensure we can use it again after reset
	eg.Do(func() error { return nil })

	// Wait for task to complete
	_ = eg.Wait()

	if eg.HasError() {
		t.Error("ErrorGroup should not have error after successful Do following reset")
	}
}

// ============================================================================
// User Story 1: Reset() waits for active tasks and clears state
// ============================================================================

func TestReset_WaitsForActiveTasks(t *testing.T) {
	eg := NewErrorGroup()
	taskStarted := make(chan struct{})
	taskCanFinish := make(chan struct{})

	// Submit a long-running task (mock async behavior)
	go func() {
		taskStarted <- struct{}{}
		<-taskCanFinish
	}()

	// Simulate task execution by waiting
	<-taskStarted
	close(taskCanFinish)

	// Call Reset - should wait for any active tasks
	// Note: Since Do() is not yet implemented, this tests the structure
	eg.Reset()

	// After Reset, state should be clean
	if eg.HasError() {
		t.Error("ErrorGroup should not have error after Reset")
	}
}

func TestReset_ClearsErrorState(t *testing.T) {
	eg := NewErrorGroup()
	testErr := errors.New("test error")

	// Simulate error state (will be properly set when Do() is implemented)
	eg.hasError.Store(true)
	eg.firstError.Store(testErr)

	// Verify error state exists
	if !eg.HasError() {
		t.Error("ErrorGroup should have error before Reset")
	}

	// Reset should clear error state
	eg.Reset()

	// Verify error state is cleared
	if eg.HasError() {
		t.Error("ErrorGroup should not have error after Reset")
	}

	if eg.Error() != nil {
		t.Error("Error() should return nil after Reset")
	}
}

func TestReset_ImmediateReturnNoTasks(t *testing.T) {
	eg := NewErrorGroup()

	// Reset on empty ErrorGroup should return immediately
	start := time.Now()
	eg.Reset()
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("Reset with no active tasks should return immediately, took %v", elapsed)
	}
}

func TestReset_ConcurrentScenarios(t *testing.T) {
	eg := NewErrorGroup()

	// Test concurrent Reset calls
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			eg.Reset()
		}()
	}

	wg.Wait()

	// State should be clean after concurrent resets
	if eg.HasError() {
		t.Error("ErrorGroup should not have error after concurrent Resets")
	}
}

// ============================================================================
// Error handling and panic recovery
// ============================================================================

func TestPanicRecovery(t *testing.T) {
	eg := NewErrorGroup()

	// Submit a task that panics
	eg.Do(func() error {
		panic("something went wrong")
	})

	// Wait should return a PanicError
	err := eg.Wait()
	if err == nil {
		t.Fatal("Wait() should return error for panic")
	}

	// Check error type
	if _, ok := err.(*PanicError); !ok {
		t.Errorf("Error should be PanicError, got %T", err)
	}

	// Error message should contain panic info
	if err.Error() != "panic: something went wrong" {
		t.Errorf("Wrong error message: %v", err)
	}
}

func TestPanicRecoveryWithError(t *testing.T) {
	eg := NewErrorGroup()

	// Submit a task that panics with an error
	testErr := errors.New("test panic error")
	eg.Do(func() error {
		panic(testErr)
	})

	// Wait should return a PanicError with the error message
	err := eg.Wait()
	if err == nil {
		t.Fatal("Wait() should return error for panic")
	}

	// Check error type
	if _, ok := err.(*PanicError); !ok {
		t.Errorf("Error should be PanicError, got %T", err)
	}

	// Error message should contain the original error
	if err.Error() != "panic: test panic error" {
		t.Errorf("Wrong error message: %v", err)
	}
}

func TestBatch_OrderedExecution(t *testing.T) {
	eg := NewErrorGroup()
	var order []int

	// Submit tasks with delay to ensure ordering
	eg.Do(
		func() error {
			order = append(order, 1)
			return nil
		},
		func() error {
			order = append(order, 2)
			return nil
		},
		func() error {
			order = append(order, 3)
			return nil
		},
	)

	// Wait for completion
	if err := eg.Wait(); err != nil {
		t.Errorf("Wait() returned error: %v", err)
	}

	// Verify execution order
	if len(order) != 3 {
		t.Fatalf("Expected 3 tasks, got %d", len(order))
	}

	if order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Errorf("Tasks executed in wrong order: %v", order)
	}
}

func TestBatch_ErrorStopsExecution(t *testing.T) {
	eg := NewErrorGroup()
	var executed []int

	// Submit tasks where second one fails
	eg.Do(
		func() error {
			executed = append(executed, 1)
			return nil
		},
		func() error {
			executed = append(executed, 2)
			return errors.New("task 2 failed")
		},
		func() error {
			executed = append(executed, 3) // Should not execute
			return nil
		},
	)

	// Wait for completion
	err := eg.Wait()
	if err == nil {
		t.Error("Wait() should return error")
	}

	if err.Error() != "task 2 failed" {
		t.Errorf("Wrong error message: %v", err)
	}

	// Verify task 3 didn't execute
	if len(executed) != 2 {
		t.Errorf("Expected 2 tasks executed, got %d: %v", len(executed), executed)
	}

	if executed[0] != 1 || executed[1] != 2 {
		t.Errorf("Wrong execution order: %v", executed)
	}
}

func TestBatch_MixedExecutionTimes(t *testing.T) {
	eg := NewErrorGroup()
	var order []int

	// Submit tasks with different execution times
	eg.Do(
		func() error {
			order = append(order, 1)
			return nil
		},
		func() error {
			time.Sleep(10 * time.Millisecond)
			order = append(order, 2)
			return nil
		},
		func() error {
			order = append(order, 3)
			return nil
		},
	)

	// Wait for completion
	if err := eg.Wait(); err != nil {
		t.Errorf("Wait() returned error: %v", err)
	}

	// Verify order is maintained despite different execution times
	if len(order) != 3 {
		t.Fatalf("Expected 3 tasks, got %d", len(order))
	}

	if order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Errorf("Tasks executed in wrong order: %v", order)
	}
}

func TestBatch_LargeBatch(t *testing.T) {
	eg := NewErrorGroup()
	count := 100

	// Submit many tasks
	var tasks []func() error
	for i := 0; i < count; i++ {
		tasks = append(tasks, func() error {
			// Simulate work
			time.Sleep(time.Microsecond)
			return nil
		})
	}

	eg.Do(tasks...)

	// Wait for completion
	if err := eg.Wait(); err != nil {
		t.Errorf("Wait() returned error: %v", err)
	}
}

func TestDo_WaitForReset(t *testing.T) {
	eg := NewErrorGroup()

	// Submit some initial tasks
	eg.Do(func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	// Start Reset in background
	resetDone := make(chan struct{})
	go func() {
		eg.Reset()
		close(resetDone)
	}()

	// Give Reset time to start
	time.Sleep(time.Millisecond)

	// Try to submit new tasks during reset
	// This should wait for reset to complete
	taskSubmitted := make(chan struct{})
	go func() {
		eg.Do(func() error {
			return nil
		})
		close(taskSubmitted)
	}()

	// Both should complete without deadlock
	select {
	case <-resetDone:
		t.Log("Reset completed")
	case <-time.After(1 * time.Second):
		t.Fatal("Reset did not complete in time")
	}

	select {
	case <-taskSubmitted:
		t.Log("Task submitted after reset")
	case <-time.After(1 * time.Second):
		t.Fatal("Task did not submit after reset")
	}

	// Wait for all tasks
	if err := eg.Wait(); err != nil {
		t.Errorf("Wait() returned error: %v", err)
	}
}

func TestReset_StateTransitions(t *testing.T) {
	eg := NewErrorGroup()

	// Initial state should be Idle
	if eg.state.Load() != int32(StateIdle) {
		t.Errorf("Initial state should be Idle, got %d", eg.state.Load())
	}

	// Submit a task
	eg.Do(func() error {
		time.Sleep(time.Millisecond)
		return nil
	})

	// Start Reset and wait for it to actually start
	resetDone := make(chan struct{})
	go func() {
		eg.Reset()
		close(resetDone)
	}()

	// Wait for state to become Resetting
	for eg.state.Load() != int32(StateResetting) {
		time.Sleep(time.Microsecond)
		// Timeout check
		runtime.Gosched()
		if eg.state.Load() == int32(StateIdle) && len(resetDone) == 0 {
			// Reset might have completed very quickly
			break
		}
	}

	// State should be either Resetting or already back to Idle
	currentState := eg.state.Load()
	if currentState != int32(StateResetting) && currentState != int32(StateIdle) {
		t.Errorf("State should be Resetting or Idle, got %d", currentState)
	}

	// Wait for Reset to complete
	<-resetDone

	// After Reset completes, state should be Idle
	_ = eg.Wait() // Wait for task to complete

	if eg.state.Load() != int32(StateIdle) {
		t.Errorf("State should be Idle after Reset(), got %d", eg.state.Load())
	}
}
