// Code generated by ent, DO NOT EDIT.

package platform

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/gofrs/uuid"
)

const (
	// Label holds the string label denoting the platform type in the database.
	Label = "platform"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldCreatedByID holds the string denoting the created_by_id field in the database.
	FieldCreatedByID = "created_by_id"
	// FieldUpdatedByID holds the string denoting the updated_by_id field in the database.
	FieldUpdatedByID = "updated_by_id"
	// FieldOrgID holds the string denoting the org_id field in the database.
	FieldOrgID = "org_id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDisplayName holds the string denoting the display_name field in the database.
	FieldDisplayName = "display_name"
	// FieldForm holds the string denoting the form field in the database.
	FieldForm = "form"
	// FieldModel holds the string denoting the model field in the database.
	FieldModel = "model"
	// FieldCue holds the string denoting the cue field in the database.
	FieldCue = "cue"
	// FieldCueDefinition holds the string denoting the cue_definition field in the database.
	FieldCueDefinition = "cue_definition"
	// EdgeCreator holds the string denoting the creator edge name in mutations.
	EdgeCreator = "creator"
	// EdgeEditor holds the string denoting the editor edge name in mutations.
	EdgeEditor = "editor"
	// EdgeOrganization holds the string denoting the organization edge name in mutations.
	EdgeOrganization = "organization"
	// Table holds the table name of the platform in the database.
	Table = "platforms"
	// CreatorTable is the table that holds the creator relation/edge.
	CreatorTable = "platforms"
	// CreatorInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	CreatorInverseTable = "users"
	// CreatorColumn is the table column denoting the creator relation/edge.
	CreatorColumn = "created_by_id"
	// EditorTable is the table that holds the editor relation/edge.
	EditorTable = "platforms"
	// EditorInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	EditorInverseTable = "users"
	// EditorColumn is the table column denoting the editor relation/edge.
	EditorColumn = "updated_by_id"
	// OrganizationTable is the table that holds the organization relation/edge.
	OrganizationTable = "platforms"
	// OrganizationInverseTable is the table name for the Organization entity.
	// It exists in this package in order to avoid circular dependency with the "organization" package.
	OrganizationInverseTable = "organizations"
	// OrganizationColumn is the table column denoting the organization relation/edge.
	OrganizationColumn = "org_id"
)

// Columns holds all SQL columns for platform fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldCreatedByID,
	FieldUpdatedByID,
	FieldOrgID,
	FieldName,
	FieldDisplayName,
	FieldForm,
	FieldModel,
	FieldCue,
	FieldCueDefinition,
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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Platform queries.
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

// ByCreatedByID orders the results by the created_by_id field.
func ByCreatedByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedByID, opts...).ToFunc()
}

// ByUpdatedByID orders the results by the updated_by_id field.
func ByUpdatedByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedByID, opts...).ToFunc()
}

// ByOrgID orders the results by the org_id field.
func ByOrgID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOrgID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDisplayName orders the results by the display_name field.
func ByDisplayName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDisplayName, opts...).ToFunc()
}

// ByCueDefinition orders the results by the cue_definition field.
func ByCueDefinition(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCueDefinition, opts...).ToFunc()
}

// ByCreatorField orders the results by creator field.
func ByCreatorField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCreatorStep(), sql.OrderByField(field, opts...))
	}
}

// ByEditorField orders the results by editor field.
func ByEditorField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEditorStep(), sql.OrderByField(field, opts...))
	}
}

// ByOrganizationField orders the results by organization field.
func ByOrganizationField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOrganizationStep(), sql.OrderByField(field, opts...))
	}
}
func newCreatorStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CreatorInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, CreatorTable, CreatorColumn),
	)
}
func newEditorStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EditorInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, EditorTable, EditorColumn),
	)
}
func newOrganizationStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OrganizationInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, OrganizationTable, OrganizationColumn),
	)
}