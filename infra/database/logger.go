package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

func parseLogLevel(level string) (logger.LogLevel, error) {

	level = strings.ToLower(level)

	switch level {
	case "silent":
		return logger.Silent, nil
	case "error":
		return logger.Error, nil
	case "warn":
		return logger.Warn, nil
	case "info":
		return logger.Info, nil
	case "":
		return logger.Silent, fmt.Errorf("database log level is empty")
	default:
		return logger.Silent, fmt.Errorf("unknown database log level %q", level)
	}
}

// jsonLogWriter adapts gorm's logger.Writer interface to emit JSON-formatted log lines.
type jsonLogWriter struct {
	inner *logrus.Logger
}

func newJSONLogWriter() *jsonLogWriter {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
	})

	return &jsonLogWriter{inner: l}
}

func (w *jsonLogWriter) Printf(format string, args ...any) {
	w.inner.Infof(format, args...)
}
