package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/gofrs/uuid"
)

func newUUID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

type BaseMixin struct {
	mixin.Schema
}

func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		// org_id represents the tenant id of the entity.
		// field.UUID("org_id", uuid.UUID{}).Default(newUUID),
		// id represents the identity of the entity.
		field.UUID("id", uuid.UUID{}).Default(newUUID),
	}
}

// TimeMixin adds created_at and updated_at fields.
type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
