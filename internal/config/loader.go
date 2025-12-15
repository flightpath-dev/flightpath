package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bluenviron/gomavlib/v3"
)

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

	// Load MAVLink configuration from environment variables
	loadMAVLinkConfig(cfg)

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
// endpoint type must be provided via environment variables (no defaults used).
func loadMAVLinkConfig(cfg *Config) {
	endpointType := os.Getenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE")
	if endpointType == "" {
		// No override - use default from Default()
		return
	}

	switch endpointType {
	case "serial":
		device := os.Getenv("FLIGHTPATH_MAVLINK_SERIAL_DEVICE")
		baudStr := os.Getenv("FLIGHTPATH_MAVLINK_SERIAL_BAUD")

		if device == "" || baudStr == "" {
			// Required parameters missing - don't override
			return
		}

		baud, err := strconv.Atoi(baudStr)
		if err != nil || baud <= 0 {
			// Invalid baud rate - don't override
			return
		}

		cfg.MAVLink.Endpoint = gomavlib.EndpointSerial{
			Device: device,
			Baud:   baud,
		}

	case "udp-server":
		address := os.Getenv("FLIGHTPATH_MAVLINK_UDP_ADDRESS")
		if address == "" {
			// Required parameter missing - don't override
			return
		}
		cfg.MAVLink.Endpoint = gomavlib.EndpointUDPServer{Address: address}

	case "udp-client":
		address := os.Getenv("FLIGHTPATH_MAVLINK_UDP_ADDRESS")
		if address == "" {
			// Required parameter missing - don't override
			return
		}
		cfg.MAVLink.Endpoint = gomavlib.EndpointUDPClient{Address: address}

	case "tcp-server":
		address := os.Getenv("FLIGHTPATH_MAVLINK_TCP_ADDRESS")
		if address == "" {
			// Required parameter missing - don't override
			return
		}
		cfg.MAVLink.Endpoint = gomavlib.EndpointTCPServer{Address: address}

	case "tcp-client":
		address := os.Getenv("FLIGHTPATH_MAVLINK_TCP_ADDRESS")
		if address == "" {
			// Required parameter missing - don't override
			return
		}
		cfg.MAVLink.Endpoint = gomavlib.EndpointTCPClient{Address: address}
	}
}

// logConfig
// Logs the loaded configuration for debugging and transparency.
// Shows server configuration and MAVLink endpoint details.
func logConfig(cfg *Config) {
	log.Println("=== Configuration ===")
	log.Printf("Server: %s:%d", cfg.Server.Host, cfg.Server.Port)
	if len(cfg.Server.CORSOrigins) > 0 {
		log.Printf("CORS Origins: %s", strings.Join(cfg.Server.CORSOrigins, ", "))
	}

	if cfg.MAVLink.Endpoint == nil {
		log.Println("MAVLink: Not configured")
		return
	}

	// Log endpoint details based on type
	switch endpoint := cfg.MAVLink.Endpoint.(type) {
	case gomavlib.EndpointSerial:
		log.Printf("MAVLink: Serial - Device: %s, Baud: %d", endpoint.Device, endpoint.Baud)
	case gomavlib.EndpointUDPServer:
		log.Printf("MAVLink: UDP Server - Address: %s", endpoint.Address)
	case gomavlib.EndpointUDPClient:
		log.Printf("MAVLink: UDP Client - Address: %s", endpoint.Address)
	case gomavlib.EndpointTCPServer:
		log.Printf("MAVLink: TCP Server - Address: %s", endpoint.Address)
	case gomavlib.EndpointTCPClient:
		log.Printf("MAVLink: TCP Client - Address: %s", endpoint.Address)
	default:
		log.Printf("MAVLink: Unknown endpoint type: %T", endpoint)
	}
	log.Println("====================")
}
