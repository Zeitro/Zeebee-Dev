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
