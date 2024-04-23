package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/gofrs/uuid"
)

// Organization represents an organization account.
type Organization struct {
	ent.Schema
}

func (Organization) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		TimeMixin{},
	}
}

func (Organization) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
		field.String("display_name"),
		field.UUID("creator_id", uuid.UUID{}),
	}
}

func (Organization) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("creator", User.Type).
			Field("creator_id").
			Unique().
			Required(),
		edge.To("users", User.Type),
	}
}
