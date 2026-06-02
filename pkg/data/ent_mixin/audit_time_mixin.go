package ent_mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"vn.vato.zora.be.api/pkg/constant"
	"vn.vato.zora.be.api/pkg/data/ent_annotation"
)

// MixinAuditTime holds the schema definition for entities that require audit timestamps.
type MixinAuditTime struct {
	mixin.Schema
}

func (MixinAuditTime) Fields() []ent.Field {
	return []ent.Field{
		field.Time(constant.CreatedAtField).
			Default(time.Now).
			Immutable().
			Annotations(ent_annotation.AuditTimeAnnotation{}).
			Comment("Record creation timestamp"),

		field.Time(constant.UpdatedAtField).
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(ent_annotation.AuditTimeAnnotation{}).
			Comment("Record last update timestamp"),
	}
}
