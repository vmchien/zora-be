package ent_annotation

import (
	"entgo.io/ent/schema"
)

type AuditTimeAnnotation struct{}

var _ schema.Annotation = (*AuditTimeAnnotation)(nil)

func (T AuditTimeAnnotation) Name() string {
	return AuditTimeAnnotationName
}
