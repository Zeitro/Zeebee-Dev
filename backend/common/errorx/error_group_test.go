package errorx

import (
	"errors"
	"testing"
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

	if !eg.HasError() {
		t.Error("ErrorGroup should have an error after error function")
	}

	if eg.Error() != testErr {
		t.Errorf("ErrorGroup.Error() = %v, want %v", eg.Error(), testErr)
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

	// Second Do should not execute since first Do already ran
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

	if !eg.HasError() {
		t.Error("ErrorGroup should have an error after Do with error function")
	}
}

func TestErrorGroupReset(t *testing.T) {
	eg := NewErrorGroup()
	testErr := errors.New("test error")

	eg.Do(func() error { return testErr })

	if !eg.HasError() {
		t.Error("ErrorGroup should have an error")
	}

	eg.Reset()

	if eg.HasError() {
		t.Error("ErrorGroup should not have an error after reset")
	}

	// Ensure we can use it again after reset
	eg.Do(func() error { return nil })

	if eg.HasError() {
		t.Error("ErrorGroup should not have error after successful Do following reset")
	}
}
