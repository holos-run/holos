package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/gofrs/uuid"
	storage "github.com/holos-run/holos/service/gen/holos/storage/v1alpha1"
)

type Platform struct {
	ent.Schema
}

func (Platform) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimestampMixin{},
		EditorMixin{},
	}
}

func (Platform) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("org_id", uuid.UUID{}),
		field.String("name").NotEmpty(),
		field.String("display_name"),
		field.JSON("form", &storage.Form{}).
			Optional().
			Comment("JSON representation of FormlyFormConfig[] refer to https://github.com/holos-run/holos/issues/161"),
		field.JSON("model", &storage.Model{}).
			Optional().
			Comment("JSON representation of the form model which holds user input values refer to https://github.com/holos-run/holos/issues/161"),
		field.Bytes("cue").
			Optional().
			Comment("CUE definition to vet the model against e.g. #PlatformConfig"),
		field.String("cue_definition").
			Optional().
			Comment("The definition name to vet config_values against config_cue e.g. '#PlatformSpec'"),
	}
}

func (Platform) Edges() []ent.Edge {
	return []ent.Edge{
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
