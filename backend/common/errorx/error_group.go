package errorx

import (
	"sync"
)

// ErrorGroup is a synchronized error type
type ErrorGroup struct {
	mu   sync.Mutex
	once sync.Once
	err  error
}

// NewSyncError returns a new SyncError instance
func NewErrorGroup() *ErrorGroup {
	return &ErrorGroup{}
}

// Do executes the given functions in a synchronized manner
func (se *ErrorGroup) Do(fns ...func() error) {
	if se.HasError() {
		return
	}

	se.mu.Lock()
	defer se.mu.Unlock()

	if se.err != nil {
		return
	}

	for _, fn := range fns {
		err := fn()
		if err != nil {
			se.once.Do(func() {
				se.err = err
			})
			return
		}
	}
}

// Error returns the first error that was encountered
func (se *ErrorGroup) Error() error {
	se.mu.Lock()
	defer se.mu.Unlock()
	return se.err
}

// HasError returns true if an error has been recorded
func (se *ErrorGroup) HasError() bool {
	se.mu.Lock()
	defer se.mu.Unlock()
	return se.err != nil
}

// Reset resets the error
func (se *ErrorGroup) Reset() {
	se.mu.Lock()
	defer se.mu.Unlock()
	se.once = sync.Once{}
	se.err = nil
}
