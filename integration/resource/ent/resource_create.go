// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/woocoos/kocli/integration/resource/ent/resource"
)

// ResourceCreate is the builder for creating a Resource entity.
type ResourceCreate struct {
	config
	mutation *ResourceMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (rc *ResourceCreate) SetName(s string) *ResourceCreate {
	rc.mutation.SetName(s)
	return rc
}

// SetDescription sets the "description" field.
func (rc *ResourceCreate) SetDescription(s string) *ResourceCreate {
	rc.mutation.SetDescription(s)
	return rc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (rc *ResourceCreate) SetNillableDescription(s *string) *ResourceCreate {
	if s != nil {
		rc.SetDescription(*s)
	}
	return rc
}

// SetTenantID sets the "tenant_id" field.
func (rc *ResourceCreate) SetTenantID(s string) *ResourceCreate {
	rc.mutation.SetTenantID(s)
	return rc
}

// SetNillableTenantID sets the "tenant_id" field if the given value is not nil.
func (rc *ResourceCreate) SetNillableTenantID(s *string) *ResourceCreate {
	if s != nil {
		rc.SetTenantID(*s)
	}
	return rc
}

// Mutation returns the ResourceMutation object of the builder.
func (rc *ResourceCreate) Mutation() *ResourceMutation {
	return rc.mutation
}

// Save creates the Resource in the database.
func (rc *ResourceCreate) Save(ctx context.Context) (*Resource, error) {
	return withHooks(ctx, rc.sqlSave, rc.mutation, rc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (rc *ResourceCreate) SaveX(ctx context.Context) *Resource {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rc *ResourceCreate) Exec(ctx context.Context) error {
	_, err := rc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rc *ResourceCreate) ExecX(ctx context.Context) {
	if err := rc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *ResourceCreate) check() error {
	if _, ok := rc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Resource.name"`)}
	}
	if v, ok := rc.mutation.Name(); ok {
		if err := resource.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Resource.name": %w`, err)}
		}
	}
	return nil
}

func (rc *ResourceCreate) sqlSave(ctx context.Context) (*Resource, error) {
	if err := rc.check(); err != nil {
		return nil, err
	}
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	rc.mutation.id = &_node.ID
	rc.mutation.done = true
	return _node, nil
}

func (rc *ResourceCreate) createSpec() (*Resource, *sqlgraph.CreateSpec) {
	var (
		_node = &Resource{config: rc.config}
		_spec = sqlgraph.NewCreateSpec(resource.Table, sqlgraph.NewFieldSpec(resource.FieldID, field.TypeInt))
	)
	if value, ok := rc.mutation.Name(); ok {
		_spec.SetField(resource.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := rc.mutation.Description(); ok {
		_spec.SetField(resource.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := rc.mutation.TenantID(); ok {
		_spec.SetField(resource.FieldTenantID, field.TypeString, value)
		_node.TenantID = value
	}
	return _node, _spec
}

// ResourceCreateBulk is the builder for creating many Resource entities in bulk.
type ResourceCreateBulk struct {
	config
	err      error
	builders []*ResourceCreate
}

// Save creates the Resource entities in the database.
func (rcb *ResourceCreateBulk) Save(ctx context.Context) ([]*Resource, error) {
	if rcb.err != nil {
		return nil, rcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Resource, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ResourceMutation)
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
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
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
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *ResourceCreateBulk) SaveX(ctx context.Context) []*Resource {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rcb *ResourceCreateBulk) Exec(ctx context.Context) error {
	_, err := rcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcb *ResourceCreateBulk) ExecX(ctx context.Context) {
	if err := rcb.Exec(ctx); err != nil {
		panic(err)
	}
}
