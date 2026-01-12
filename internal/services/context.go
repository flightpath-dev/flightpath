package services

import (
	"github.com/bluenviron/gomavlib/v3"
	"github.com/flightpath-dev/flightpath/internal/config"
	"github.com/flightpath-dev/flightpath/internal/logger"
	"github.com/flightpath-dev/flightpath/internal/mavlink"
)

// Holds shared context for all services.
// This provides a clean way to pass common dependencies (config, logger, etc.)
// to service constructors without requiring multiple parameters.
type ServiceContext struct {
	Config            *config.Config
	Logger            *logger.Logger
	Node              *gomavlib.Node
	MessageReceiver   *mavlink.MAVLinkMessageReceiver
	CommandDispatcher *mavlink.MAVLinkCommandDispatcher
}
