package shared

type Monitor interface {
	Logger() Logger
	// Tracer() libmonitor.Tracer
	// Meter() libmonitor.Meter
}

type Logger interface {
	Infof(format string, args ...any)
	Printf(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
	Tracef(format string, args ...any)
	Warnf(format string, args ...any)
	// WithFields(map[string]interface{}) *logrus.Entry
	// WithContext(ctx context.Context) *logrus.Entry
}
