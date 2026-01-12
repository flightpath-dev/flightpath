package config

import (
	"testing"

	"github.com/bluenviron/gomavlib/v3"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear any environment variables that might affect the test
	t.Setenv("FLIGHTPATH_GRPC_PORT", "")
	t.Setenv("FLIGHTPATH_GRPC_HOST", "")
	t.Setenv("FLIGHTPATH_GRPC_CORS_ORIGINS", "")
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Verify defaults are applied
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default host '0.0.0.0', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
}

func TestLoad_OverridePort(t *testing.T) {
	t.Setenv("FLIGHTPATH_GRPC_PORT", "3000")
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Server.Port != 3000 {
		t.Errorf("expected port 3000, got %d", cfg.Server.Port)
	}
}

func TestLoad_OverrideHost(t *testing.T) {
	t.Setenv("FLIGHTPATH_GRPC_HOST", "127.0.0.1")
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected host '127.0.0.1', got '%s'", cfg.Server.Host)
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	t.Setenv("FLIGHTPATH_GRPC_PORT", "not-a-number")
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Invalid port string should be ignored, default used
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080 when invalid port string, got %d", cfg.Server.Port)
	}
}

func TestLoad_CORSOrigins(t *testing.T) {
	t.Run("single origin", func(t *testing.T) {
		t.Setenv("FLIGHTPATH_GRPC_CORS_ORIGINS", "http://example.com")
		t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}

		if len(cfg.Server.CORSOrigins) != 1 {
			t.Errorf("expected 1 CORS origin, got %d", len(cfg.Server.CORSOrigins))
		}
		if cfg.Server.CORSOrigins[0] != "http://example.com" {
			t.Errorf("expected 'http://example.com', got '%s'", cfg.Server.CORSOrigins[0])
		}
	})

	t.Run("multiple origins", func(t *testing.T) {
		t.Setenv("FLIGHTPATH_GRPC_CORS_ORIGINS", "http://localhost:3000,http://localhost:4000,http://example.com")
		t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}

		if len(cfg.Server.CORSOrigins) != 3 {
			t.Errorf("expected 3 CORS origins, got %d", len(cfg.Server.CORSOrigins))
		}
	})

	t.Run("origins with whitespace", func(t *testing.T) {
		t.Setenv("FLIGHTPATH_GRPC_CORS_ORIGINS", " http://localhost:3000 , http://localhost:4000 ")
		t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}

		if len(cfg.Server.CORSOrigins) != 2 {
			t.Errorf("expected 2 CORS origins, got %d", len(cfg.Server.CORSOrigins))
		}
		// Whitespace should be trimmed
		if cfg.Server.CORSOrigins[0] != "http://localhost:3000" {
			t.Errorf("expected trimmed origin 'http://localhost:3000', got '%s'", cfg.Server.CORSOrigins[0])
		}
	})

	t.Run("empty values filtered out", func(t *testing.T) {
		t.Setenv("FLIGHTPATH_GRPC_CORS_ORIGINS", "http://localhost:3000,,http://localhost:4000,")
		t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}

		if len(cfg.Server.CORSOrigins) != 2 {
			t.Errorf("expected 2 CORS origins (empty filtered), got %d: %v", len(cfg.Server.CORSOrigins), cfg.Server.CORSOrigins)
		}
	})
}

func TestLoad_MAVLinkSerial(t *testing.T) {
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "serial")
	t.Setenv("FLIGHTPATH_MAVLINK_SERIAL_DEVICE", "/dev/ttyUSB0")
	t.Setenv("FLIGHTPATH_MAVLINK_SERIAL_BAUD", "115200")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	endpoint, ok := cfg.MAVLink.Endpoint.(gomavlib.EndpointSerial)
	if !ok {
		t.Fatalf("expected EndpointSerial, got %T", cfg.MAVLink.Endpoint)
	}
	if endpoint.Device != "/dev/ttyUSB0" {
		t.Errorf("expected device '/dev/ttyUSB0', got '%s'", endpoint.Device)
	}
	if endpoint.Baud != 115200 {
		t.Errorf("expected baud 115200, got %d", endpoint.Baud)
	}
}

func TestLoad_MAVLinkSerial_MissingDevice(t *testing.T) {
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "serial")
	t.Setenv("FLIGHTPATH_MAVLINK_SERIAL_DEVICE", "")
	t.Setenv("FLIGHTPATH_MAVLINK_SERIAL_BAUD", "115200")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when serial device is missing, got nil")
	}
}

func TestLoad_MAVLinkSerial_MissingBaud(t *testing.T) {
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "serial")
	t.Setenv("FLIGHTPATH_MAVLINK_SERIAL_DEVICE", "/dev/ttyUSB0")
	t.Setenv("FLIGHTPATH_MAVLINK_SERIAL_BAUD", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when serial baud is missing, got nil")
	}
}

