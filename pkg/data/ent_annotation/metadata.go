package ent_annotation

import (
	"entgo.io/ent/schema"
)

type MetadataAnnotation struct{}

// Annotations returns the annotations for the MetadataAnnotation type.
var _ schema.Annotation = (*MetadataAnnotation)(nil)

func (T MetadataAnnotation) Name() string {
	return MetadataAnnotationName
}
