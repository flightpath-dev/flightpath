package services

import (
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/flightpath-dev/flightpath/internal/config"
)

// Holds shared context for all services.
// This provides a clean way to pass common dependencies (config, logger, etc.)
// to service constructors without requiring multiple parameters.
type ServiceContext struct {
	Config     *config.Config
	Logger     *log.Logger
	Node       *gomavlib.Node
	Dispatcher *MessageDispatcher
}
