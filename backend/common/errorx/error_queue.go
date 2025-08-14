package errorx

import (
	"sync"
)

// SyncError is a synchronized error type
type SyncError struct {
	mu   sync.Mutex
	once sync.Once
	err  error
}

// NewSyncError returns a new SyncError instance
func NewSyncError() *SyncError {
	return &SyncError{}
}

// Do executes the given functions in a synchronized manner
func (se *SyncError) Do(fns ...func() error) {
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
func (se *SyncError) Error() error {
	se.mu.Lock()
	defer se.mu.Unlock()
	return se.err
}

// HasError returns true if an error has been recorded
func (se *SyncError) HasError() bool {
	se.mu.Lock()
	defer se.mu.Unlock()
	return se.err != nil
}

// Reset resets the error
func (se *SyncError) Reset() {
	se.mu.Lock()
	defer se.mu.Unlock()
	se.once = sync.Once{}
	se.err = nil
}
