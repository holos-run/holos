// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/holos-run/holos/internal/ent/organization"
	"github.com/holos-run/holos/internal/ent/predicate"
)

// OrganizationUpdate is the builder for updating Organization entities.
type OrganizationUpdate struct {
	config
	hooks    []Hook
	mutation *OrganizationMutation
}

// Where appends a list predicates to the OrganizationUpdate builder.
func (ou *OrganizationUpdate) Where(ps ...predicate.Organization) *OrganizationUpdate {
	ou.mutation.Where(ps...)
	return ou
}

// SetUpdatedAt sets the "updated_at" field.
func (ou *OrganizationUpdate) SetUpdatedAt(t time.Time) *OrganizationUpdate {
	ou.mutation.SetUpdatedAt(t)
	return ou
}

// SetName sets the "name" field.
func (ou *OrganizationUpdate) SetName(s string) *OrganizationUpdate {
	ou.mutation.SetName(s)
	return ou
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ou *OrganizationUpdate) SetNillableName(s *string) *OrganizationUpdate {
	if s != nil {
		ou.SetName(*s)
	}
	return ou
}

// SetDisplayName sets the "display_name" field.
func (ou *OrganizationUpdate) SetDisplayName(s string) *OrganizationUpdate {
	ou.mutation.SetDisplayName(s)
	return ou
}

// SetNillableDisplayName sets the "display_name" field if the given value is not nil.
func (ou *OrganizationUpdate) SetNillableDisplayName(s *string) *OrganizationUpdate {
	if s != nil {
		ou.SetDisplayName(*s)
	}
	return ou
}

// Mutation returns the OrganizationMutation object of the builder.
func (ou *OrganizationUpdate) Mutation() *OrganizationMutation {
	return ou.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ou *OrganizationUpdate) Save(ctx context.Context) (int, error) {
	ou.defaults()
	return withHooks(ctx, ou.sqlSave, ou.mutation, ou.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ou *OrganizationUpdate) SaveX(ctx context.Context) int {
	affected, err := ou.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ou *OrganizationUpdate) Exec(ctx context.Context) error {
	_, err := ou.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ou *OrganizationUpdate) ExecX(ctx context.Context) {
	if err := ou.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ou *OrganizationUpdate) defaults() {
	if _, ok := ou.mutation.UpdatedAt(); !ok {
		v := organization.UpdateDefaultUpdatedAt()
		ou.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ou *OrganizationUpdate) check() error {
	if v, ok := ou.mutation.Name(); ok {
		if err := organization.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Organization.name": %w`, err)}
		}
	}
	return nil
}

func (ou *OrganizationUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ou.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(organization.Table, organization.Columns, sqlgraph.NewFieldSpec(organization.FieldID, field.TypeUUID))
	if ps := ou.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ou.mutation.UpdatedAt(); ok {
		_spec.SetField(organization.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ou.mutation.Name(); ok {
		_spec.SetField(organization.FieldName, field.TypeString, value)
	}
	if value, ok := ou.mutation.DisplayName(); ok {
		_spec.SetField(organization.FieldDisplayName, field.TypeString, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ou.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{organization.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ou.mutation.done = true
	return n, nil
}

// OrganizationUpdateOne is the builder for updating a single Organization entity.
type OrganizationUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *OrganizationMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (ouo *OrganizationUpdateOne) SetUpdatedAt(t time.Time) *OrganizationUpdateOne {
	ouo.mutation.SetUpdatedAt(t)
	return ouo
}

// SetName sets the "name" field.
func (ouo *OrganizationUpdateOne) SetName(s string) *OrganizationUpdateOne {
	ouo.mutation.SetName(s)
	return ouo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ouo *OrganizationUpdateOne) SetNillableName(s *string) *OrganizationUpdateOne {
	if s != nil {
		ouo.SetName(*s)
	}
	return ouo
}

// SetDisplayName sets the "display_name" field.
func (ouo *OrganizationUpdateOne) SetDisplayName(s string) *OrganizationUpdateOne {
	ouo.mutation.SetDisplayName(s)
	return ouo
}

// SetNillableDisplayName sets the "display_name" field if the given value is not nil.
func (ouo *OrganizationUpdateOne) SetNillableDisplayName(s *string) *OrganizationUpdateOne {
	if s != nil {
		ouo.SetDisplayName(*s)
	}
	return ouo
}

// Mutation returns the OrganizationMutation object of the builder.
func (ouo *OrganizationUpdateOne) Mutation() *OrganizationMutation {
	return ouo.mutation
}

// Where appends a list predicates to the OrganizationUpdate builder.
func (ouo *OrganizationUpdateOne) Where(ps ...predicate.Organization) *OrganizationUpdateOne {
	ouo.mutation.Where(ps...)
	return ouo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ouo *OrganizationUpdateOne) Select(field string, fields ...string) *OrganizationUpdateOne {
	ouo.fields = append([]string{field}, fields...)
	return ouo
}

// Save executes the query and returns the updated Organization entity.
func (ouo *OrganizationUpdateOne) Save(ctx context.Context) (*Organization, error) {
	ouo.defaults()
	return withHooks(ctx, ouo.sqlSave, ouo.mutation, ouo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ouo *OrganizationUpdateOne) SaveX(ctx context.Context) *Organization {
	node, err := ouo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ouo *OrganizationUpdateOne) Exec(ctx context.Context) error {
	_, err := ouo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ouo *OrganizationUpdateOne) ExecX(ctx context.Context) {
	if err := ouo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ouo *OrganizationUpdateOne) defaults() {
	if _, ok := ouo.mutation.UpdatedAt(); !ok {
		v := organization.UpdateDefaultUpdatedAt()
		ouo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ouo *OrganizationUpdateOne) check() error {
	if v, ok := ouo.mutation.Name(); ok {
		if err := organization.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Organization.name": %w`, err)}
		}
	}
	return nil
}

func (ouo *OrganizationUpdateOne) sqlSave(ctx context.Context) (_node *Organization, err error) {
	if err := ouo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(organization.Table, organization.Columns, sqlgraph.NewFieldSpec(organization.FieldID, field.TypeUUID))
	id, ok := ouo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Organization.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ouo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, organization.FieldID)
		for _, f := range fields {
			if !organization.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != organization.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ouo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ouo.mutation.UpdatedAt(); ok {
		_spec.SetField(organization.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ouo.mutation.Name(); ok {
		_spec.SetField(organization.FieldName, field.TypeString, value)
	}
	if value, ok := ouo.mutation.DisplayName(); ok {
		_spec.SetField(organization.FieldDisplayName, field.TypeString, value)
	}
	_node = &Organization{config: ouo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ouo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{organization.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ouo.mutation.done = true
	return _node, nil
}