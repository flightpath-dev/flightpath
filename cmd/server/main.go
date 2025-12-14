package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
	"github.com/flightpath-dev/flightpath/internal/config"
	"github.com/flightpath-dev/flightpath/internal/server"
	"github.com/flightpath-dev/flightpath/internal/services"
)

func main() {
	// Load configuration from environment variables (with sensible defaults)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create server
	srv := server.NewServer(cfg)

	// Register services
	registerServices(srv)

	// Setup graceful shutdown
	go handleShutdown(srv)

	// Start server
	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		srv.Logger().Fatalf("Server error: %v", err)
	}
}

// Register all services
func registerServices(srv *server.Server) {
	// Create shared service context
	ctx := &services.ServiceContext{
		Config: srv.Config(),
		Logger: srv.Logger(),
	}

	// ConnectionService
	connectionService := services.NewConnectionService(ctx)
	connectionPath, connectionHandler := flightpathconnect.NewConnectionServiceHandler(connectionService)
	srv.RegisterService(connectionPath, connectionHandler)
}

// handleShutdown handles graceful shutdown on interrupt signals
func handleShutdown(srv *server.Server) {
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

	srv.Logger().Println("✅ Cleanup complete")
	os.Exit(0)
}
