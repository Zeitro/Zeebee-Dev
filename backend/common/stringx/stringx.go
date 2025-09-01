package stringx

import (
	"fmt"
	"strings"
)

func Formate(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func ParseUserName(fullName string) (firstName, lastName string) {
	fullName = strings.TrimSpace(fullName)
	nameArr := strings.Split(fullName, " ")
	if len(nameArr) == 0 {
		return "", ""
	}
	return nameArr[0], nameArr[len(nameArr)-1]
}

func FormatFullName(firstName, midName, lastName string) string {
	names := make([]string, 0, 3)
	if firstName != "" {
		names = append(names, firstName)
	}
	if midName != "" {
		names = append(names, midName)
	}
	if lastName != "" {
		names = append(names, lastName)
	}
	return strings.Join(names, " ")
}

func TransNumberToString[T ~int | ~int32 | ~int64 | ~int8 | ~uint | ~float64 | ~float32 | ~uint8 | ~uint16 | ~uint32 | ~uint64](num T, zeroToEmpty bool) string {
	if num == 0 && zeroToEmpty {
		return ""
	}
	return fmt.Sprintf("%v", num)
}
