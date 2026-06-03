package cast

import (
	"regexp"
	"strings"
	"unsafe"
)

// CamelToSnake converts a CamelCase string to snake_case.
// It replaces the first uppercase letter followed by lowercase letters with an underscore and the second part,
// and replaces all uppercase letters that follow lowercase letters with an underscore before them.
// For example, "CamelCase" becomes "camel_case", and "HTTPRequest" becomes "http_request".
// It also converts the entire string to lowercase.
// This function is useful for converting identifiers in programming languages that use snake_case naming conventions.
// It uses regular expressions to find patterns in the string and replace them accordingly.
func CamelToSnake(s string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func PascalToKebab(input string) string {
	re1 := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	input = re1.ReplaceAllString(input, "$1-$2")

	re2 := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	input = re2.ReplaceAllString(input, "$1-$2")

	re3 := regexp.MustCompile(`([a-zA-Z])([0-9])`)
	input = re3.ReplaceAllString(input, "$1-$2")

	re4 := regexp.MustCompile(`([0-9])([a-zA-Z])`)
	input = re4.ReplaceAllString(input, "$1-$2")

	return strings.ToLower(input)
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

func SafeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
