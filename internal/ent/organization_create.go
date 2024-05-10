// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/gofrs/uuid"
	"github.com/holos-run/holos/internal/ent/organization"
	"github.com/holos-run/holos/internal/ent/platform"
	"github.com/holos-run/holos/internal/ent/user"
)

// OrganizationCreate is the builder for creating a Organization entity.
type OrganizationCreate struct {
	config
	mutation *OrganizationMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (oc *OrganizationCreate) SetCreatedAt(t time.Time) *OrganizationCreate {
	oc.mutation.SetCreatedAt(t)
	return oc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (oc *OrganizationCreate) SetNillableCreatedAt(t *time.Time) *OrganizationCreate {
	if t != nil {
		oc.SetCreatedAt(*t)
	}
	return oc
}

// SetUpdatedAt sets the "updated_at" field.
func (oc *OrganizationCreate) SetUpdatedAt(t time.Time) *OrganizationCreate {
	oc.mutation.SetUpdatedAt(t)
	return oc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (oc *OrganizationCreate) SetNillableUpdatedAt(t *time.Time) *OrganizationCreate {
	if t != nil {
		oc.SetUpdatedAt(*t)
	}
	return oc
}

// SetCreatedByID sets the "created_by_id" field.
func (oc *OrganizationCreate) SetCreatedByID(u uuid.UUID) *OrganizationCreate {
	oc.mutation.SetCreatedByID(u)
	return oc
}

// SetUpdatedByID sets the "updated_by_id" field.
func (oc *OrganizationCreate) SetUpdatedByID(u uuid.UUID) *OrganizationCreate {
	oc.mutation.SetUpdatedByID(u)
	return oc
}

// SetName sets the "name" field.
func (oc *OrganizationCreate) SetName(s string) *OrganizationCreate {
	oc.mutation.SetName(s)
	return oc
}

// SetDisplayName sets the "display_name" field.
func (oc *OrganizationCreate) SetDisplayName(s string) *OrganizationCreate {
	oc.mutation.SetDisplayName(s)
	return oc
}

// SetID sets the "id" field.
func (oc *OrganizationCreate) SetID(u uuid.UUID) *OrganizationCreate {
	oc.mutation.SetID(u)
	return oc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (oc *OrganizationCreate) SetNillableID(u *uuid.UUID) *OrganizationCreate {
	if u != nil {
		oc.SetID(*u)
	}
	return oc
}

// SetCreatorID sets the "creator" edge to the User entity by ID.
func (oc *OrganizationCreate) SetCreatorID(id uuid.UUID) *OrganizationCreate {
	oc.mutation.SetCreatorID(id)
	return oc
}

// SetCreator sets the "creator" edge to the User entity.
func (oc *OrganizationCreate) SetCreator(u *User) *OrganizationCreate {
	return oc.SetCreatorID(u.ID)
}

// SetEditorID sets the "editor" edge to the User entity by ID.
func (oc *OrganizationCreate) SetEditorID(id uuid.UUID) *OrganizationCreate {
	oc.mutation.SetEditorID(id)
	return oc
}

// SetEditor sets the "editor" edge to the User entity.
func (oc *OrganizationCreate) SetEditor(u *User) *OrganizationCreate {
	return oc.SetEditorID(u.ID)
}

// AddUserIDs adds the "users" edge to the User entity by IDs.
func (oc *OrganizationCreate) AddUserIDs(ids ...uuid.UUID) *OrganizationCreate {
	oc.mutation.AddUserIDs(ids...)
	return oc
}

// AddUsers adds the "users" edges to the User entity.
func (oc *OrganizationCreate) AddUsers(u ...*User) *OrganizationCreate {
	ids := make([]uuid.UUID, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return oc.AddUserIDs(ids...)
}

// AddPlatformIDs adds the "platforms" edge to the Platform entity by IDs.
func (oc *OrganizationCreate) AddPlatformIDs(ids ...uuid.UUID) *OrganizationCreate {
	oc.mutation.AddPlatformIDs(ids...)
	return oc
}

// AddPlatforms adds the "platforms" edges to the Platform entity.
func (oc *OrganizationCreate) AddPlatforms(p ...*Platform) *OrganizationCreate {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return oc.AddPlatformIDs(ids...)
}

// Mutation returns the OrganizationMutation object of the builder.
func (oc *OrganizationCreate) Mutation() *OrganizationMutation {
	return oc.mutation
}

// Save creates the Organization in the database.
func (oc *OrganizationCreate) Save(ctx context.Context) (*Organization, error) {
	oc.defaults()
	return withHooks(ctx, oc.sqlSave, oc.mutation, oc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (oc *OrganizationCreate) SaveX(ctx context.Context) *Organization {
	v, err := oc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (oc *OrganizationCreate) Exec(ctx context.Context) error {
	_, err := oc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (oc *OrganizationCreate) ExecX(ctx context.Context) {
	if err := oc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (oc *OrganizationCreate) defaults() {
	if _, ok := oc.mutation.CreatedAt(); !ok {
		v := organization.DefaultCreatedAt()
		oc.mutation.SetCreatedAt(v)
	}
	if _, ok := oc.mutation.UpdatedAt(); !ok {
		v := organization.DefaultUpdatedAt()
		oc.mutation.SetUpdatedAt(v)
	}
	if _, ok := oc.mutation.ID(); !ok {
		v := organization.DefaultID()
		oc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (oc *OrganizationCreate) check() error {
	if _, ok := oc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Organization.created_at"`)}
	}
	if _, ok := oc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Organization.updated_at"`)}
	}
	if _, ok := oc.mutation.CreatedByID(); !ok {
		return &ValidationError{Name: "created_by_id", err: errors.New(`ent: missing required field "Organization.created_by_id"`)}
	}
	if _, ok := oc.mutation.UpdatedByID(); !ok {
		return &ValidationError{Name: "updated_by_id", err: errors.New(`ent: missing required field "Organization.updated_by_id"`)}
	}
	if _, ok := oc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Organization.name"`)}
	}
	if v, ok := oc.mutation.Name(); ok {
		if err := organization.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Organization.name": %w`, err)}
		}
	}
	if _, ok := oc.mutation.DisplayName(); !ok {
		return &ValidationError{Name: "display_name", err: errors.New(`ent: missing required field "Organization.display_name"`)}
	}
	if _, ok := oc.mutation.CreatorID(); !ok {
		return &ValidationError{Name: "creator", err: errors.New(`ent: missing required edge "Organization.creator"`)}
	}
	if _, ok := oc.mutation.EditorID(); !ok {
		return &ValidationError{Name: "editor", err: errors.New(`ent: missing required edge "Organization.editor"`)}
	}
	return nil
}

func (oc *OrganizationCreate) sqlSave(ctx context.Context) (*Organization, error) {
	if err := oc.check(); err != nil {
		return nil, err
	}
	_node, _spec := oc.createSpec()
	if err := sqlgraph.CreateNode(ctx, oc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	oc.mutation.id = &_node.ID
	oc.mutation.done = true
	return _node, nil
}

func (oc *OrganizationCreate) createSpec() (*Organization, *sqlgraph.CreateSpec) {
	var (
		_node = &Organization{config: oc.config}
		_spec = sqlgraph.NewCreateSpec(organization.Table, sqlgraph.NewFieldSpec(organization.FieldID, field.TypeUUID))
	)
	_spec.OnConflict = oc.conflict
	if id, ok := oc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := oc.mutation.CreatedAt(); ok {
		_spec.SetField(organization.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := oc.mutation.UpdatedAt(); ok {
		_spec.SetField(organization.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := oc.mutation.Name(); ok {
		_spec.SetField(organization.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := oc.mutation.DisplayName(); ok {
		_spec.SetField(organization.FieldDisplayName, field.TypeString, value)
		_node.DisplayName = value
	}
	if nodes := oc.mutation.CreatorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   organization.CreatorTable,
			Columns: []string{organization.CreatorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.CreatedByID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := oc.mutation.EditorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   organization.EditorTable,
			Columns: []string{organization.EditorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.UpdatedByID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := oc.mutation.UsersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   organization.UsersTable,
			Columns: organization.UsersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := oc.mutation.PlatformsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   organization.PlatformsTable,
			Columns: []string{organization.PlatformsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(platform.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Organization.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.OrganizationUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (oc *OrganizationCreate) OnConflict(opts ...sql.ConflictOption) *OrganizationUpsertOne {
	oc.conflict = opts
	return &OrganizationUpsertOne{
		create: oc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Organization.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (oc *OrganizationCreate) OnConflictColumns(columns ...string) *OrganizationUpsertOne {
	oc.conflict = append(oc.conflict, sql.ConflictColumns(columns...))
	return &OrganizationUpsertOne{
		create: oc,
	}
}

type (
	// OrganizationUpsertOne is the builder for "upsert"-ing
	//  one Organization node.
	OrganizationUpsertOne struct {
		create *OrganizationCreate
	}

	// OrganizationUpsert is the "OnConflict" setter.
	OrganizationUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *OrganizationUpsert) SetUpdatedAt(v time.Time) *OrganizationUpsert {
	u.Set(organization.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *OrganizationUpsert) UpdateUpdatedAt() *OrganizationUpsert {
	u.SetExcluded(organization.FieldUpdatedAt)
	return u
}

// SetUpdatedByID sets the "updated_by_id" field.
func (u *OrganizationUpsert) SetUpdatedByID(v uuid.UUID) *OrganizationUpsert {
	u.Set(organization.FieldUpdatedByID, v)
	return u
}

// UpdateUpdatedByID sets the "updated_by_id" field to the value that was provided on create.
func (u *OrganizationUpsert) UpdateUpdatedByID() *OrganizationUpsert {
	u.SetExcluded(organization.FieldUpdatedByID)
	return u
}

// SetName sets the "name" field.
func (u *OrganizationUpsert) SetName(v string) *OrganizationUpsert {
	u.Set(organization.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *OrganizationUpsert) UpdateName() *OrganizationUpsert {
	u.SetExcluded(organization.FieldName)
	return u
}

// SetDisplayName sets the "display_name" field.
func (u *OrganizationUpsert) SetDisplayName(v string) *OrganizationUpsert {
	u.Set(organization.FieldDisplayName, v)
	return u
}

// UpdateDisplayName sets the "display_name" field to the value that was provided on create.
func (u *OrganizationUpsert) UpdateDisplayName() *OrganizationUpsert {
	u.SetExcluded(organization.FieldDisplayName)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Organization.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(organization.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *OrganizationUpsertOne) UpdateNewValues() *OrganizationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(organization.FieldID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(organization.FieldCreatedAt)
		}
		if _, exists := u.create.mutation.CreatedByID(); exists {
			s.SetIgnore(organization.FieldCreatedByID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Organization.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *OrganizationUpsertOne) Ignore() *OrganizationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *OrganizationUpsertOne) DoNothing() *OrganizationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the OrganizationCreate.OnConflict
// documentation for more info.
func (u *OrganizationUpsertOne) Update(set func(*OrganizationUpsert)) *OrganizationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&OrganizationUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *OrganizationUpsertOne) SetUpdatedAt(v time.Time) *OrganizationUpsertOne {
	return u.Update(func(s *OrganizationUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *OrganizationUpsertOne) UpdateUpdatedAt() *OrganizationUpsertOne {
	return u.Update(func(s *OrganizationUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetUpdatedByID sets the "updated_by_id" field.
func (u *OrganizationUpsertOne) SetUpdatedByID(v uuid.UUID) *OrganizationUpsertOne {
	return u.Update(func(s *OrganizationUpsert) {
		s.SetUpdatedByID(v)
	})
}

// UpdateUpdatedByID sets the "updated_by_id" field to the value that was provided on create.
func (u *OrganizationUpsertOne) UpdateUpdatedByID() *OrganizationUpsertOne {
	return u.Update(func(s *OrganizationUpsert) {
		s.UpdateUpdatedByID()
	})
}

// SetName sets the "name" field.
func (u *OrganizationUpsertOne) SetName(v string) *OrganizationUpsertOne {
	return u.Update(func(s *OrganizationUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *OrganizationUpsertOne) UpdateName() *OrganizationUpsertOne {
	return u.Update(func(s *OrganizationUpsert) {
		s.UpdateName()
	})
}

// SetDisplayName sets the "display_name" field.
func (u *OrganizationUpsertOne) SetDisplayName(v string) *OrganizationUpsertOne {
	return u.Update(func(s *OrganizationUpsert) {
		s.SetDisplayName(v)
	})
}

// UpdateDisplayName sets the "display_name" field to the value that was provided on create.
func (u *OrganizationUpsertOne) UpdateDisplayName() *OrganizationUpsertOne {
	return u.Update(func(s *OrganizationUpsert) {
		s.UpdateDisplayName()
	})
}

// Exec executes the query.
func (u *OrganizationUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for OrganizationCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *OrganizationUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *OrganizationUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: OrganizationUpsertOne.ID is not supported by MySQL driver. Use OrganizationUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *OrganizationUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// OrganizationCreateBulk is the builder for creating many Organization entities in bulk.
type OrganizationCreateBulk struct {
	config
	err      error
	builders []*OrganizationCreate
	conflict []sql.ConflictOption
}

// Save creates the Organization entities in the database.
func (ocb *OrganizationCreateBulk) Save(ctx context.Context) ([]*Organization, error) {
	if ocb.err != nil {
		return nil, ocb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ocb.builders))
	nodes := make([]*Organization, len(ocb.builders))
	mutators := make([]Mutator, len(ocb.builders))
	for i := range ocb.builders {
		func(i int, root context.Context) {
			builder := ocb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*OrganizationMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ocb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ocb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ocb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ocb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ocb *OrganizationCreateBulk) SaveX(ctx context.Context) []*Organization {
	v, err := ocb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ocb *OrganizationCreateBulk) Exec(ctx context.Context) error {
	_, err := ocb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ocb *OrganizationCreateBulk) ExecX(ctx context.Context) {
	if err := ocb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Organization.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.OrganizationUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (ocb *OrganizationCreateBulk) OnConflict(opts ...sql.ConflictOption) *OrganizationUpsertBulk {
	ocb.conflict = opts
	return &OrganizationUpsertBulk{
		create: ocb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Organization.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ocb *OrganizationCreateBulk) OnConflictColumns(columns ...string) *OrganizationUpsertBulk {
	ocb.conflict = append(ocb.conflict, sql.ConflictColumns(columns...))
	return &OrganizationUpsertBulk{
		create: ocb,
	}
}

// OrganizationUpsertBulk is the builder for "upsert"-ing
// a bulk of Organization nodes.
type OrganizationUpsertBulk struct {
	create *OrganizationCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Organization.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(organization.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *OrganizationUpsertBulk) UpdateNewValues() *OrganizationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(organization.FieldID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(organization.FieldCreatedAt)
			}
			if _, exists := b.mutation.CreatedByID(); exists {
				s.SetIgnore(organization.FieldCreatedByID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Organization.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *OrganizationUpsertBulk) Ignore() *OrganizationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *OrganizationUpsertBulk) DoNothing() *OrganizationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the OrganizationCreateBulk.OnConflict
// documentation for more info.
func (u *OrganizationUpsertBulk) Update(set func(*OrganizationUpsert)) *OrganizationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&OrganizationUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *OrganizationUpsertBulk) SetUpdatedAt(v time.Time) *OrganizationUpsertBulk {
	return u.Update(func(s *OrganizationUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *OrganizationUpsertBulk) UpdateUpdatedAt() *OrganizationUpsertBulk {
	return u.Update(func(s *OrganizationUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetUpdatedByID sets the "updated_by_id" field.
func (u *OrganizationUpsertBulk) SetUpdatedByID(v uuid.UUID) *OrganizationUpsertBulk {
	return u.Update(func(s *OrganizationUpsert) {
		s.SetUpdatedByID(v)
	})
}

// UpdateUpdatedByID sets the "updated_by_id" field to the value that was provided on create.
func (u *OrganizationUpsertBulk) UpdateUpdatedByID() *OrganizationUpsertBulk {
	return u.Update(func(s *OrganizationUpsert) {
		s.UpdateUpdatedByID()
	})
}

// SetName sets the "name" field.
func (u *OrganizationUpsertBulk) SetName(v string) *OrganizationUpsertBulk {
	return u.Update(func(s *OrganizationUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *OrganizationUpsertBulk) UpdateName() *OrganizationUpsertBulk {
	return u.Update(func(s *OrganizationUpsert) {
		s.UpdateName()
	})
}

// SetDisplayName sets the "display_name" field.
func (u *OrganizationUpsertBulk) SetDisplayName(v string) *OrganizationUpsertBulk {
	return u.Update(func(s *OrganizationUpsert) {
		s.SetDisplayName(v)
	})
}

// UpdateDisplayName sets the "display_name" field to the value that was provided on create.
func (u *OrganizationUpsertBulk) UpdateDisplayName() *OrganizationUpsertBulk {
	return u.Update(func(s *OrganizationUpsert) {
		s.UpdateDisplayName()
	})
}

// Exec executes the query.
func (u *OrganizationUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the OrganizationCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for OrganizationCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *OrganizationUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
