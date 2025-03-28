// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sifu-tool/ent/ddns"
	"sifu-tool/models"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// DDNSCreate is the builder for creating a DDNS entity.
type DDNSCreate struct {
	config
	mutation *DDNSMutation
	hooks    []Hook
}

// SetV4method sets the "v4method" field.
func (dc *DDNSCreate) SetV4method(i int) *DDNSCreate {
	dc.mutation.SetV4method(i)
	return dc
}

// SetNillableV4method sets the "v4method" field if the given value is not nil.
func (dc *DDNSCreate) SetNillableV4method(i *int) *DDNSCreate {
	if i != nil {
		dc.SetV4method(*i)
	}
	return dc
}

// SetV6method sets the "v6method" field.
func (dc *DDNSCreate) SetV6method(i int) *DDNSCreate {
	dc.mutation.SetV6method(i)
	return dc
}

// SetNillableV6method sets the "v6method" field if the given value is not nil.
func (dc *DDNSCreate) SetNillableV6method(i *int) *DDNSCreate {
	if i != nil {
		dc.SetV6method(*i)
	}
	return dc
}

// SetIpv6 sets the "ipv6" field.
func (dc *DDNSCreate) SetIpv6(s string) *DDNSCreate {
	dc.mutation.SetIpv6(s)
	return dc
}

// SetNillableIpv6 sets the "ipv6" field if the given value is not nil.
func (dc *DDNSCreate) SetNillableIpv6(s *string) *DDNSCreate {
	if s != nil {
		dc.SetIpv6(*s)
	}
	return dc
}

// SetRev6 sets the "rev6" field.
func (dc *DDNSCreate) SetRev6(s string) *DDNSCreate {
	dc.mutation.SetRev6(s)
	return dc
}

// SetNillableRev6 sets the "rev6" field if the given value is not nil.
func (dc *DDNSCreate) SetNillableRev6(s *string) *DDNSCreate {
	if s != nil {
		dc.SetRev6(*s)
	}
	return dc
}

// SetIpv4 sets the "ipv4" field.
func (dc *DDNSCreate) SetIpv4(s string) *DDNSCreate {
	dc.mutation.SetIpv4(s)
	return dc
}

// SetNillableIpv4 sets the "ipv4" field if the given value is not nil.
func (dc *DDNSCreate) SetNillableIpv4(s *string) *DDNSCreate {
	if s != nil {
		dc.SetIpv4(*s)
	}
	return dc
}

// SetRev4 sets the "rev4" field.
func (dc *DDNSCreate) SetRev4(s string) *DDNSCreate {
	dc.mutation.SetRev4(s)
	return dc
}

// SetNillableRev4 sets the "rev4" field if the given value is not nil.
func (dc *DDNSCreate) SetNillableRev4(s *string) *DDNSCreate {
	if s != nil {
		dc.SetRev4(*s)
	}
	return dc
}

// SetV4script sets the "v4script" field.
func (dc *DDNSCreate) SetV4script(s string) *DDNSCreate {
	dc.mutation.SetV4script(s)
	return dc
}

// SetNillableV4script sets the "v4script" field if the given value is not nil.
func (dc *DDNSCreate) SetNillableV4script(s *string) *DDNSCreate {
	if s != nil {
		dc.SetV4script(*s)
	}
	return dc
}

// SetV4interface sets the "v4interface" field.
func (dc *DDNSCreate) SetV4interface(s string) *DDNSCreate {
	dc.mutation.SetV4interface(s)
	return dc
}

// SetNillableV4interface sets the "v4interface" field if the given value is not nil.
func (dc *DDNSCreate) SetNillableV4interface(s *string) *DDNSCreate {
	if s != nil {
		dc.SetV4interface(*s)
	}
	return dc
}

// SetV6script sets the "v6script" field.
func (dc *DDNSCreate) SetV6script(s string) *DDNSCreate {
	dc.mutation.SetV6script(s)
	return dc
}

// SetNillableV6script sets the "v6script" field if the given value is not nil.
func (dc *DDNSCreate) SetNillableV6script(s *string) *DDNSCreate {
	if s != nil {
		dc.SetV6script(*s)
	}
	return dc
}

