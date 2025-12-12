package services

import (
	"log"

	"github.com/flightpath-dev/flightpath/internal/config"
)

// Holds shared context for all services.
// This provides a clean way to pass common dependencies (config, logger, etc.)
// to service constructors without requiring multiple parameters.
type ServiceContext struct {
	Config *config.Config
	Logger *log.Logger
}
