package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/gofrs/uuid"
	holos "github.com/holos-run/holos/service/gen/holos/v1alpha1"
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
		field.JSON("config_form", &holos.PlatformForm{}).
			Optional().
			Comment("JSON holos.PlatformForm representing the platform data entry form."),
		field.JSON("config_values", &holos.UserDefinedConfig{}).
			Optional().
			Comment("JSON holos.ConfigValues representing the platform config values."),
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
