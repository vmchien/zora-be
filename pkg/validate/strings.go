package validate

import "strings"

func IsStringValid(s *string) bool {
	return s != nil && *s != ""
}

func IsStringEmpty(val any) bool {
	if s, ok := val.(string); ok {
		return s == ""
	}
	return false
}
func IsStringBlank(val any) bool {
	if s, ok := val.(string); ok {
		return strings.TrimSpace(s) == ""
	}
	return false
}
