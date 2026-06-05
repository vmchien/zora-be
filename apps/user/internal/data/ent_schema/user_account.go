package ent_schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type UserAccount struct {
	ent.Schema
}

func (UserAccount) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_login_id", uuid.UUID{}).
			Unique(),
		field.String("full_name").
			Optional().
			Nillable(),
		field.String("avatar_url").
			Optional().
			Nillable(),
		field.String("phone").
			Optional().
			Nillable(),
		field.Time("date_of_birth").
			Optional().
			Nillable(),
		field.String("gender").
			Optional().
			Nillable(),
		field.String("address").
			Optional().
			Nillable(),
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

func (UserAccount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account_login", AccountLogin.Type).
			Ref("user_account").
			Field("user_login_id").
			Unique().
			Required(),
	}
}
