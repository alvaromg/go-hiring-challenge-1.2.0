package monitor

import (
	"github.com/mytheresa/go-hiring-challenge/shared"
)

type monitor struct {
	logger shared.Logger
	// tracer libmonitor.Tracer
	// meter  libmonitor.Meter
}

func (m *monitor) Logger() shared.Logger {
	return m.logger
}

// func (m *monitor) Tracer() libmonitor.Tracer {
// 	return m.tracer
// }

// func (m *monitor) Meter() libmonitor.Meter {
// 	return m.meter
// }

func NewMonitor(log shared.Logger /*tr libmonitor.Tracer, mt libmonitor.Meter*/) *monitor {
	return &monitor{
		logger: log,
		// tracer: tr,
		// meter:  mt,
	}
}
