package ent_mixin

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"vn.vato.zora.be.api/pkg/constant"
	"vn.vato.zora.be.api/pkg/data/ent_annotation"
)

// MixinSoftDelete provides soft delete functionality for entities.
type MixinSoftDelete struct {
	mixin.Schema
}

func (MixinSoftDelete) Fields() []ent.Field {
	return []ent.Field{
		field.Bool(constant.IsDeletedField).
			Default(false).
			Annotations(ent_annotation.SoftDeleteAnnotation{}).
			Comment("Indicates if the record is soft-deleted"),

		field.Time(constant.DeletedAtField).
			Optional().
			Nillable().
			Annotations(ent_annotation.SoftDeleteAnnotation{}).
			Comment("Timestamp when the record was soft-deleted"),
	}
}
func (MixinSoftDelete) Annotations() []schema.Annotation {
	return []schema.Annotation{
		ent_annotation.SoftDeleteAnnotation{},
	}
}

type softDeleteKey struct{}

// SkipSoftDelete returns a new context that skips the soft-delete interceptor/mutators.
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

// P adds a storage-level predicate to the queries and mutations.
func (d MixinSoftDelete) P(w interface{ WhereP(...func(*sql.Selector)) }) {
	w.WhereP(
		sql.FieldEQ(constant.IsDeletedField, false),
	)
}

// Interceptors of the SoftDeleteMixin.
func (d MixinSoftDelete) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			// Skip soft-delete, means include soft-deleted entities.
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}
			if qp, ok := q.(interface{ WhereP(...func(*sql.Selector)) }); ok {
				d.P(qp)
			}
			return nil
		}),
	}
}
