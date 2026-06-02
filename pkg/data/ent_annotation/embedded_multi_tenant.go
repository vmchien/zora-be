package ent_annotation

import "entgo.io/ent/schema"

type TenantAnnotation struct{}

var _ schema.Annotation = (*TenantAnnotation)(nil)

func (a TenantAnnotation) Name() string {
	return TenantAnnotationName
}
