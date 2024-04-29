package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/gofrs/uuid"
)

type Platform struct {
	ent.Schema
}

func (Platform) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		TimeMixin{},
	}
}

func (Platform) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("org_id", uuid.UUID{}),
		field.String("name").NotEmpty(),
		field.String("display_name"),
		field.UUID("creator_id", uuid.UUID{}),
	}
}

func (Platform) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("creator", User.Type).
			Field("creator_id").
			Unique().
			Required(),
		edge.To("organization", Organization.Type).
			Field("org_id").
			Unique().
			Required(),
	}
}
