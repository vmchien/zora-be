package ent_mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"vn.vato.zora.be.api/pkg/constant"
	"vn.vato.zora.be.api/pkg/data/ent_annotation"
	"vn.vato.zora.be.api/pkg/def"
)

type MixinTags struct {
	mixin.Schema
}

func (MixinTags) Fields() []ent.Field {
	return []ent.Field{
		field.JSON(constant.TagsField, def.JsonArrayString).
			// SchemaType(map[string]string{
			// 	dialect.Postgres: def.JSONB,
			// }).
			Annotations(ent_annotation.TagsAnnotation{}).
			Optional().
			Default(def.DefaultArrayString()).
			Comment("Additional tags stored as JSONB"),
	}
}
func (MixinTags) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(constant.TagsField).
			Annotations(
				entsql.IndexTypes(map[string]string{
					dialect.MySQL:    def.MySqlJsonIndexType,
					dialect.Postgres: def.PostgresJsonIndexType,
				}),
			),
	}
}
