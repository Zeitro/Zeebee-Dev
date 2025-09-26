package errorx

import (
	"bytes"
	"fmt"
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
			lines = append(lines, fmt.Sprintf("%s", line))
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
