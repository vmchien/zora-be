package validate

import (
	"fmt"
	"strings"
	"unicode"
)

func ValidatePassword(password string, requireUppercase, requireLowercase, requireNumber, requireSymbol bool, minLength int) error {
	var errs []string

	// Length check
	if len(password) < minLength {
		errs = append(errs, fmt.Sprintf("password must be at least %d characters long", minLength))
	}

	var hasUpper, hasLower, hasNumber, hasSymbol bool
	// Character category checks
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}

	if requireUppercase && !hasUpper {
		errs = append(errs, "password must contain at least one uppercase letter")
	}
	if requireLowercase && !hasLower {
		errs = append(errs, "password must contain at least one lowercase letter")
	}
	if requireNumber && !hasNumber {
		errs = append(errs, "password must contain at least one number")
	}
	if requireSymbol && !hasSymbol {
		errs = append(errs, "password must contain at least one special symbol")
	}

	// TODO: Implement history check (prevent reuse of last N passwords)
	// You may need to query your storage for the user's previous passwords and compare hashes.

	// Aggregate errors
	if len(errs) > 0 {
		return fmt.Errorf("password validation failed: %s", strings.Join(errs, "; "))
	}

	return nil
}
