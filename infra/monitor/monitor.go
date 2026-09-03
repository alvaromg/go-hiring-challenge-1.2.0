package monitor

import (
	"github.com/mytheresa/go-hiring-challenge/shared"
	"go.opentelemetry.io/otel/trace"
)

type monitor struct {
	logger shared.Logger
	tracer trace.Tracer
}

func (m *monitor) Logger() shared.Logger {
	return m.logger
}

func (m *monitor) Tracer() trace.Tracer {
	return m.tracer
}

func NewMonitor(log shared.Logger, tracer trace.Tracer) *monitor {
	return &monitor{
		logger: log,
		tracer: tracer,
	}
}
