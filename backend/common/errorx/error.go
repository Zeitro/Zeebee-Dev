package errorx

import (
	"bytes"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

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

// GetErrorStack get error stack
func GetErrorStack(err error) []string {
	if e, ok := err.(interface {
		StackTrace() errors.StackTrace
	}); ok {
		st := e.StackTrace()
		var lines []string
		for _, frame := range st {
			line, err := frame.MarshalText()
			if err != nil {
				continue
			}
			lines = append(lines, string(line))
		}
		return lines
	}

	// fallback: debug.Stack()
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	// replace \n\t with \n
	buf = bytes.ReplaceAll(buf, []byte("\n\t"), []byte("\n"))
	// split lines with \n
	lines := bytes.Split(buf[:n], []byte("\n"))
	lineArr := []string{}
	for i := 0; i < len(lines); i++ {
		if i == 0 {
			continue
		}
		lineArr = append(lineArr, string(lines[i]))
	}
	// to string
	return lineArr
}

func New(text string) (err error) {
	return WithStack(errors.New(text))
}

func Errorf(format string, args ...interface{}) (err error) {
	return errors.Errorf(format, args...)
}

func Wrap(err error, text string) (newErr error) {
	return errors.Wrap(err, text)
}
func Wrapf(err error, format string, args ...interface{}) (newErr error) {
	return errors.Wrapf(err, format, args...)
}
func Unwrap(err error) (newErr error) {
	if err == nil {
		return nil
	}
	if u := errors.Unwrap(err); u != nil {
		return u
	}
	return errors.Cause(err)
}

func Cause(err error) (newErr error) {
	return errors.Cause(err)
}

func WithStack(err error) (newErr error) {
	return errors.WithStack(err)
}

func NotFound(record string) (err error) {
	if len(record) == 0 {
		return New("Record Not Found")
	}
	return Errorf("%s Not Found", record)
}
