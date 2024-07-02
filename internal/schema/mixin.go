package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/gofrs/uuid"
)

func newUUID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

// IDMixin mixes in an id field with a server generated default.  All resource
// objects in the storage system should be identified by a uuid.
type IDMixin struct {
	mixin.Schema
}

func (IDMixin) Fields() []ent.Field {
	return []ent.Field{
		// id represents the identity of the entity.
		field.UUID("id", uuid.UUID{}).Default(newUUID),
	}
}

// TimestampMixin adds created_at and updated_at fields.
type TimestampMixin struct {
	mixin.Schema
}

func (TimestampMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// EditorMixin adds created_by_id and updated_by_id fields representing the user
// who created or last modified the resource.
type EditorMixin struct {
	mixin.Schema
}

func (EditorMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("created_by_id", uuid.UUID{}).Immutable(),
		field.UUID("updated_by_id", uuid.UUID{}),
	}
}

func (EditorMixin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("creator", User.Type).
			Field("created_by_id").
			Immutable().
			Unique().
			Required(),
		edge.To("editor", User.Type).
			Field("updated_by_id").
			Unique().
			Required(),
	}
}
