package shared

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Monitor interface {
	Logger() Logger
	Tracer() trace.Tracer
}

type Logger interface {
	Infof(format string, args ...any)
	Printf(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
	Tracef(format string, args ...any)
	Warnf(format string, args ...any)
	WithFields(map[string]any) *logrus.Entry
	WithContext(ctx context.Context) *logrus.Entry
}
