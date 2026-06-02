package ent_mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
	"vn.vato.zora.be.api/pkg/constant"
	"vn.vato.zora.be.api/pkg/data/ent_annotation"
	"vn.vato.zora.be.api/pkg/def"
)

// MixinTenant holds the schema definition for tenant-related entities.
type MixinTenant struct {
	mixin.Schema
}

func (MixinTenant) Fields() []ent.Field {
	return []ent.Field{
		field.UUID(constant.TenantIDField, uuid.UUID{}).
			Default(def.DefaultTenantID).
			Optional().
			Annotations(ent_annotation.TenantAnnotation{}).
			Immutable(),
	}
}
func (MixinTenant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(constant.TenantIDField),
	}
}
func (MixinTenant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		ent_annotation.TenantAnnotation{},
	}
}
