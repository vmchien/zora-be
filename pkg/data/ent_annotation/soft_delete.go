package ent_annotation

import "entgo.io/ent/schema"

type SoftDeleteAnnotation struct{}

var _ schema.Annotation = (*SoftDeleteAnnotation)(nil)

func (a SoftDeleteAnnotation) Name() string {
	return SoftDeleteAnnotationName
}
