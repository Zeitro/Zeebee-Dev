package stringx

import (
	"fmt"
	"strings"

	"github.com/Rhymond/go-money"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Format(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// ParseUserName parse user name
func ParseUserName(fullName string) (firstName, lastName string) {
	fullName = strings.TrimSpace(fullName)
	nameArr := strings.Split(fullName, " ")
	if len(nameArr) == 0 {
		return "", ""
	}
	return nameArr[0], nameArr[len(nameArr)-1]
}

// FormatFullName
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

func FormatUsPhoneNumber(phoneNumber string) (newPhoneNumber string) {
	phoneNumber = strings.TrimSpace(phoneNumber)
	if len(phoneNumber) != 10 {
		return phoneNumber
	}
	newPhoneNumber = fmt.Sprintf("(%s) %s-%s", phoneNumber[0:3], phoneNumber[3:6], phoneNumber[6:10])
	return newPhoneNumber
}

func FormatUsSsn(ssn string) (newSsn string) {
	ssn = strings.TrimSpace(ssn)
	if len(ssn) != 9 {
		return ssn
	}
	newSsn = fmt.Sprintf("%s-%s-%s", ssn[0:3], ssn[3:5], ssn[5:9])
	return newSsn
}

func FormatCurrency(symbol string, amount float64) string {
	if amount == 0 {
		return ""
	}
	return fmt.Sprintf("%v%.2f", symbol, amount)
}

func FormatUsDollarCurrency(amount float64) string {
	usdAmount := money.NewFromFloat(amount, money.USD)
	return strings.TrimSuffix(usdAmount.Display(), ".00")
}

func ToTitle(s string) string {
	caser := cases.Title(language.Und)
	return caser.String(s)
}
