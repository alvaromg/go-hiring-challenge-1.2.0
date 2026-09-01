package monitor

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
)

type logger struct {
	inner *logrus.Logger
}

// NewLogger creates a new logger with the specified level and formatting options.
// Use LoggerOption functions to configure additional features like OpenTelemetry integration.
//
// Example:
//
//	logger, err := logger.NewLogger("info", false, logger.WithOtelLogger(otelLogger))
func NewLogger(level string, pretty bool) (*logger, error) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	l := &logger{
		logrus.New(),
	}

	l.inner.SetLevel(logLevel)
	// l.inner.AddHook(newLogrusTraceHook(cfg.otelLogger))
	l.inner.SetFormatter(NewLogFormatter(pretty))

	return l, nil
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.inner.Logf(logrus.InfoLevel, format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.inner.Logf(logrus.InfoLevel, format, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.inner.Logf(logrus.DebugLevel, format, args...)
}

func (l *logger) Tracef(format string, args ...interface{}) {
	l.inner.Logf(logrus.TraceLevel, format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.inner.Logf(logrus.WarnLevel, format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.inner.Logf(logrus.ErrorLevel, format, args...)
}

func (l *logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	entry := l.inner.WithFields(fields)
	return entry
}

func (l *logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.inner.WithContext(ctx)
	return entry
}

func NewLogFormatter(pretty bool) logrus.Formatter {
	return UTCFormatter{
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
			PrettyPrint:     pretty,
		},
	}
}

type UTCFormatter struct {
	Formatter logrus.Formatter
}

func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}
