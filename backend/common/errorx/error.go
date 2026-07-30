// Package errorx wraps errors with stack traces and provides go-zero/gRPC
// error helpers. Wrap/New capture the caller's stack so GetErrorStack can
// report it for logging.
package errorx

import (
	"bytes"
	"errors"
	stderrors "errors"
	"fmt"
	"runtime"
	"strings"

	"google.golang.org/grpc/status"
)

// stackDepth is the max number of frames captured per error.
const stackDepth = 32

// Error is one layer of a wrapped error chain. It holds the message added at
// this layer, the inner cause, the stack captured where it was created, and a
// flag selecting how Error() renders:
//   - inlined == false (Wrap/Wrapf): "msg: cause.Error()"
//   - inlined == true (WithStack): "msg" verbatim (msg already holds cause text)
type Error struct {
	message string
	cause   error
	pcs     []uintptr
	inlined bool
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.cause != nil && !e.inlined {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

// Unwrap returns the inner cause, supporting errors.Is/As and chain walking.
// Nil-safe so stdlib errors.Unwrap never panics on a typed-nil *Error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// Stack renders this layer's captured call stack. Returns "" if no stack.
func (e *Error) Stack() string {
	if e == nil || len(e.pcs) == 0 {
		return ""
	}
	return renderStack(e.pcs)
}

// newErr builds an *Error and captures the caller's stack. Frame layout at
// capture: 0=runtime.Callers, 1=newErr, 2=New/Wrap, 3=user caller — so skip=3.
func newErr(message string, cause error) *Error {
	pcs := make([]uintptr, stackDepth)
	n := runtime.Callers(3, pcs)
	return &Error{message: message, cause: cause, pcs: pcs[:n]}
}

// New returns a leaf error with no cause. Captures the caller's stack.
func New(text string) error { return newErr(text, nil) }

// Errorf formats according to a printf-style format specifier and returns a
// leaf error. Captures the caller's stack.
func Errorf(format string, args ...any) error {
	return newErr(fmt.Sprintf(format, args...), nil)
}

// Wrap adds a message layer around err and captures the caller's stack.
// Returns nil if err is nil.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return newErr(msg, err)
}

// Wrapf adds a formatted message layer around err and captures the caller's
// stack. Returns nil if err is nil.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return newErr(fmt.Sprintf(format, args...), err)
}

// Unwrap returns the inner cause of err, or nil if err is a leaf.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// Cause returns the innermost cause of err by walking its Unwrap chain.
func Cause(err error) error {
	for err != nil {
		u := stderrors.Unwrap(err)
		if u == nil {
			return err
		}
		err = u
	}
	return nil
}

func GetErrorMsg(err error) string {
	if err == nil {
		return ""
	}

	// Try to extract the original gRPC error message
	if s, ok := status.FromError(err); ok {
		msg := s.Message()
		if strings.HasPrefix(msg, "rpc error: code = ") {
			msgArr := strings.Split(msg, "=")
			if len(msgArr) > 2 {
				msg = strings.Join(msgArr[2:], "=")
			}
		}
		msg = strings.TrimSpace(msg)
		return msg
	}

	// Fallback to the original error message
	return err.Error()
}

// GetErrorStack get error stack as a slice of "function\n\tfile:line" lines.
// Reads the stack captured on *Error; falls back to runtime.Stack() for
// foreign errors without one.
func GetErrorStack(err error) []string {
	if err == nil {
		return nil
	}

	// 1. Our *Error: render the stack captured at construction.
	var our *Error
	if errors.As(err, &our) && our != nil {
		if s := our.Stack(); s != "" {
			return splitStackLines(s)
		}
	}

	// 2. fallback: runtime.Stack() for foreign errors without a captured stack.
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	// replace \n\t with \n
	buf = bytes.ReplaceAll(buf, []byte("\n\t"), []byte("\n"))
	// split lines with \n
	lines := bytes.Split(buf[:n], []byte("\n"))
	lineArr := []string{}
	for i := range lines {
		if i == 0 {
			continue
		}
		lineArr = append(lineArr, string(lines[i]))
	}
	return lineArr
}

// splitStackLines splits a rendered stack block into non-empty lines.
func splitStackLines(s string) []string {
	var out []string
	for line := range strings.SplitSeq(s, "\n") {
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

// renderStack formats program counters as the standard
// "function\n\tfile:line" listing.
func renderStack(pcs []uintptr) string {
	if len(pcs) == 0 {
		return ""
	}
	var b strings.Builder
	frames := runtime.CallersFrames(pcs)
	for {
		f, more := frames.Next()
		if f.Function != "runtime.goexit" {
			fmt.Fprintf(&b, "%s\n\t%s:%d\n", f.Function, f.File, f.Line)
		}
		if !more {
			break
		}
	}
	return b.String()
}

// WithStack attaches a stack trace captured at the call site to err,
// preserving err's message and cause chain. Returns nil if err is nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	e := newErr(err.Error(), err)
	e.inlined = true
	return e
}

// NotFound returns a not-found error for record, or a generic one when record
// is empty.
func NotFound(record string) error {
	if len(record) == 0 {
		return New("record not found")
	}
	return Errorf("%s not found", record)
}
