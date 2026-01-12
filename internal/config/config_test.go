package config

import (
	"testing"

	"github.com/bluenviron/gomavlib/v3"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	// Verify server defaults
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default host '0.0.0.0', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if len(cfg.Server.CORSOrigins) != 1 || cfg.Server.CORSOrigins[0] != "http://localhost:3000" {
		t.Errorf("expected default CORS origins ['http://localhost:3000'], got %v", cfg.Server.CORSOrigins)
	}

	// Verify MAVLink defaults
	if cfg.MAVLink.Endpoint == nil {
		t.Error("expected default MAVLink endpoint to be set")
	}
	udpEndpoint, ok := cfg.MAVLink.Endpoint.(gomavlib.EndpointUDPServer)
	if !ok {
		t.Errorf("expected default endpoint to be EndpointUDPServer, got %T", cfg.MAVLink.Endpoint)
	}
	if udpEndpoint.Address != "0.0.0.0:14550" {
		t.Errorf("expected default UDP address '0.0.0.0:14550', got '%s'", udpEndpoint.Address)
	}
}

func TestConfig_Validate_ValidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"minimum valid port", 1},
		{"common port", 8080},
		{"maximum valid port", 65535},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Default()
			cfg.Server.Port = tt.port
			if err := cfg.Validate(); err != nil {
				t.Errorf("expected valid config for port %d, got error: %v", tt.port, err)
			}
		})
	}
}

func TestConfig_Validate_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"zero port", 0},
		{"negative port", -1},
		{"port too high", 65536},
		{"port way too high", 100000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Default()
			cfg.Server.Port = tt.port
			if err := cfg.Validate(); err == nil {
				t.Errorf("expected error for invalid port %d, got nil", tt.port)
			}
		})
	}
}

func TestMAVLinkConfig_Validate_Serial(t *testing.T) {
	t.Run("valid serial config", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid serial config, got error: %v", err)
		}
	})

	t.Run("missing device", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointSerial{
				Device: "",
				Baud:   57600,
			},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for missing device, got nil")
		}
	})

	t.Run("invalid baud rate", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   0,
			},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for invalid baud rate, got nil")
		}
	})

	t.Run("negative baud rate", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   -1,
			},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for negative baud rate, got nil")
		}
	})
}

func TestMAVLinkConfig_Validate_UDP(t *testing.T) {
	t.Run("valid UDP server config", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointUDPServer{Address: "0.0.0.0:14550"},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid UDP server config, got error: %v", err)
		}
	})

	t.Run("missing UDP server address", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointUDPServer{Address: ""},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for missing UDP server address, got nil")
		}
	})

	t.Run("valid UDP client config", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointUDPClient{Address: "127.0.0.1:14550"},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid UDP client config, got error: %v", err)
		}
	})

	t.Run("missing UDP client address", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointUDPClient{Address: ""},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for missing UDP client address, got nil")
		}
	})
}

func TestMAVLinkConfig_Validate_TCP(t *testing.T) {
	t.Run("valid TCP server config", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointTCPServer{Address: "0.0.0.0:5760"},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid TCP server config, got error: %v", err)
		}
	})

	t.Run("missing TCP server address", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointTCPServer{Address: ""},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for missing TCP server address, got nil")
		}
	})

	t.Run("valid TCP client config", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointTCPClient{Address: "127.0.0.1:5760"},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid TCP client config, got error: %v", err)
		}
	})

	t.Run("missing TCP client address", func(t *testing.T) {
		cfg := &MAVLinkConfig{
			Endpoint: gomavlib.EndpointTCPClient{Address: ""},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for missing TCP client address, got nil")
		}
	})
}

func TestMAVLinkConfig_Validate_NilEndpoint(t *testing.T) {
	cfg := &MAVLinkConfig{
		Endpoint: nil,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected nil endpoint to be valid, got error: %v", err)
	}
}

func TestServerAddr(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     int
		expected string
	}{
		{"default config", "0.0.0.0", 8080, "0.0.0.0:8080"},
		{"localhost", "localhost", 3000, "localhost:3000"},
		{"specific IP", "192.168.1.1", 9090, "192.168.1.1:9090"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					Host: tt.host,
					Port: tt.port,
				},
			}
			if got := cfg.ServerAddr(); got != tt.expected {
				t.Errorf("ServerAddr() = '%s', expected '%s'", got, tt.expected)
			}
		})
	}
}
