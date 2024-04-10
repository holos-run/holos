package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
		field.Bool("email_verified").Default(false),
		field.String("name"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("identities", UserIdentity.Type).
			StorageKey(edge.Column("user_id")).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// UserIdentity holds the schema definition for the UserIdentity entity, a
// representation of an identity from an oidc provider typically implemented as
// an oidc id token.
type UserIdentity struct {
	ent.Schema
}

func (UserIdentity) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		TimeMixin{},
	}
}

// Fields of the UserIdentity. Assumes an oidc id token is available for the user with
// the openid scope `iss`, `aud`, `exp`, `iat`, and `at_hash` claims. Assumes
// the profile scope as well, with claims family_name, given_name, middle_name,
// nickname, picture, and updated_at. Claims are defined by the [standard
// claims][1].
//
// [1]: https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
func (UserIdentity) Fields() []ent.Field {
	return []ent.Field{
		field.String("iss").NotEmpty(),
		field.String("sub").NotEmpty(),
		field.String("email").NotEmpty(),
		field.Bool("email_verified").Default(false),
		field.String("name").Optional().Nillable(),
	}
}

func (UserIdentity) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("iss", "sub").Unique(),
	}
}

// Edges of the UserIdentity.
func (UserIdentity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("identities").
			Required().
			Immutable().
			Unique(),
	}
}
