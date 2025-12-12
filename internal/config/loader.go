package config

import (
	"log"
	"os"
	"strconv"
	"strings"
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
//   - FLIGHTPATH_PORT: Server port (integer, 1-65535)
//   - FLIGHTPATH_HOST: Server host (string, default: "0.0.0.0")
//   - FLIGHTPATH_CORS_ORIGINS: Comma-separated list of allowed CORS origins
//
// Example usage:
//
//	export FLIGHTPATH_PORT=3000
//	./server
func Load() *Config {
	cfg := Default()

	// Override with environment variables if present
	if port := os.Getenv("FLIGHTPATH_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}

	if host := os.Getenv("FLIGHTPATH_HOST"); host != "" {
		cfg.Server.Host = host
	}

	if corsOrigins := os.Getenv("FLIGHTPATH_CORS_ORIGINS"); corsOrigins != "" {
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

	// Validate configuration (fail-fast)
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	return cfg
}
