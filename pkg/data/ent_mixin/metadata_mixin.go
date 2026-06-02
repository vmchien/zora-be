package ent_mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"vn.vato.zora.be.api/pkg/data/ent_annotation"
	"vn.vato.zora.be.api/pkg/def"
)

type MixinMetadata struct {
	mixin.Schema
}

func (MixinMetadata) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("metadata", def.JsonMap).
			// SchemaType(map[string]string{
			// 	dialect.Postgres: def.JSONB,
			// }).
			Optional().
			Annotations(ent_annotation.MetadataAnnotation{}).
			Comment("Additional metadata stored as JSONB"),
	}
}
func (MixinMetadata) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("metadata").
			Annotations(
				entsql.IndexTypes(map[string]string{
					dialect.MySQL:    def.MySqlJsonIndexType,
					dialect.Postgres: def.PostgresJsonIndexType,
				}),
			),
	}
}

func (MixinMetadata) Annotations() []schema.Annotation {
	return []schema.Annotation{
		ent_annotation.MetadataAnnotation{},
	}
}
