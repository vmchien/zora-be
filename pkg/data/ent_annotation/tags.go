package ent_annotation

import (
	"entgo.io/ent/schema"
)

type TagsAnnotation struct{}

var _ schema.Annotation = (*TagsAnnotation)(nil)

func (T TagsAnnotation) Name() string {
	return TagsAnnotationName
}
