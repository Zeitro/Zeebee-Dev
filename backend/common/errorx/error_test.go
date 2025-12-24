package errorx

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetErrorMsg(t *testing.T) {
	// Test case: nil error
	if msg := GetErrorMsg(nil); msg != "" {
		t.Errorf("GetErrorMsg(nil) = %q, want empty string", msg)
	}

	// Test case: regular error
	err := errors.New("test error")
	if msg := GetErrorMsg(err); msg != "test error" {
		t.Errorf("GetErrorMsg(err) = %q, want 'test error'", msg)
	}

	// Test case: gRPC error
	grpcErr := status.New(codes.NotFound, "not found").Err()
	if msg := GetErrorMsg(grpcErr); msg != "not found" {
		t.Errorf("GetErrorMsg(grpcErr) = %q, want 'not found'", msg)
	}

	// Test case: gRPC error with rpc error prefix
	grpcErrWithPrefix := status.New(codes.Unknown, "rpc error: code = Unknown desc = something went wrong").Err()
	expected := "something went wrong"
	if msg := GetErrorMsg(grpcErrWithPrefix); msg != expected {
		t.Errorf("GetErrorMsg(grpcErrWithPrefix) = %q, want %q", msg, expected)
	}
}

func TestGetErrorStack(t *testing.T) {
	// Test case: error with stack trace
	err := errors.New("test error")
	errWithStack := WithStack(err)
	stack := GetErrorStack(errWithStack)

	if len(stack) == 0 {
		t.Error("GetErrorStack returned empty stack for error with stack trace")
	}

	// Check that the stack contains relevant information
	stackStr := strings.Join(stack, "\n")
	if !strings.Contains(stackStr, "TestGetErrorStack") {
		t.Error("Stack trace should contain test function name")
	}

	// Test case: regular error without stack trace
	regErr := errors.New("regular error")
	regStack := GetErrorStack(regErr)
	if len(regStack) == 0 {
		t.Error("GetErrorStack should return fallback stack for regular errors")
	}
}

func TestNew(t *testing.T) {
	err := New("test error")
	if err == nil {
		t.Fatal("New returned nil")
	}

	if err.Error() != "test error" {
		t.Errorf("err.Error() = %q, want 'test error'", err.Error())
	}

	// Check that the error has a stack trace
	stack := GetErrorStack(err)
	if len(stack) == 0 {
		t.Error("New should create error with stack trace")
	}
	t.Logf("%v\n", err)
}

func TestErrorf(t *testing.T) {
	err := Errorf("error %d: %s", 42, "test")
	if err == nil {
		t.Fatal("Errorf returned nil")
	}

	expected := "error 42: test"
	if err.Error() != expected {
		t.Errorf("err.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestWrap(t *testing.T) {
	cause := New("cause error")
	err := Wrap(cause, "wrapper")

	if err == nil {
		t.Fatal("Wrap returned nil")
	}

	if !strings.Contains(err.Error(), "cause error") {
		t.Errorf("err.Error() = %q, should contain 'cause error'", err.Error())
	}

	if !strings.Contains(err.Error(), "wrapper") {
		t.Errorf("err.Error() = %q, should contain 'wrapper'", err.Error())
	}
}

func TestWrapf(t *testing.T) {
	cause := New("cause error")
	err := Wrapf(cause, "wrapper %d", 42)

	if err == nil {
		t.Fatal("Wrapf returned nil")
	}

	if !strings.Contains(err.Error(), "cause error") {
		t.Errorf("err.Error() = %q, should contain 'cause error'", err.Error())
	}

	if !strings.Contains(err.Error(), "wrapper 42") {
		t.Errorf("err.Error() = %q, should contain 'wrapper 42'", err.Error())
	}

}

func TestUnwrap(t *testing.T) {
	cause := New("cause error")
	wrapped := Wrap(cause, "wrapper")

	unwrapped := Unwrap(wrapped)
	if unwrapped == nil {
		t.Fatal("expected Unwrap to return a non-nil error")
	}
	if errors.Cause(unwrapped).Error() != cause.Error() {
		t.Errorf("expected cause message %q, got %q",
			cause.Error(), errors.Cause(unwrapped).Error())
	}
}

func TestWithStack(t *testing.T) {
	err := errors.New("test error")
	errWithStack := WithStack(err)

	if errWithStack == nil {
		t.Fatal("WithStack returned nil")
	}

	if errWithStack.Error() != "test error" {
		t.Errorf("errWithStack.Error() = %q, want 'test error'", errWithStack.Error())
	}

	// Check that the error has a stack trace
	stack := GetErrorStack(errWithStack)
	if len(stack) == 0 {
		t.Error("WithStack should add stack trace to error")
	}
}

func TestNotFound(t *testing.T) {
	// Test case: empty record name
	err := NotFound("")
	if err == nil {
		t.Fatal("NotFound returned nil")
	}

	expected := "Record Not Found"
	if err.Error() != expected {
		t.Errorf("err.Error() = %q, want %q", err.Error(), expected)
	}

	// Test case: specific record name
	recordName := "User"
	err = NotFound(recordName)
	if err == nil {
		t.Fatal("NotFound returned nil")
	}

	expected = fmt.Sprintf("%s Not Found", recordName)
	if err.Error() != expected {
		t.Errorf("err.Error() = %q, want %q", err.Error(), expected)
	}
}
