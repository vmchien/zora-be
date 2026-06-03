package def

import (
	"github.com/google/uuid"
	"vn.vato.zora.be.api/pkg/constant"
)

func DefaultStruct() struct{} {
	return struct{}{}
}

func DefaultJsonMap() map[string]any {
	return map[string]any{}
}
func DefaultJsonArray() []map[string]any {
	return []map[string]any{}
}
func DefaultArray() []any {
	return []any{}
}
func DefaultArrayString() []string {
	return []string{}
}
func DefaultTenantID() uuid.UUID {
	return constant.DefaultTenantID
}
