package errorx

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ExecutionState represents the current execution state of ErrorGroup
type ExecutionState int32

const (
	StateIdle      ExecutionState = iota // Idle, no active tasks
	StateExecuting                       // Executing, has active tasks
	StateResetting                       // Resetting, waiting for tasks to complete
)

// taskBatch represents a batch of tasks submitted in a single Do() call
type taskBatch struct {
	tasks []func() error // Task functions (preserve order)
	done  chan struct{}  // Batch completion notification (optional)
}

// ErrorGroup provides synchronized error handling for concurrent operations
// with guaranteed task ordering and fast-fail error propagation.
type ErrorGroup struct {
	// Execution state (atomic operations)
	state atomic.Int32

	// Task queue (buffered channel, capacity 1000)
	taskQueue chan taskBatch

	// Active task counter
	activeTasks sync.WaitGroup

	// Error state (atomic operations)
	hasError   atomic.Bool  // Error flag (fast check)
	firstError atomic.Value // Stores first error encountered

	// Reset coordination
	resetMu   sync.Mutex
	resetCond *sync.Cond

	// Startup flag (ensures dispatcher starts once)
	started atomic.Bool
}

// NewErrorGroup creates a new ErrorGroup instance in idle state.
func NewErrorGroup() *ErrorGroup {
	eg := &ErrorGroup{
		taskQueue: make(chan taskBatch, 1000),
	}
	eg.resetCond = sync.NewCond(&eg.resetMu)
	return eg
}

// Do submits tasks for asynchronous execution in order.
//
// The method returns immediately without waiting for tasks to complete.
// Tasks within the batch execute in submission order.
// If an error has already been recorded, new tasks are not executed.
func (eg *ErrorGroup) Do(fns ...func() error) {
	// Fast-fail: if error already recorded, don't execute new tasks
	if eg.HasError() {
		return
	}

	// Wait if Reset() is in progress
	// This prevents submitting tasks during reset operation
	eg.resetMu.Lock()
	defer eg.resetMu.Unlock()
	for eg.state.Load() == int32(StateResetting) {
		eg.resetCond.Wait()
	}

	// Check for empty batch
	if len(fns) == 0 {
		return
	}

	// Filter out nil functions
	var validTasks []func() error
	for _, fn := range fns {
		if fn != nil {
			validTasks = append(validTasks, fn)
		}
	}

	if len(validTasks) == 0 {
		return
	}

	// Track tasks in activeTasks BEFORE sending to queue
	// This ensures Wait() will block even if dispatcher hasn't started yet
	eg.activeTasks.Add(len(validTasks))

	// Start dispatcher goroutine on first Do() call
	if eg.started.CompareAndSwap(false, true) {
		go eg.dispatcher()
	}

	// Create batch and send to queue
	batch := taskBatch{
		tasks: validTasks,
		done:  nil,
	}

	// Non-blocking send (buffered channel)
	eg.taskQueue <- batch
}

// Wait blocks until all tasks are completed or an error occurs.
//
// Returns the first error if any, nil otherwise.
// If an error has already been recorded, returns immediately.
func (eg *ErrorGroup) Wait() error {
	// Fast path: if error already set, return immediately
	if eg.HasError() {
		return eg.Error()
	}

	// Wait for all active tasks to complete
	eg.activeTasks.Wait()

	// After tasks complete, check error one more time
	if eg.HasError() {
		return eg.Error()
	}

	return nil
}

// Reset waits for all active tasks to complete, then clears error state.
//
// This method blocks until all active tasks finish, then resets the
// ErrorGroup to initial state. New tasks submitted during reset
// will wait until reset completes.
func (eg *ErrorGroup) Reset() {
	eg.resetMu.Lock()
	defer eg.resetMu.Unlock()

	// Set state to Resetting before waiting
	// This signals Do() to wait for reset to complete
	eg.state.Store(int32(StateResetting))

	// Wait for all active tasks to complete
	eg.activeTasks.Wait()

	// Clear error state - firstError is kept but hasError flag controls access
	eg.hasError.Store(false)

	// Reset execution state to Idle
	eg.state.Store(int32(StateIdle))

	// Wake up all goroutines waiting for reset to complete
	eg.resetCond.Broadcast()
}

// Error returns the first error encountered (non-blocking).
//
// Returns:
//   - error: The first error encountered, or nil if no error
//
// Thread Safety: Safe to call concurrently (lock-free atomic operation)
// Performance: ~10ns
//
// Example:
//
//	if err := eg.Error(); err != nil {
//	    log.Printf("Error: %v", err)
//	}
func (eg *ErrorGroup) Error() error {
	if !eg.HasError() {
		return nil
	}
	err, _ := eg.firstError.Load().(error)
	return err
}

// HasError returns true if an error has been recorded (non-blocking).
//
// This is a fast check using atomic operations.
//
// Returns:
//   - bool: true if an error has been recorded, false otherwise
//
// Thread Safety: Safe to call concurrently (lock-free atomic operation)
// Performance: ~5ns
//
// Example:
//
//	if eg.HasError() {
//	    log.Printf("Error: %v", eg.Error())
//	}
func (eg *ErrorGroup) HasError() bool {
	return eg.hasError.Load()
}

// setError sets the first error using atomic CAS (compare-and-swap)
func (eg *ErrorGroup) setError(err error) {
	if eg.hasError.CompareAndSwap(false, true) {
		eg.firstError.Store(err)
	}
}

// dispatcher processes tasks from the queue in order
func (eg *ErrorGroup) dispatcher() {
	for batch := range eg.taskQueue {
		// Execute tasks in batch sequentially
		for i, task := range batch.tasks {
			// Check if error already set (fast-fail)
			if eg.HasError() {
				// Still need to call Done for this and all remaining tasks in batch
				for j := i; j < len(batch.tasks); j++ {
					eg.activeTasks.Done()
				}
				break // Exit task loop
			}

			eg.state.Store(int32(StateExecuting))

			// Execute task (with panic recovery)
			err := eg.runTask(task)

			// Task completed - mark as done
			eg.activeTasks.Done()

			// Handle error
			if err != nil {
				eg.setError(err)
				// Call Done for remaining tasks in this batch
				for j := i + 1; j < len(batch.tasks); j++ {
					eg.activeTasks.Done()
				}
				break // Skip remaining tasks in batch
			}
		}

		// Mark batch completion if needed
		if batch.done != nil {
			close(batch.done)
		}
	}

	// All tasks completed
	eg.state.Store(int32(StateIdle))
}

// runTask executes a single task with panic recovery
func (eg *ErrorGroup) runTask(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert panic to error
			err = &PanicError{value: r}
		}
	}()
	return fn()
}

// PanicError represents a recovered panic
type PanicError struct {
	value interface{}
}

func (e *PanicError) Error() string {
	return "panic: " + stringifyPanic(e.value)
}

func stringifyPanic(v interface{}) string {
	switch v := v.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
