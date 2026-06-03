package ent_mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
	"vn.vato.zora.be.api/pkg/constant"
	"vn.vato.zora.be.api/pkg/data/ent_annotation"
)

// MixinAuditUser holds the schema definition for entities that require audit user information.
type MixinAuditUser struct {
	mixin.Schema
}

func (MixinAuditUser) Fields() []ent.Field {
	return []ent.Field{
		field.UUID(constant.CreatedByField, uuid.UUID{}).
			Immutable().
			Optional().
			Annotations(ent_annotation.AuditUserAnnotation{}).
			Comment("ID of the employee who created this record - for audit trail"),

		field.UUID(constant.UpdatedByField, uuid.UUID{}).
			Optional().
			Nillable().
			Annotations(ent_annotation.AuditUserAnnotation{}).
			Comment("ID of the employee who last updated this record - for audit trail"),
	}
}

func (MixinAuditUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		ent_annotation.AuditUserAnnotation{},
	}
}
