package database

import (
	"errors"

	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// tracingPlugin wraps every SQL statement GORM sends to the database in its
// own OpenTelemetry span, so each query shows up in traces alongside the
// HTTP span (infra/rest/tracing.go) that triggered it.
type tracingPlugin struct {
	tracer trace.Tracer
}

func newTracingPlugin(tracer trace.Tracer) *tracingPlugin {
	return &tracingPlugin{tracer: tracer}
}

func (p *tracingPlugin) Name() string {
	return "otel:tracing"
}

// Initialize hooks a start/end callback pair around each of GORM's callback
// chains (create, query, update, delete, row, raw), which together cover
// every way GORM issues a SQL statement. Each pair is registered around the
// chain's own main callback (e.g. "gorm:query"), so the span covers exactly
// the work GORM does to build and run that statement.
func (p *tracingPlugin) Initialize(db *gorm.DB) error {
	create := db.Callback().Create()
	before, after := p.callbacks("INSERT")
	if err := create.Before("gorm:before_create").Register("otel:before_create", before); err != nil {
		return err
	}
	if err := create.After("gorm:after_create").Register("otel:after_create", after); err != nil {
		return err
	}

	query := db.Callback().Query()
	before, after = p.callbacks("SELECT")
	if err := query.Before("gorm:query").Register("otel:before_query", before); err != nil {
		return err
	}
	if err := query.After("gorm:after_query").Register("otel:after_query", after); err != nil {
		return err
	}

	update := db.Callback().Update()
	before, after = p.callbacks("UPDATE")
	if err := update.Before("gorm:before_update").Register("otel:before_update", before); err != nil {
		return err
	}
	if err := update.After("gorm:after_update").Register("otel:after_update", after); err != nil {
		return err
	}

	del := db.Callback().Delete()
	before, after = p.callbacks("DELETE")
	if err := del.Before("gorm:before_delete").Register("otel:before_delete", before); err != nil {
		return err
	}
	if err := del.After("gorm:after_delete").Register("otel:after_delete", after); err != nil {
		return err
	}

	row := db.Callback().Row()
	before, after = p.callbacks("SELECT")
	if err := row.Before("gorm:row").Register("otel:before_row", before); err != nil {
		return err
	}
	if err := row.After("gorm:row").Register("otel:after_row", after); err != nil {
		return err
	}

	raw := db.Callback().Raw()
	before, after = p.callbacks("SQL")
	if err := raw.Before("gorm:raw").Register("otel:before_raw", before); err != nil {
		return err
	}
	if err := raw.After("gorm:raw").Register("otel:after_raw", after); err != nil {
		return err
	}

	return nil
}

// callbacks returns the before/after handler pair for a given SQL operation.
// before starts the span and stashes it on the statement's context; after
// enriches it with the table and final SQL text (only known once the main
// callback has run) and ends it.
func (p *tracingPlugin) callbacks(operation string) (before, after func(*gorm.DB)) {
	before = func(tx *gorm.DB) {
		ctx, span := p.tracer.Start(tx.Statement.Context, operation, trace.WithSpanKind(trace.SpanKindClient))
		span.SetAttributes(semconv.DBSystemNamePostgreSQL, semconv.DBOperationName(operation))
		tx.Statement.Context = ctx
	}

	after = func(tx *gorm.DB) {
		span := trace.SpanFromContext(tx.Statement.Context)
		defer span.End()

		if table := tx.Statement.Table; table != "" {
			span.SetName(operation + " " + table)
			span.SetAttributes(semconv.DBCollectionName(table))
		}
		if sql := tx.Statement.SQL.String(); sql != "" {
			span.SetAttributes(semconv.DBQueryText(sql))
		}
		if err := tx.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}

	return before, after
}
