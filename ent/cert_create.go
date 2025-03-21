// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sifu-tool/ent/cert"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CertCreate is the builder for creating a Cert entity.
type CertCreate struct {
	config
	mutation *CertMutation
	hooks    []Hook
}

// SetTag sets the "tag" field.
func (cc *CertCreate) SetTag(s string) *CertCreate {
	cc.mutation.SetTag(s)
	return cc
}

// Mutation returns the CertMutation object of the builder.
func (cc *CertCreate) Mutation() *CertMutation {
	return cc.mutation
}

// Save creates the Cert in the database.
func (cc *CertCreate) Save(ctx context.Context) (*Cert, error) {
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *CertCreate) SaveX(ctx context.Context) *Cert {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *CertCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *CertCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *CertCreate) check() error {
	if _, ok := cc.mutation.Tag(); !ok {
		return &ValidationError{Name: "tag", err: errors.New(`ent: missing required field "Cert.tag"`)}
	}
	if v, ok := cc.mutation.Tag(); ok {
		if err := cert.TagValidator(v); err != nil {
			return &ValidationError{Name: "tag", err: fmt.Errorf(`ent: validator failed for field "Cert.tag": %w`, err)}
		}
	}
	return nil
}

func (cc *CertCreate) sqlSave(ctx context.Context) (*Cert, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *CertCreate) createSpec() (*Cert, *sqlgraph.CreateSpec) {
	var (
		_node = &Cert{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(cert.Table, sqlgraph.NewFieldSpec(cert.FieldID, field.TypeInt))
	)
	if value, ok := cc.mutation.Tag(); ok {
		_spec.SetField(cert.FieldTag, field.TypeString, value)
		_node.Tag = value
	}
	return _node, _spec
}

// CertCreateBulk is the builder for creating many Cert entities in bulk.
type CertCreateBulk struct {
	config
	err      error
	builders []*CertCreate
}

// Save creates the Cert entities in the database.
func (ccb *CertCreateBulk) Save(ctx context.Context) ([]*Cert, error) {
	if ccb.err != nil {
		return nil, ccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Cert, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CertMutation)
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
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *CertCreateBulk) SaveX(ctx context.Context) []*Cert {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *CertCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *CertCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
