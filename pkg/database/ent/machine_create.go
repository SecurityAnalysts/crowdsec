// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/crowdsecurity/crowdsec/pkg/database/ent/alert"
	"github.com/crowdsecurity/crowdsec/pkg/database/ent/machine"
)

// MachineCreate is the builder for creating a Machine entity.
type MachineCreate struct {
	config
	mutation *MachineMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (mc *MachineCreate) SetCreatedAt(t time.Time) *MachineCreate {
	mc.mutation.SetCreatedAt(t)
	return mc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (mc *MachineCreate) SetNillableCreatedAt(t *time.Time) *MachineCreate {
	if t != nil {
		mc.SetCreatedAt(*t)
	}
	return mc
}

// SetUpdatedAt sets the "updated_at" field.
func (mc *MachineCreate) SetUpdatedAt(t time.Time) *MachineCreate {
	mc.mutation.SetUpdatedAt(t)
	return mc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (mc *MachineCreate) SetNillableUpdatedAt(t *time.Time) *MachineCreate {
	if t != nil {
		mc.SetUpdatedAt(*t)
	}
	return mc
}

// SetLastPush sets the "last_push" field.
func (mc *MachineCreate) SetLastPush(t time.Time) *MachineCreate {
	mc.mutation.SetLastPush(t)
	return mc
}

// SetNillableLastPush sets the "last_push" field if the given value is not nil.
func (mc *MachineCreate) SetNillableLastPush(t *time.Time) *MachineCreate {
	if t != nil {
		mc.SetLastPush(*t)
	}
	return mc
}

// SetLastHeartbeat sets the "last_heartbeat" field.
func (mc *MachineCreate) SetLastHeartbeat(t time.Time) *MachineCreate {
	mc.mutation.SetLastHeartbeat(t)
	return mc
}

// SetNillableLastHeartbeat sets the "last_heartbeat" field if the given value is not nil.
func (mc *MachineCreate) SetNillableLastHeartbeat(t *time.Time) *MachineCreate {
	if t != nil {
		mc.SetLastHeartbeat(*t)
	}
	return mc
}

// SetMachineId sets the "machineId" field.
func (mc *MachineCreate) SetMachineId(s string) *MachineCreate {
	mc.mutation.SetMachineId(s)
	return mc
}

// SetPassword sets the "password" field.
func (mc *MachineCreate) SetPassword(s string) *MachineCreate {
	mc.mutation.SetPassword(s)
	return mc
}

// SetIpAddress sets the "ipAddress" field.
func (mc *MachineCreate) SetIpAddress(s string) *MachineCreate {
	mc.mutation.SetIpAddress(s)
	return mc
}

// SetScenarios sets the "scenarios" field.
func (mc *MachineCreate) SetScenarios(s string) *MachineCreate {
	mc.mutation.SetScenarios(s)
	return mc
}

// SetNillableScenarios sets the "scenarios" field if the given value is not nil.
func (mc *MachineCreate) SetNillableScenarios(s *string) *MachineCreate {
	if s != nil {
		mc.SetScenarios(*s)
	}
	return mc
}

// SetVersion sets the "version" field.
func (mc *MachineCreate) SetVersion(s string) *MachineCreate {
	mc.mutation.SetVersion(s)
	return mc
}

// SetNillableVersion sets the "version" field if the given value is not nil.
func (mc *MachineCreate) SetNillableVersion(s *string) *MachineCreate {
	if s != nil {
		mc.SetVersion(*s)
	}
	return mc
}

// SetIsValidated sets the "isValidated" field.
func (mc *MachineCreate) SetIsValidated(b bool) *MachineCreate {
	mc.mutation.SetIsValidated(b)
	return mc
}

// SetNillableIsValidated sets the "isValidated" field if the given value is not nil.
func (mc *MachineCreate) SetNillableIsValidated(b *bool) *MachineCreate {
	if b != nil {
		mc.SetIsValidated(*b)
	}
	return mc
}

// SetStatus sets the "status" field.
func (mc *MachineCreate) SetStatus(s string) *MachineCreate {
	mc.mutation.SetStatus(s)
	return mc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (mc *MachineCreate) SetNillableStatus(s *string) *MachineCreate {
	if s != nil {
		mc.SetStatus(*s)
	}
	return mc
}

// AddAlertIDs adds the "alerts" edge to the Alert entity by IDs.
func (mc *MachineCreate) AddAlertIDs(ids ...int) *MachineCreate {
	mc.mutation.AddAlertIDs(ids...)
	return mc
}

// AddAlerts adds the "alerts" edges to the Alert entity.
func (mc *MachineCreate) AddAlerts(a ...*Alert) *MachineCreate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return mc.AddAlertIDs(ids...)
}

// Mutation returns the MachineMutation object of the builder.
func (mc *MachineCreate) Mutation() *MachineMutation {
	return mc.mutation
}

// Save creates the Machine in the database.
func (mc *MachineCreate) Save(ctx context.Context) (*Machine, error) {
	var (
		err  error
		node *Machine
	)
	mc.defaults()
	if len(mc.hooks) == 0 {
		if err = mc.check(); err != nil {
			return nil, err
		}
		node, err = mc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MachineMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = mc.check(); err != nil {
				return nil, err
			}
			mc.mutation = mutation
			if node, err = mc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(mc.hooks) - 1; i >= 0; i-- {
			if mc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = mc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, mc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (mc *MachineCreate) SaveX(ctx context.Context) *Machine {
	v, err := mc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mc *MachineCreate) Exec(ctx context.Context) error {
	_, err := mc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mc *MachineCreate) ExecX(ctx context.Context) {
	if err := mc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mc *MachineCreate) defaults() {
	if _, ok := mc.mutation.CreatedAt(); !ok {
		v := machine.DefaultCreatedAt()
		mc.mutation.SetCreatedAt(v)
	}
	if _, ok := mc.mutation.UpdatedAt(); !ok {
		v := machine.DefaultUpdatedAt()
		mc.mutation.SetUpdatedAt(v)
	}
	if _, ok := mc.mutation.LastPush(); !ok {
		v := machine.DefaultLastPush()
		mc.mutation.SetLastPush(v)
	}
	if _, ok := mc.mutation.LastHeartbeat(); !ok {
		v := machine.DefaultLastHeartbeat()
		mc.mutation.SetLastHeartbeat(v)
	}
	if _, ok := mc.mutation.IsValidated(); !ok {
		v := machine.DefaultIsValidated
		mc.mutation.SetIsValidated(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mc *MachineCreate) check() error {
	if _, ok := mc.mutation.MachineId(); !ok {
		return &ValidationError{Name: "machineId", err: errors.New(`ent: missing required field "Machine.machineId"`)}
	}
	if _, ok := mc.mutation.Password(); !ok {
		return &ValidationError{Name: "password", err: errors.New(`ent: missing required field "Machine.password"`)}
	}
	if _, ok := mc.mutation.IpAddress(); !ok {
		return &ValidationError{Name: "ipAddress", err: errors.New(`ent: missing required field "Machine.ipAddress"`)}
	}
	if v, ok := mc.mutation.Scenarios(); ok {
		if err := machine.ScenariosValidator(v); err != nil {
			return &ValidationError{Name: "scenarios", err: fmt.Errorf(`ent: validator failed for field "Machine.scenarios": %w`, err)}
		}
	}
	if _, ok := mc.mutation.IsValidated(); !ok {
		return &ValidationError{Name: "isValidated", err: errors.New(`ent: missing required field "Machine.isValidated"`)}
	}
	return nil
}

func (mc *MachineCreate) sqlSave(ctx context.Context) (*Machine, error) {
	_node, _spec := mc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (mc *MachineCreate) createSpec() (*Machine, *sqlgraph.CreateSpec) {
	var (
		_node = &Machine{config: mc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: machine.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: machine.FieldID,
			},
		}
	)
	if value, ok := mc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: machine.FieldCreatedAt,
		})
		_node.CreatedAt = &value
	}
	if value, ok := mc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: machine.FieldUpdatedAt,
		})
		_node.UpdatedAt = &value
	}
	if value, ok := mc.mutation.LastPush(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: machine.FieldLastPush,
		})
		_node.LastPush = &value
	}
	if value, ok := mc.mutation.LastHeartbeat(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: machine.FieldLastHeartbeat,
		})
		_node.LastHeartbeat = &value
	}
	if value, ok := mc.mutation.MachineId(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: machine.FieldMachineId,
		})
		_node.MachineId = value
	}
	if value, ok := mc.mutation.Password(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: machine.FieldPassword,
		})
		_node.Password = value
	}
	if value, ok := mc.mutation.IpAddress(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: machine.FieldIpAddress,
		})
		_node.IpAddress = value
	}
	if value, ok := mc.mutation.Scenarios(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: machine.FieldScenarios,
		})
		_node.Scenarios = value
	}
	if value, ok := mc.mutation.Version(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: machine.FieldVersion,
		})
		_node.Version = value
	}
	if value, ok := mc.mutation.IsValidated(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: machine.FieldIsValidated,
		})
		_node.IsValidated = value
	}
	if value, ok := mc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: machine.FieldStatus,
		})
		_node.Status = value
	}
	if nodes := mc.mutation.AlertsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   machine.AlertsTable,
			Columns: []string{machine.AlertsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: alert.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// MachineCreateBulk is the builder for creating many Machine entities in bulk.
type MachineCreateBulk struct {
	config
	builders []*MachineCreate
}

// Save creates the Machine entities in the database.
func (mcb *MachineCreateBulk) Save(ctx context.Context) ([]*Machine, error) {
	specs := make([]*sqlgraph.CreateSpec, len(mcb.builders))
	nodes := make([]*Machine, len(mcb.builders))
	mutators := make([]Mutator, len(mcb.builders))
	for i := range mcb.builders {
		func(i int, root context.Context) {
			builder := mcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MachineMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, mcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, mcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mcb *MachineCreateBulk) SaveX(ctx context.Context) []*Machine {
	v, err := mcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mcb *MachineCreateBulk) Exec(ctx context.Context) error {
	_, err := mcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mcb *MachineCreateBulk) ExecX(ctx context.Context) {
	if err := mcb.Exec(ctx); err != nil {
		panic(err)
	}
}