// SetV6interface sets the "v6interface" field.
func (dc *DDNSCreate) SetV6interface(s string) *DDNSCreate {
	dc.mutation.SetV6interface(s)
	return dc
}

// SetNillableV6interface sets the "v6interface" field if the given value is not nil.
func (dc *DDNSCreate) SetNillableV6interface(s *string) *DDNSCreate {
	if s != nil {
		dc.SetV6interface(*s)
	}
	return dc
}

// SetDomains sets the "domains" field.
func (dc *DDNSCreate) SetDomains(m []models.Domain) *DDNSCreate {
	dc.mutation.SetDomains(m)
	return dc
}

// SetConfig sets the "config" field.
func (dc *DDNSCreate) SetConfig(m map[string]string) *DDNSCreate {
	dc.mutation.SetConfig(m)
	return dc
}

// SetWebhook sets the "webhook" field.
func (dc *DDNSCreate) SetWebhook(m map[string]string) *DDNSCreate {
	dc.mutation.SetWebhook(m)
	return dc
}

// Mutation returns the DDNSMutation object of the builder.
func (dc *DDNSCreate) Mutation() *DDNSMutation {
	return dc.mutation
}

// Save creates the DDNS in the database.
func (dc *DDNSCreate) Save(ctx context.Context) (*DDNS, error) {
	return withHooks(ctx, dc.sqlSave, dc.mutation, dc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (dc *DDNSCreate) SaveX(ctx context.Context) *DDNS {
	v, err := dc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (dc *DDNSCreate) Exec(ctx context.Context) error {
	_, err := dc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (dc *DDNSCreate) ExecX(ctx context.Context) {
	if err := dc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (dc *DDNSCreate) check() error {
	if v, ok := dc.mutation.Ipv6(); ok {
		if err := ddns.Ipv6Validator(v); err != nil {
			return &ValidationError{Name: "ipv6", err: fmt.Errorf(`ent: validator failed for field "DDNS.ipv6": %w`, err)}
		}
	}
	if v, ok := dc.mutation.Rev6(); ok {
		if err := ddns.Rev6Validator(v); err != nil {
			return &ValidationError{Name: "rev6", err: fmt.Errorf(`ent: validator failed for field "DDNS.rev6": %w`, err)}
		}
	}
	if v, ok := dc.mutation.Ipv4(); ok {
		if err := ddns.Ipv4Validator(v); err != nil {
			return &ValidationError{Name: "ipv4", err: fmt.Errorf(`ent: validator failed for field "DDNS.ipv4": %w`, err)}
		}
	}
	if v, ok := dc.mutation.Rev4(); ok {
		if err := ddns.Rev4Validator(v); err != nil {
			return &ValidationError{Name: "rev4", err: fmt.Errorf(`ent: validator failed for field "DDNS.rev4": %w`, err)}
		}
	}
	if v, ok := dc.mutation.V4script(); ok {
		if err := ddns.V4scriptValidator(v); err != nil {
			return &ValidationError{Name: "v4script", err: fmt.Errorf(`ent: validator failed for field "DDNS.v4script": %w`, err)}
		}
	}
	if v, ok := dc.mutation.V4interface(); ok {
		if err := ddns.V4interfaceValidator(v); err != nil {
			return &ValidationError{Name: "v4interface", err: fmt.Errorf(`ent: validator failed for field "DDNS.v4interface": %w`, err)}
		}
	}
	if v, ok := dc.mutation.V6script(); ok {
		if err := ddns.V6scriptValidator(v); err != nil {
			return &ValidationError{Name: "v6script", err: fmt.Errorf(`ent: validator failed for field "DDNS.v6script": %w`, err)}
		}
	}
	if v, ok := dc.mutation.V6interface(); ok {
		if err := ddns.V6interfaceValidator(v); err != nil {
			return &ValidationError{Name: "v6interface", err: fmt.Errorf(`ent: validator failed for field "DDNS.v6interface": %w`, err)}
		}
	}
	if _, ok := dc.mutation.Domains(); !ok {
		return &ValidationError{Name: "domains", err: errors.New(`ent: missing required field "DDNS.domains"`)}
	}
	if _, ok := dc.mutation.Config(); !ok {
		return &ValidationError{Name: "config", err: errors.New(`ent: missing required field "DDNS.config"`)}
	}
	return nil
}

func (dc *DDNSCreate) sqlSave(ctx context.Context) (*DDNS, error) {
	if err := dc.check(); err != nil {
		return nil, err
	}
	_node, _spec := dc.createSpec()
	if err := sqlgraph.CreateNode(ctx, dc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	dc.mutation.id = &_node.ID
	dc.mutation.done = true
	return _node, nil
}

func (dc *DDNSCreate) createSpec() (*DDNS, *sqlgraph.CreateSpec) {
	var (
		_node = &DDNS{config: dc.config}
		_spec = sqlgraph.NewCreateSpec(ddns.Table, sqlgraph.NewFieldSpec(ddns.FieldID, field.TypeInt))
	)
	if value, ok := dc.mutation.V4method(); ok {
		_spec.SetField(ddns.FieldV4method, field.TypeInt, value)
		_node.V4method = value
	}
	if value, ok := dc.mutation.V6method(); ok {
		_spec.SetField(ddns.FieldV6method, field.TypeInt, value)
		_node.V6method = value
	}
	if value, ok := dc.mutation.Ipv6(); ok {
		_spec.SetField(ddns.FieldIpv6, field.TypeString, value)
		_node.Ipv6 = value
	}
	if value, ok := dc.mutation.Rev6(); ok {
		_spec.SetField(ddns.FieldRev6, field.TypeString, value)
		_node.Rev6 = value
	}
	if value, ok := dc.mutation.Ipv4(); ok {
		_spec.SetField(ddns.FieldIpv4, field.TypeString, value)
		_node.Ipv4 = value
	}
	if value, ok := dc.mutation.Rev4(); ok {
		_spec.SetField(ddns.FieldRev4, field.TypeString, value)
		_node.Rev4 = value
	}
	if value, ok := dc.mutation.V4script(); ok {
		_spec.SetField(ddns.FieldV4script, field.TypeString, value)
		_node.V4script = value
	}
	if value, ok := dc.mutation.V4interface(); ok {
		_spec.SetField(ddns.FieldV4interface, field.TypeString, value)
		_node.V4interface = value
	}
	if value, ok := dc.mutation.V6script(); ok {
		_spec.SetField(ddns.FieldV6script, field.TypeString, value)
		_node.V6script = value
	}
	if value, ok := dc.mutation.V6interface(); ok {
		_spec.SetField(ddns.FieldV6interface, field.TypeString, value)
		_node.V6interface = value
	}
	if value, ok := dc.mutation.Domains(); ok {
		_spec.SetField(ddns.FieldDomains, field.TypeJSON, value)
		_node.Domains = value
	}
	if value, ok := dc.mutation.Config(); ok {
		_spec.SetField(ddns.FieldConfig, field.TypeJSON, value)
		_node.Config = value
	}
	if value, ok := dc.mutation.Webhook(); ok {
		_spec.SetField(ddns.FieldWebhook, field.TypeJSON, value)
		_node.Webhook = value
	}
	return _node, _spec
}

// DDNSCreateBulk is the builder for creating many DDNS entities in bulk.
type DDNSCreateBulk struct {
	config
	err      error
	builders []*DDNSCreate
}

// Save creates the DDNS entities in the database.
func (dcb *DDNSCreateBulk) Save(ctx context.Context) ([]*DDNS, error) {
	if dcb.err != nil {
		return nil, dcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(dcb.builders))
	nodes := make([]*DDNS, len(dcb.builders))
	mutators := make([]Mutator, len(dcb.builders))
	for i := range dcb.builders {
		func(i int, root context.Context) {
			builder := dcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*DDNSMutation)
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
					_, err = mutators[i+1].Mutate(root, dcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, dcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, dcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (dcb *DDNSCreateBulk) SaveX(ctx context.Context) []*DDNS {
	v, err := dcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (dcb *DDNSCreateBulk) Exec(ctx context.Context) error {
	_, err := dcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (dcb *DDNSCreateBulk) ExecX(ctx context.Context) {
	if err := dcb.Exec(ctx); err != nil {
		panic(err)
	}
}
