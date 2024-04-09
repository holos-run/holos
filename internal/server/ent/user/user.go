// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/gofrs/uuid"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldEmailVerified holds the string denoting the email_verified field in the database.
	FieldEmailVerified = "email_verified"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgeIdentities holds the string denoting the identities edge name in mutations.
	EdgeIdentities = "identities"
	// Table holds the table name of the user in the database.
	Table = "users"
	// IdentitiesTable is the table that holds the identities relation/edge.
	IdentitiesTable = "user_identities"
	// IdentitiesInverseTable is the table name for the UserIdentity entity.
	// It exists in this package in order to avoid circular dependency with the "useridentity" package.
	IdentitiesInverseTable = "user_identities"
	// IdentitiesColumn is the table column denoting the identities relation/edge.
	IdentitiesColumn = "user_id"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldEmail,
	FieldEmailVerified,
	FieldName,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator func(string) error
	// DefaultEmailVerified holds the default value on creation for the "email_verified" field.
	DefaultEmailVerified bool
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByEmail orders the results by the email field.
func ByEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmail, opts...).ToFunc()
}

// ByEmailVerified orders the results by the email_verified field.
func ByEmailVerified(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmailVerified, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByIdentitiesCount orders the results by identities count.
func ByIdentitiesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newIdentitiesStep(), opts...)
	}
}

// ByIdentities orders the results by identities terms.
func ByIdentities(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newIdentitiesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newIdentitiesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(IdentitiesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, IdentitiesTable, IdentitiesColumn),
	)
}