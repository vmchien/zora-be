package validate

import "github.com/google/uuid"

func IsUuidValid(u *uuid.UUID) bool {
	return u != nil && *u != uuid.Nil
}
func IsUuidStringValid(u *string) bool {
	if u == nil || len(*u) == 0 {
		return false
	}
	if v, err := uuid.Parse(*u); err != nil || v == uuid.Nil {
		return false
	}
	return true
}
func IsUuidV7Valid(u *uuid.UUID) bool {
	return u != nil && *u != uuid.Nil && u.Version() == uuid.Version(7)
}
