package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
	"github.com/flightpath-dev/flightpath/internal/config"
	"github.com/flightpath-dev/flightpath/internal/server"
	"github.com/flightpath-dev/flightpath/internal/services"
)

// ------------------------------------------------------------------------------------------------
// Flightpath Server
// ------------------------------------------------------------------------------------------------
// This is the main entry point for the Flightpath server.
// It loads configuration from environment variables, and connects to the drone on the configured
// MAVLink endpoint. It then starts the gRPC server, exposing the various services.
//
// See config.Load() function for all the available environment variables.
//
//  1. To run the server using the default configuration (MAVLink running as a UDP server on port 14550)
//     go run cmd/server/main.go
//
//  2. Or configure a serial connection via environment variables:
//     export FLIGHTPATH_MAVLINK_ENDPOINT_TYPE=serial
//     export FLIGHTPATH_MAVLINK_SERIAL_DEVICE=/dev/cu.usbserial-D30JAXGS
//     export FLIGHTPATH_MAVLINK_SERIAL_BAUD=57600
//
//     go run cmd/server/main.go
//
//  3. Or configure a UDP server connection via environment variables:
//     export FLIGHTPATH_MAVLINK_ENDPOINT_TYPE=udp-server
//     export FLIGHTPATH_MAVLINK_UDP_ADDRESS=0.0.0.0:14550
//
//     go run cmd/server/main.go
//
// ------------------------------------------------------------------------------------------------
func main() {
	// Load configuration from environment variables (with sensible defaults)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize MAVLink node
	// We use system ID 254 to coexist with QGroundControl (which uses 255).
	log.Println("📡 Initializing MAVLink node...")
	node := &gomavlib.Node{
		Endpoints:   []gomavlib.EndpointConf{cfg.MAVLink.Endpoint},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 254,
	}
	err = node.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize MAVLink node: %v", err)
	}
	log.Println("✅ MAVLink node initialized successfully")

	// Use sync.Once to ensure node is closed exactly once
	var nodeCloseOnce sync.Once
	closeNode := func() {
		nodeCloseOnce.Do(func() {
			log.Println("🔌 Closing MAVLink node...")
			if node != nil {
				node.Close()
			}
		})
	}

	// Ensure node is closed on any exit path
	defer closeNode()

	// Create message dispatcher and start it
	dispatcher := services.NewMessageDispatcher(node)
	dispatcher.Start()
	defer dispatcher.Stop()

	// Create server
	srv := server.NewServer(cfg)

	// Register services
	registerServices(srv, node, dispatcher)

	// Setup graceful shutdown
	go handleShutdown(srv, dispatcher, closeNode)

	// Start server
	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		srv.Logger().Fatalf("Server error: %v", err)
	}
}

// Register all services
func registerServices(srv *server.Server, node *gomavlib.Node, dispatcher *services.MessageDispatcher) {
	// Create shared service context
	ctx := &services.ServiceContext{
		Config:     srv.Config(),
		Logger:     srv.Logger(),
		Node:       node,
		Dispatcher: dispatcher,
	}

	// HeartbeatService
	heartbeatService := services.NewHeartbeatService(ctx)
	heartbeatPath, heartbeatHandler := flightpathconnect.NewHeartbeatServiceHandler(heartbeatService)
	srv.RegisterService(heartbeatPath, heartbeatHandler)

	// GpsRawIntService
	gpsRawIntService := services.NewGpsRawIntService(ctx)
	gpsRawIntPath, gpsRawIntHandler := flightpathconnect.NewGpsRawIntServiceHandler(gpsRawIntService)
	srv.RegisterService(gpsRawIntPath, gpsRawIntHandler)

	// SysStatusService
	sysStatusService := services.NewSysStatusService(ctx)
	sysStatusPath, sysStatusHandler := flightpathconnect.NewSysStatusServiceHandler(sysStatusService)
	srv.RegisterService(sysStatusPath, sysStatusHandler)

	// ExtendedSysStateService
	extendedSysStateService := services.NewExtendedSysStateService(ctx)
	extendedSysStatePath, extendedSysStateHandler := flightpathconnect.NewExtendedSysStateServiceHandler(extendedSysStateService)
	srv.RegisterService(extendedSysStatePath, extendedSysStateHandler)

	// StatusTextService
	statusTextService := services.NewStatusTextService(ctx)
	statusTextPath, statusTextHandler := flightpathconnect.NewStatusTextServiceHandler(statusTextService)
	srv.RegisterService(statusTextPath, statusTextHandler)

	// RadioStatusService
	radioStatusService := services.NewRadioStatusService(ctx)
	radioStatusPath, radioStatusHandler := flightpathconnect.NewRadioStatusServiceHandler(radioStatusService)
	srv.RegisterService(radioStatusPath, radioStatusHandler)

	// GlobalPositionIntService
	globalPositionIntService := services.NewGlobalPositionIntService(ctx)
	globalPositionIntPath, globalPositionIntHandler := flightpathconnect.NewGlobalPositionIntServiceHandler(globalPositionIntService)
	srv.RegisterService(globalPositionIntPath, globalPositionIntHandler)
}

// handleShutdown handles graceful shutdown on interrupt signals
func handleShutdown(srv *server.Server, dispatcher *services.MessageDispatcher, closeNode func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	srv.Logger().Println("🛑 Shutting down server gracefully...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		srv.Logger().Printf("Error during server shutdown: %v", err)
	}

	// Stop message dispatcher
	dispatcher.Stop()

	// Close MAVLink node (sync.Once ensures this is only called once)
	closeNode()

	srv.Logger().Println("✅ Cleanup complete")
	os.Exit(0)
}
