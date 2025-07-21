package stringx

import (
	"fmt"
)

func Formate(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
