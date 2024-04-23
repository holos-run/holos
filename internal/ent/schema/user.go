package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity, the internal
// representation and identity of a single human user. Users are scoped
// globally.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").NotEmpty().Unique(),
		field.String("iss"),
		field.String("sub"),
		field.String("name"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organizations", Organization.Type).
			Ref("users"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("iss", "sub").Unique(),
	}
}
