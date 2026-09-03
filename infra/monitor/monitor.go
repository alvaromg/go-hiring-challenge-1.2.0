package monitor

import (
	"github.com/mytheresa/go-hiring-challenge/shared"
)

type monitor struct {
	logger shared.Logger
}

func (m *monitor) Logger() shared.Logger {
	return m.logger
}

func NewMonitor(log shared.Logger) *monitor {
	return &monitor{
		logger: log,
	}
}
