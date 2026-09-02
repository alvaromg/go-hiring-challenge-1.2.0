package monitor

import (
	"context"
	"io"

	"github.com/mytheresa/go-hiring-challenge/shared"
	"github.com/sirupsen/logrus"
)

var _ shared.Monitor = (*noopMonitor)(nil)
var _ shared.Logger = (*noopLogger)(nil)

// noopMonitor and noopLogger satisfy shared.Monitor/shared.Logger without
// producing any output, so test runs stay quiet.
type noopMonitor struct {
	logger shared.Logger
}

func NewNoopMonitor() *noopMonitor {
	return &noopMonitor{logger: newNoopLogger()}
}

func (m *noopMonitor) Logger() shared.Logger {
	return m.logger
}

type noopLogger struct {
	inner *logrus.Logger
}

func newNoopLogger() *noopLogger {
	l := logrus.New()
	l.SetOutput(io.Discard)

	return &noopLogger{inner: l}
}

func (l *noopLogger) Infof(format string, args ...any)  {}
func (l *noopLogger) Printf(format string, args ...any) {}
func (l *noopLogger) Errorf(format string, args ...any) {}
func (l *noopLogger) Debugf(format string, args ...any) {}
func (l *noopLogger) Tracef(format string, args ...any) {}
func (l *noopLogger) Warnf(format string, args ...any)  {}

func (l *noopLogger) WithFields(fields map[string]any) *logrus.Entry {
	return l.inner.WithFields(fields)
}

func (l *noopLogger) WithContext(ctx context.Context) *logrus.Entry {
	return l.inner.WithContext(ctx)
}
