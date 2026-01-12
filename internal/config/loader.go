package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/flightpath-dev/flightpath/internal/logger"
)

// parseLogLevel
// Parses log level string to LogLevel enum.
// Valid values: DEBUG, INFO, WARN, ERROR (case-insensitive)
// Defaults to LogLevelInfo if invalid value.
func parseLogLevel(levelStr string) logger.LogLevel {
	levelStr = strings.ToUpper(levelStr)
	switch levelStr {
	case "DEBUG":
		return logger.LogLevelDebug
	case "INFO":
		return logger.LogLevelInfo
	case "WARN":
		return logger.LogLevelWarn
	case "ERROR":
		return logger.LogLevelError
	default:
		return logger.LogLevelInfo
	}
}

// parseLogFormat
// Parses log format string to LogFormat enum.
// Valid values: text, json (case-insensitive)
// Defaults to LogFormatText if invalid value.
func parseLogFormat(formatStr string) logger.LogFormat {
	formatStr = strings.ToLower(formatStr)
	switch formatStr {
	case "json":
		return logger.LogFormatJSON
	case "text":
		fallthrough
	default:
		return logger.LogFormatText
	}
}

// Load loads configuration from environment variables, falling back to defaults
// for any missing values. This implements the 12-factor app configuration pattern.
//
// The configuration is loaded in the following order:
// 1. Start with Default() configuration (developer-friendly defaults)
// 2. Override with environment variables if present
// 3. Validate the final configuration (fail-fast)
//
// Environment Variables:
//   - FLIGHTPATH_GRPC_PORT: gRPC server port (integer, 1-65535)
//   - FLIGHTPATH_GRPC_HOST: gRPC server host (string, default: "0.0.0.0")
//   - FLIGHTPATH_GRPC_CORS_ORIGINS: Comma-separated list of allowed CORS origins
//   - FLIGHTPATH_MAVLINK_ENDPOINT_TYPE: MAVLink endpoint type (serial, udp-server, udp-client, tcp-server, tcp-client)
//   - FLIGHTPATH_MAVLINK_SERIAL_DEVICE: Serial device path (required if type is "serial")
//   - FLIGHTPATH_MAVLINK_SERIAL_BAUD: Serial baud rate (default: 57600, required if type is "serial")
//   - FLIGHTPATH_MAVLINK_UDP_ADDRESS: UDP address in "host:port" format (default: "0.0.0.0:14550")
//   - FLIGHTPATH_MAVLINK_TCP_ADDRESS: TCP address in "host:port" format (required if type is "tcp-server" or "tcp-client")
//   - FLIGHTPATH_LOG_LEVEL: Log level (string, case-insensitive, default: "INFO")
//   - FLIGHTPATH_LOG_FORMAT: Log format (string, case-insensitive, default: "text")
//
// Example usage:
//
//	export FLIGHTPATH_GRPC_PORT=3000
//	./server
func Load() (*Config, error) {
	cfg := Default()

	// Override with environment variables if present
	if port := os.Getenv("FLIGHTPATH_GRPC_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}

	if host := os.Getenv("FLIGHTPATH_GRPC_HOST"); host != "" {
		cfg.Server.Host = host
	}

	if corsOrigins := os.Getenv("FLIGHTPATH_GRPC_CORS_ORIGINS"); corsOrigins != "" {
		// Split comma-separated values and trim whitespace
		origins := strings.Split(corsOrigins, ",")
		cfg.Server.CORSOrigins = make([]string, 0, len(origins))
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				cfg.Server.CORSOrigins = append(cfg.Server.CORSOrigins, trimmed)
			}
		}
	}

	// Load log level from environment variable
	if logLevel := os.Getenv("FLIGHTPATH_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = parseLogLevel(logLevel)
	}

	// Load log format from environment variable
	if logFormat := os.Getenv("FLIGHTPATH_LOG_FORMAT"); logFormat != "" {
		cfg.LogFormat = parseLogFormat(logFormat)
	}

	// Load MAVLink configuration from environment variables
	if err := loadMAVLinkConfig(cfg); err != nil {
		return nil, err
	}

	// Validate configuration (fail-fast)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Log configuration
	logConfig(cfg)

	return cfg, nil
}

