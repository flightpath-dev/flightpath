package config

import (
	"fmt"
)

// Config holds all application configuration.
//
// Configuration Strategy:
//
// This package follows the 12-factor app methodology for configuration:
// - Configuration is provided via environment variables
// - Sensible defaults are provided for developer-friendly local development
// - Production deployments override defaults via environment variables
// - Configuration is validated at startup (fail-fast)
// - Type-safe structured types instead of maps/strings
//
// Environment Variable Naming:
// All environment variables are prefixed with FLIGHTPATH_ to avoid conflicts.
// Examples:
//   - FLIGHTPATH_PORT: Server port (default: 8080)
//   - FLIGHTPATH_HOST: Server host (default: 0.0.0.0)
//   - FLIGHTPATH_CORS_ORIGINS: Comma-separated CORS origins (default: localhost:5173,localhost:3000)
//
// This follows the convention over configuration principle: sensible defaults
// with optional overrides for production environments.
type Config struct {
	Server ServerConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host        string
	Port        int
	CORSOrigins []string
}

// Default returns a Config with sensible defaults for local development.
// These defaults work out of the box without any configuration.
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
			CORSOrigins: []string{
				"http://localhost:5173", // Vite dev server
				"http://localhost:3000",
			},
		},
	}
}

// Validate checks if the configuration is valid and returns an error if invalid.
// This implements fail-fast validation to catch configuration errors at startup.
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be between 1 and 65535)", c.Server.Port)
	}

	return nil
}

// ServerAddr returns the server address as "host:port" format.
// This is a convenience method for use with http.ListenAndServe.
func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
