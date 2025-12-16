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
	go handleShutdown(srv, node, dispatcher, closeNode)

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

	// ConnectionService
	connectionService := services.NewConnectionService(ctx)
	connectionPath, connectionHandler := flightpathconnect.NewConnectionServiceHandler(connectionService)
	srv.RegisterService(connectionPath, connectionHandler)

	// TelemetryService
	telemetryService := services.NewTelemetryService(ctx)
	telemetryPath, telemetryHandler := flightpathconnect.NewTelemetryServiceHandler(telemetryService)
	srv.RegisterService(telemetryPath, telemetryHandler)
}

// handleShutdown handles graceful shutdown on interrupt signals
func handleShutdown(srv *server.Server, node *gomavlib.Node, dispatcher *services.MessageDispatcher, closeNode func()) {
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
