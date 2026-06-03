package validate

import (
	"regexp"

	"vn.vato.zora.be.api/pkg/constant"
)

func IsEmailValid(email string) bool {
	if len(email) == 0 {
		return false
	}
	if ok, err := regexp.MatchString(constant.REGEXP_EMAIL, email); err == nil {
		return ok
	}
	return false
}
