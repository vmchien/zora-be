package ent_schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type AccountLogin struct {
	ent.Schema
}

func (AccountLogin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("account").
			Unique().
			NotEmpty(),
		field.String("account_type").
			NotEmpty(),
		field.String("password_salt").
			NotEmpty(),
		field.String("password_hash").
			NotEmpty(),
		field.Bool("is_active").
			Default(true),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.UUID("created_by", uuid.UUID{}),
		field.UUID("updated_by", uuid.UUID{}),
	}
}

func (AccountLogin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_account", UserAccount.Type).
			Unique(),
	}
}
