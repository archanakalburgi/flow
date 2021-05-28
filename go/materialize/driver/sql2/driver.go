package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/estuary/flow/go/materialize/lifecycle"
	pf "github.com/estuary/flow/go/protocols/flow"
	pm "github.com/estuary/flow/go/protocols/materialize"
)

// Driver implements the pm.DriverServer interface.
type Driver struct {
	// NewEndpoint returns an Endpoint, which will be used to handle interactions with the database.
	NewEndpoint func(_ context.Context, name string, config json.RawMessage) (*Endpoint, error)
	// NewTransactor returns a Transactor ready for lifecycle.RunTransactions.
	NewTransactor func(*Endpoint, *pf.MaterializationSpec, *Fence) (lifecycle.Transactor, error)
}

var _ pm.DriverServer = &Driver{}

// Validate implements the DriverServer interface.
func (d *Driver) Validate(ctx context.Context, req *pm.ValidateRequest) (*pm.ValidateResponse, error) {
	endpoint, err := d.NewEndpoint(
		ctx,
		req.EndpointName,
		json.RawMessage(req.EndpointSpecJson),
	)
	if err != nil {
		return nil, fmt.Errorf("building endpoint: %w", err)
	} else if err = req.Collection.Validate(); err != nil {
		return nil, fmt.Errorf("validating collection: %w", err)
	}

	current, err := endpoint.LoadSpec(false)
	if err != nil {
		return nil, fmt.Errorf("loading current spec: %w", err)
	}

	var constraints map[string]*pm.Constraint
	if current != nil {
		constraints = ValidateMatchesExisting(current, req.Collection)
	} else {
		constraints = ValidateNewSQLProjections(req.Collection)
	}

	return &pm.ValidateResponse{
		Constraints:  constraints,
		ResourcePath: endpoint.TablePath,
	}, nil
}

// Apply implements the DriverServer interface.
func (d *Driver) Apply(ctx context.Context, req *pm.ApplyRequest) (*pm.ApplyResponse, error) {
	endpoint, err := d.NewEndpoint(
		ctx,
		req.Materialization.EndpointName,
		json.RawMessage(req.Materialization.EndpointSpecJson),
	)
	if err != nil {
		return nil, fmt.Errorf("building endpoint: %w", err)
	} else if err = req.Materialization.Validate(); err != nil {
		return nil, fmt.Errorf("validating materialization: %w", err)
	}

	current, err := endpoint.LoadSpec(false)
	if err != nil {
		return nil, fmt.Errorf("loading current spec: %w", err)
	}

	var constraints map[string]*pm.Constraint
	if current != nil {
		constraints = ValidateMatchesExisting(current, &req.Materialization.Collection)
	} else {
		constraints = ValidateNewSQLProjections(&req.Materialization.Collection)
	}

	// Validate the request materialization is a valid solution for its own constraints.
	if err = ValidateSelectedFields(constraints, req.Materialization); err != nil {
		return nil, fmt.Errorf("re-validating materialization: %w", err)
	}

	// We don't handle any form of schema migrations, so we require that the list of
	// fields in the request is identical to the current fields. doValidate doesn't handle that
	// because the list of fields isn't known until Apply is called.
	if current != nil && !req.Materialization.FieldSelection.Equal(current.FieldSelection) {
		return nil, fmt.Errorf(
			"the set of fields in the request differs from the existing fields,"+
				"which is disallowed because this driver does not perform schema migrations. "+
				"Request fields: %s , Existing fields: %s",
			req.Materialization.FieldSelection.String(),
			current.FieldSelection.String(),
		)
	}

	// If schema was already applied, there's no further work to be done.
	if current != nil {
		return new(pm.ApplyResponse), nil
	}

	// Generate statements to be applied.
	statements, err := generateApplyStatements(endpoint, req.Materialization)
	if err != nil {
		return nil, err
	}

	// Apply the statements if not in DryRun.
	if !req.DryRun {
		if err = endpoint.ApplyStatements(statements); err != nil {
			return nil, fmt.Errorf("applying schema updates: %w", err)
		}
	}

	// Build and return a description of what happened (or would have happened).
	return &pm.ApplyResponse{
		ActionDescription: fmt.Sprintf(
			"%s\nBEGIN;\n%s\nCOMMIT;\n",
			endpoint.Generator.Comment(fmt.Sprintf(
				"Generated by Flow for materializing collection '%s'\nto table: %s",
				req.Materialization.Collection.Collection,
				endpoint.TargetName(),
			)),
			strings.Join(statements, "\n\n"),
		),
	}, nil
}

// Transactions implements the DriverServer interface.
func (d *Driver) Transactions(stream pm.Driver_TransactionsServer) error {
	var open, err = stream.Recv()
	if err != nil {
		return fmt.Errorf("read Open: %w", err)
	} else if open.Open == nil {
		return fmt.Errorf("expected Open, got %#v", open)
	}

	endpoint, err := d.NewEndpoint(
		stream.Context(),
		open.Open.Materialization.EndpointName,
		json.RawMessage(open.Open.Materialization.EndpointSpecJson),
	)
	if err != nil {
		return fmt.Errorf("building endpoint: %w", err)
	}

	spec, err := endpoint.LoadSpec(true)
	if err != nil {
		return fmt.Errorf("loading materialization spec: %w", err)
	}
	fence, err := endpoint.NewFence(open.Open.ShardFqn)
	if err != nil {
		return fmt.Errorf("installing fence: %w", err)
	}
	transactor, err := d.NewTransactor(endpoint, spec, fence)
	if err != nil {
		return err
	}

	if err = stream.Send(&pm.TransactionResponse{
		Opened: &pm.TransactionResponse_Opened{
			FlowCheckpoint: fence.Checkpoint,
			DeltaUpdates:   false,
		},
	}); err != nil {
		return fmt.Errorf("sending Opened: %w", err)
	}

	return lifecycle.RunTransactions(stream, transactor, fence.LogEntry())
}
