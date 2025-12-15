package config

import (
	"fmt"

	"github.com/bluenviron/gomavlib/v3"
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
//   - FLIGHTPATH_GRPC_PORT: gRPC server port (default: 8080)
//   - FLIGHTPATH_GRPC_HOST: gRPC server host (default: 0.0.0.0)
//   - FLIGHTPATH_GRPC_CORS_ORIGINS: Comma-separated CORS origins (default: localhost:5173,localhost:3000)
//
// This follows the convention over configuration principle: sensible defaults
// with optional overrides for production environments.
type Config struct {
	Server  ServerConfig
	MAVLink MAVLinkConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host        string
	Port        int
	CORSOrigins []string
}

// MAVLinkConfig holds MAVLink connection configuration.
// Uses gomavlib's EndpointConf interface directly, which provides a discriminated union
// pattern with type-safe endpoint configurations.
//
// gomavlib.EndpointConf is implemented by:
//   - gomavlib.EndpointSerial
//   - gomavlib.EndpointUDPServer / gomavlib.EndpointUDPClient
//   - gomavlib.EndpointTCPServer / gomavlib.EndpointTCPClient
//   - gomavlib.EndpointUDPBroadcast
//   - gomavlib.EndpointCustom / gomavlib.EndpointCustomServer / gomavlib.EndpointCustomClient
type MAVLinkConfig struct {
	Endpoint gomavlib.EndpointConf
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
		MAVLink: MAVLinkConfig{
			// Default to UDP server on port 14550 (standard PX4 SITL port)
			Endpoint: gomavlib.EndpointUDPServer{Address: "0.0.0.0:14550"},
		},
	}
}

// Validate checks if the configuration is valid and returns an error if invalid.
// This implements fail-fast validation to catch configuration errors at startup.
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be between 1 and 65535)", c.Server.Port)
	}

	// Validate MAVLink configuration
	if err := c.MAVLink.Validate(); err != nil {
		return fmt.Errorf("invalid MAVLink configuration: %w", err)
	}

	return nil
}

// Validate checks if the MAVLink configuration is valid.
// Uses type switch to validate the specific endpoint configuration type.
func (m *MAVLinkConfig) Validate() error {
	if m.Endpoint == nil {
		// nil endpoint is allowed (no MAVLink connection)
		return nil
	}

	switch cfg := m.Endpoint.(type) {
	case gomavlib.EndpointSerial:
		if cfg.Device == "" {
			return fmt.Errorf("serial device path is required")
		}
		if cfg.Baud <= 0 {
			return fmt.Errorf("serial baud rate must be greater than 0")
		}
	case gomavlib.EndpointUDPServer:
		if cfg.Address == "" {
			return fmt.Errorf("UDP server address is required")
		}
	case gomavlib.EndpointUDPClient:
		if cfg.Address == "" {
			return fmt.Errorf("UDP client address is required")
		}
	case gomavlib.EndpointTCPServer:
		if cfg.Address == "" {
			return fmt.Errorf("TCP server address is required")
		}
	case gomavlib.EndpointTCPClient:
		if cfg.Address == "" {
			return fmt.Errorf("TCP client address is required")
		}
	}
	return nil
}

// ServerAddr returns the server address as "host:port" format.
// This is a convenience method for use with http.ListenAndServe.
func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