// loadMAVLinkConfig
// Loads MAVLink configuration from environment variables.
//
// Only overrides defaults if environment variables are present.
// If FLIGHTPATH_MAVLINK_ENDPOINT_TYPE is set, all required parameters for that
// endpoint type must be provided via environment variables.
// Returns an error if the endpoint type is set but required parameters are missing.
func loadMAVLinkConfig(cfg *Config) error {
	endpointType := os.Getenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE")
	if endpointType == "" {
		// No override - use default from Default()
		return nil
	}

	switch endpointType {
	case "serial":
		device := os.Getenv("FLIGHTPATH_MAVLINK_SERIAL_DEVICE")
		baudStr := os.Getenv("FLIGHTPATH_MAVLINK_SERIAL_BAUD")

		if device == "" {
			return fmt.Errorf("FLIGHTPATH_MAVLINK_SERIAL_DEVICE is required when endpoint type is 'serial'")
		}
		if baudStr == "" {
			return fmt.Errorf("FLIGHTPATH_MAVLINK_SERIAL_BAUD is required when endpoint type is 'serial'")
		}

		baud, err := strconv.Atoi(baudStr)
		if err != nil {
			return fmt.Errorf("FLIGHTPATH_MAVLINK_SERIAL_BAUD must be a valid integer: %w", err)
		}
		if baud <= 0 {
			return fmt.Errorf("FLIGHTPATH_MAVLINK_SERIAL_BAUD must be greater than 0, got %d", baud)
		}

		cfg.MAVLink.Endpoint = gomavlib.EndpointSerial{
			Device: device,
			Baud:   baud,
		}

	case "udp-server":
		address := os.Getenv("FLIGHTPATH_MAVLINK_UDP_ADDRESS")
		if address == "" {
			return fmt.Errorf("FLIGHTPATH_MAVLINK_UDP_ADDRESS is required when endpoint type is 'udp-server'")
		}
		cfg.MAVLink.Endpoint = gomavlib.EndpointUDPServer{Address: address}

	case "udp-client":
		address := os.Getenv("FLIGHTPATH_MAVLINK_UDP_ADDRESS")
		if address == "" {
			return fmt.Errorf("FLIGHTPATH_MAVLINK_UDP_ADDRESS is required when endpoint type is 'udp-client'")
		}
		cfg.MAVLink.Endpoint = gomavlib.EndpointUDPClient{Address: address}

	case "tcp-server":
		address := os.Getenv("FLIGHTPATH_MAVLINK_TCP_ADDRESS")
		if address == "" {
			return fmt.Errorf("FLIGHTPATH_MAVLINK_TCP_ADDRESS is required when endpoint type is 'tcp-server'")
		}
		cfg.MAVLink.Endpoint = gomavlib.EndpointTCPServer{Address: address}

	case "tcp-client":
		address := os.Getenv("FLIGHTPATH_MAVLINK_TCP_ADDRESS")
		if address == "" {
			return fmt.Errorf("FLIGHTPATH_MAVLINK_TCP_ADDRESS is required when endpoint type is 'tcp-client'")
		}
		cfg.MAVLink.Endpoint = gomavlib.EndpointTCPClient{Address: address}

	default:
		return fmt.Errorf("unknown FLIGHTPATH_MAVLINK_ENDPOINT_TYPE: %q (valid types: serial, udp-server, udp-client, tcp-server, tcp-client)", endpointType)
	}

	return nil
}

// logConfig
// Logs the loaded configuration for debugging and transparency.
// Shows server configuration and MAVLink endpoint details.
func logConfig(cfg *Config) {
	configLogger := logger.New(cfg.LogLevel, cfg.LogFormat).WithPrefix("config")

	// Convert log level to string for logging
	logLevelStr := "INFO"
	switch cfg.LogLevel {
	case logger.LogLevelDebug:
		logLevelStr = "DEBUG"
	case logger.LogLevelInfo:
		logLevelStr = "INFO"
	case logger.LogLevelWarn:
		logLevelStr = "WARN"
	case logger.LogLevelError:
		logLevelStr = "ERROR"
	}

	configLogger.Info("Configuration loaded",
		"server_host", cfg.Server.Host,
		"server_port", cfg.Server.Port,
		"log_level", logLevelStr,
	)

	if len(cfg.Server.CORSOrigins) > 0 {
		configLogger.Info("CORS origins configured", "origins", strings.Join(cfg.Server.CORSOrigins, ", "))
	}

	if cfg.MAVLink.Endpoint == nil {
		configLogger.Info("MAVLink endpoint not configured")
		return
	}

	// Log endpoint details based on type
	switch endpoint := cfg.MAVLink.Endpoint.(type) {
	case gomavlib.EndpointSerial:
		configLogger.Info("MAVLink endpoint configured",
			"type", "serial",
			"device", endpoint.Device,
			"baud", endpoint.Baud,
		)
	case gomavlib.EndpointUDPServer:
		configLogger.Info("MAVLink endpoint configured",
			"type", "udp-server",
			"address", endpoint.Address,
		)
	case gomavlib.EndpointUDPClient:
		configLogger.Info("MAVLink endpoint configured",
			"type", "udp-client",
			"address", endpoint.Address,
		)
	case gomavlib.EndpointTCPServer:
		configLogger.Info("MAVLink endpoint configured",
			"type", "tcp-server",
			"address", endpoint.Address,
		)
	case gomavlib.EndpointTCPClient:
		configLogger.Info("MAVLink endpoint configured",
			"type", "tcp-client",
			"address", endpoint.Address,
		)
	default:
		configLogger.Warn("MAVLink endpoint configured with unknown type", "type", fmt.Sprintf("%T", endpoint))
	}
}
