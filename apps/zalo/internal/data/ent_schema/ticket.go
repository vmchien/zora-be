package ent_schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"vn.vato.zora.be.api/pkg/data/ent_mixin"
	"vn.vato.zora.be.api/pkg/data/enums"
)

type Ticket struct {
	ent.Schema
}

func (Ticket) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ent_mixin.MixinID{},
		ent_mixin.MixinAuditTime{},
		ent_mixin.MixinTenant{},
	}
}

func (Ticket) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			MaxLen(20).
			NotEmpty(),
		field.Enum("status").
			Values(enums.TicketStatusValues()...).
			Default(enums.TicketStatusUnknown.String()),
		field.Enum("way_type").
			Values(enums.WayTypeValues()...).
			Default(enums.OneWay.String()),
		field.Enum("booking_channel").
			Values(enums.ChannelValues()...).
			Default(enums.ChannelUnknown.String()),
		field.Int("payment_method").
			Default(0).
			Comment("Payment Method"),
		field.Bool("is_active").
			Default(true),
		field.String("promotion_code").
			Optional().
			Nillable(),

		field.Strings("departure_times").
			Default([]string{}).
			Comment("Departure Times: [outbound, return]"),
		field.Float("origin_amount").
			SchemaType(map[string]string{
				dialect.Postgres: "numeric(18,2)",
			}).
			Default(0).
			Min(0).
			Comment("Origin amount for the ticket"),
		field.Float("discount_amount").
			SchemaType(map[string]string{
				dialect.Postgres: "numeric(18,2)",
			}).
			Default(0).
			Min(0).
			Comment("Discount amount applied to the ticket"),
		field.Float("refund_amount").
			SchemaType(map[string]string{
				dialect.Postgres: "numeric(18,2)",
			}).
			Default(0).
			Min(0).
			Comment("Refund amount for the ticket"),

		field.String("remarks").
			Optional().
			Nillable(),

		field.String("reference_user_id").
			NotEmpty(),
		field.String("reference_phone"),
		field.String("reference_migration_id").
			Optional().
			Nillable(),
		field.String("force_updated_user_id").
			Optional().
			Nillable(),
	}
}

func (Ticket) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "code").Unique(),
		index.Fields("code"),
		index.Fields("code", "status"),
		index.Fields("reference_user_id"),
	}
}
func (Ticket) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ticketSeats", TicketSeat.Type),
	}
}
