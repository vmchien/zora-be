package ent_schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Balance struct {
	ent.Schema
}

func (Balance) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id").Unique(),
		field.Float("balance").Default(0),
		field.Time("created_at").Immutable(),
		field.Time("updated_at"),
	}
}

func (Balance) Edges() []ent.Edge {
	return nil
}
