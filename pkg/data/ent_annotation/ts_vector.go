package ent_annotation

import (
	"entgo.io/ent/schema"
)

type TSVectorAnnotation struct{}

// Annotations returns the annotations for the TSVectorAnnotation type.
var _ schema.Annotation = (*TSVectorAnnotation)(nil)

func (T TSVectorAnnotation) Name() string {
	return TSVectorAnnotationName
}
