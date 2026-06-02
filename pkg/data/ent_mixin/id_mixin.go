package ent_mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
	"vn.vato.zora.be.api/pkg/guid"
)

type MixinID struct {
	mixin.Schema
}

func (MixinID) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(guid.New).
			Immutable().
			Comment("Primary key"),
	}
}
