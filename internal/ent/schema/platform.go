package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.Bytes("config_form").
			Optional().
			Comment("Opaque JSON bytes representing the platform config form."),
		field.Bytes("config_values").
			Optional().
			Comment("Opaque JSON bytes representing the platform config values."),
		field.Bytes("config_cue").
			Optional().
			Comment("Opaque bytes representing the CUE definition of the config struct."),
		field.String("config_definition").
			Optional().
			Comment("The definition name to vet config_values against config_cue e.g. '#PlatformSpec'"),
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

func (Platform) Indexes() []ent.Index {
	return []ent.Index{
		// One org cannot have two platforms with the same name.
		index.Fields("org_id", "name").Unique(),
	}
}
