package ent_mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"vn.vato.zora.be.api/pkg/constant"
)

// MixinTSVector provides full-text search capabilities for entities.
type MixinTSVector struct {
	mixin.Schema
}

func (MixinTSVector) Fields() []ent.Field {
	return []ent.Field{
		field.String(constant.TSVectorField).
			SchemaType(map[string]string{
				dialect.Postgres: "tsvector",
			}).
			Optional().
			Nillable().
			Immutable().
			Sensitive().
			Annotations(entsql.Skip()).
			Comment("Full-text search vector for tenant fields"),
	}
}
