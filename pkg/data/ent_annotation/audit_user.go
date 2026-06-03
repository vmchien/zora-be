package ent_annotation

import (
	"entgo.io/ent/schema"
)

type AuditUserAnnotation struct{}

var _ schema.Annotation = (*AuditUserAnnotation)(nil)

func (T AuditUserAnnotation) Name() string {
	return AuditUserAnnotationName
}
