package utils

import (
	"regexp"
	"strings"
)

// EscapeTSQuery cleans up a string for use in a PostgreSQL full-text search query.
func EscapeTSQuery(input string) string {
	// Remove all characters that are not alphanumeric or whitespace
	allowed := regexp.MustCompile(`[^a-zA-Z0-9\s\|\&\:\*\!']+`)
	cleaned := allowed.ReplaceAllString(input, "")

	// Escape special characters
	cleaned = strings.ReplaceAll(cleaned, "'", "''")

	return strings.TrimSpace(cleaned)
}
