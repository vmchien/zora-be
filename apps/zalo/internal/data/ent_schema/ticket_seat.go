package ent_schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"vn.vato.zora.be.api/pkg/data/ent_mixin"
	"vn.vato.zora.be.api/pkg/data/enums"
	"vn.vato.zora.be.api/pkg/def"
)

type TicketSeat struct {
	ent.Schema
}

func (TicketSeat) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ent_mixin.MixinID{},
		ent_mixin.MixinAuditTime{},
	}
}

func (TicketSeat) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("ticket_id", uuid.UUID{}),
		field.String("route_id").
			NotEmpty().
			MaxLen(50),
		field.Enum("seat_type").
			Values(enums.TicketSeatTypeValues()...).
			Default(enums.TicketSeatTypeOutbound.String()),
		field.String("route_name").
			NotEmpty(),
		field.Time("departure_time").
			Comment("Scheduled departure time for the ticket"),
		field.String("trip_id").
			Optional().
			Nillable().
			MaxLen(50),
		field.String("seat_id").
			NotEmpty().
			MaxLen(50),
		field.String("seat_name").
			NotEmpty(),
		field.Enum("status").
			Values(enums.TicketSeatStatusValues()...).
			Default(enums.TicketSeatStatusInitial.String()),
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

		field.Time("paid_at").
			Optional().
			Nillable().
			Comment("Timestamp when the ticket was paid"),
		field.Time("exported_at").
			Optional().
			Nillable().
			Comment("Timestamp when the ticket was exported"),
		field.Time("cancelled_at").
			Optional().
			Nillable().
			Comment("Timestamp when the ticket was cancelled"),
		field.Time("refunded_at").
			Optional().
			Nillable().
			Comment("Timestamp when the ticket was refunded"),

		field.String("kind").
			Optional().
			Nillable().
			MaxLen(50),

		field.String("origin_code").
			Optional().
			Nillable().
			MaxLen(50),
		field.String("origin_name").
			Optional().
			Nillable(),
		field.String("dest_code").
			Optional().
			Nillable().
			MaxLen(50),
		field.String("dest_name").
			Optional().
			Nillable(),
		field.Float("distance_km").
			SchemaType(map[string]string{
				dialect.Postgres: "numeric(18,4)",
			}).
			Optional().
			Nillable(),
		field.Int64("duration_ms").
			Optional().
			Nillable(),
		field.String("way_id").
			Optional().
			Nillable().
			MaxLen(50),
		field.Text("way_name").
			Optional().
			Nillable(),
		field.Enum("pickup_type").
			Values(enums.PickupDropOffValues()...).
			Nillable().
			Optional(),
		field.String("pickup_id").
			Optional().
			Nillable().
			MaxLen(50),
		field.Text("pickup_name").
			Optional().
			Nillable(),
		field.Text("pickup_address").
			Optional().
			Nillable(),
		field.JSON("pickup_info", def.JsonMap).
			Optional().
			Default(def.DefaultJsonMap()),
		field.Enum("drop_off_type").
			Values(enums.PickupDropOffValues()...).
			Nillable().
			Optional(),
		field.String("drop_off_id").
			Optional().
			Nillable().
			MaxLen(50),
		field.Text("drop_off_name").
			Optional().
			Nillable(),
		field.Text("drop_off_address").
			Optional().
			Nillable(),
		field.JSON("drop_off_info", def.JsonMap).
			Optional().
			Default(def.DefaultJsonMap()),

		field.JSON("extra_data", def.JsonMap).
			Optional().
			Default(def.DefaultJsonMap()),
	}
}

func (TicketSeat) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ticket_id"),
	}
}
func (TicketSeat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ticket", Ticket.Type).
			Ref("ticketSeats").
			Field("ticket_id").
			Required().
			Unique(),
	}
}