func TestLoad_MAVLinkSerial_InvalidBaud(t *testing.T) {
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "serial")
	t.Setenv("FLIGHTPATH_MAVLINK_SERIAL_DEVICE", "/dev/ttyUSB0")
	t.Setenv("FLIGHTPATH_MAVLINK_SERIAL_BAUD", "not-a-number")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when serial baud is invalid, got nil")
	}
}

func TestLoad_MAVLinkSerial_NegativeBaud(t *testing.T) {
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "serial")
	t.Setenv("FLIGHTPATH_MAVLINK_SERIAL_DEVICE", "/dev/ttyUSB0")
	t.Setenv("FLIGHTPATH_MAVLINK_SERIAL_BAUD", "-1")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when serial baud is negative, got nil")
	}
}

func TestLoad_MAVLinkUDP_MissingAddress(t *testing.T) {
	t.Run("udp-server", func(t *testing.T) {
		t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "udp-server")
		t.Setenv("FLIGHTPATH_MAVLINK_UDP_ADDRESS", "")

		_, err := Load()
		if err == nil {
			t.Fatal("expected error when UDP address is missing for udp-server, got nil")
		}
	})

	t.Run("udp-client", func(t *testing.T) {
		t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "udp-client")
		t.Setenv("FLIGHTPATH_MAVLINK_UDP_ADDRESS", "")

		_, err := Load()
		if err == nil {
			t.Fatal("expected error when UDP address is missing for udp-client, got nil")
		}
	})
}

func TestLoad_MAVLinkTCP_MissingAddress(t *testing.T) {
	t.Run("tcp-server", func(t *testing.T) {
		t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "tcp-server")
		t.Setenv("FLIGHTPATH_MAVLINK_TCP_ADDRESS", "")

		_, err := Load()
		if err == nil {
			t.Fatal("expected error when TCP address is missing for tcp-server, got nil")
		}
	})

	t.Run("tcp-client", func(t *testing.T) {
		t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "tcp-client")
		t.Setenv("FLIGHTPATH_MAVLINK_TCP_ADDRESS", "")

		_, err := Load()
		if err == nil {
			t.Fatal("expected error when TCP address is missing for tcp-client, got nil")
		}
	})
}

func TestLoad_MAVLinkUDPServer(t *testing.T) {
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "udp-server")
	t.Setenv("FLIGHTPATH_MAVLINK_UDP_ADDRESS", "0.0.0.0:14551")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	endpoint, ok := cfg.MAVLink.Endpoint.(gomavlib.EndpointUDPServer)
	if !ok {
		t.Fatalf("expected EndpointUDPServer, got %T", cfg.MAVLink.Endpoint)
	}
	if endpoint.Address != "0.0.0.0:14551" {
		t.Errorf("expected address '0.0.0.0:14551', got '%s'", endpoint.Address)
	}
}

func TestLoad_MAVLinkUDPClient(t *testing.T) {
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "udp-client")
	t.Setenv("FLIGHTPATH_MAVLINK_UDP_ADDRESS", "127.0.0.1:14550")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	endpoint, ok := cfg.MAVLink.Endpoint.(gomavlib.EndpointUDPClient)
	if !ok {
		t.Fatalf("expected EndpointUDPClient, got %T", cfg.MAVLink.Endpoint)
	}
	if endpoint.Address != "127.0.0.1:14550" {
		t.Errorf("expected address '127.0.0.1:14550', got '%s'", endpoint.Address)
	}
}

func TestLoad_MAVLinkTCPServer(t *testing.T) {
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "tcp-server")
	t.Setenv("FLIGHTPATH_MAVLINK_TCP_ADDRESS", "0.0.0.0:5760")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	endpoint, ok := cfg.MAVLink.Endpoint.(gomavlib.EndpointTCPServer)
	if !ok {
		t.Fatalf("expected EndpointTCPServer, got %T", cfg.MAVLink.Endpoint)
	}
	if endpoint.Address != "0.0.0.0:5760" {
		t.Errorf("expected address '0.0.0.0:5760', got '%s'", endpoint.Address)
	}
}

func TestLoad_MAVLinkTCPClient(t *testing.T) {
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "tcp-client")
	t.Setenv("FLIGHTPATH_MAVLINK_TCP_ADDRESS", "192.168.1.100:5760")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	endpoint, ok := cfg.MAVLink.Endpoint.(gomavlib.EndpointTCPClient)
	if !ok {
		t.Fatalf("expected EndpointTCPClient, got %T", cfg.MAVLink.Endpoint)
	}
	if endpoint.Address != "192.168.1.100:5760" {
		t.Errorf("expected address '192.168.1.100:5760', got '%s'", endpoint.Address)
	}
}

func TestLoad_UnknownEndpointType(t *testing.T) {
	t.Setenv("FLIGHTPATH_MAVLINK_ENDPOINT_TYPE", "unknown-type")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for unknown endpoint type, got nil")
	}
}
