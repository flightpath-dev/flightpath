package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

func main() {
	// Load configuration from environment variables (with sensible defaults)
	cfg := Load()

	// Create connection service client
	connectionService := createClient(cfg.ServerURL)

	// Setup graceful shutdown on Ctrl+C
	ctx := handleShutdown()

	// Stream heartbeats
	streamHeartbeats(ctx, connectionService, cfg.ServerURL)
}

// createClient creates the HTTP client and connection service client
func createClient(serverURL string) flightpathconnect.ConnectionServiceClient {
	// Create a single HTTP client to share across all service clients
	// This client uses the default transport which provides connection pooling
	httpClient := &http.Client{}

	// Create connection service client to communicate with the gRPC server
	return flightpathconnect.NewConnectionServiceClient(
		httpClient,
		serverURL,
		connect.WithProtoJSON(), // Use JSON codec for readability
	)
}

// handleShutdown handles Ctrl+C gracefully by canceling the context
func handleShutdown() context.Context {
	// Create a cancellable context – cancel() stops operations
	ctx, cancel := context.WithCancel(context.Background())

	// Handle Ctrl+C gracefully: when user presses Ctrl+C, cancel the context to stop the stream
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nStopping...")
		cancel() // Cancel the context, which signals StreamHeartbeats to stop
	}()

	return ctx
}

// streamHeartbeats connects to the server and streams heartbeat messages
func streamHeartbeats(ctx context.Context, connectionService flightpathconnect.ConnectionServiceClient, serverURL string) {
	fmt.Printf("Connecting to StreamHeartbeats endpoint: %s\n", serverURL)
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println("")

	// Create StreamHeartbeats request
	req := connect.NewRequest(&flightpath.StreamHeartbeatsRequest{})

	// Call StreamHeartbeats to start the stream, pass ctx for cancellation when user presses Ctrl+C
	stream, err := connectionService.StreamHeartbeats(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling StreamHeartbeats: %v\n", err)
		os.Exit(1)
	}

	// Receive messages from the stream (stream.Receive() is a blocking call)
	for stream.Receive() {
		// Get the message from the stream
		msg := stream.Msg()
		// Convert the timestamp to a human-readable format
		timestamp := time.Unix(0, msg.TimestampMs*int64(time.Millisecond))
		fmt.Printf("Received heartbeat: timestamp = %d ms (%s)\n",
			msg.TimestampMs,
			timestamp.Format("2006-01-02 15:04:05"),
		)
	}

	// Receive loop exited, check if there was an error from the stream
	if err := stream.Err(); err != nil {
		// Check if the error is due to context cancellation (user pressed Ctrl+C)
		if err == context.Canceled {
			fmt.Println("Stream canceled by user")
			return
		}
		fmt.Fprintf(os.Stderr, "Stream error: %v\n", err)
		os.Exit(1)
	}
}

// ----------------------------------------------------------------------------
// Configuration handling
// ----------------------------------------------------------------------------

// Config holds configuration for this example.
// This follows the 12-factor app methodology: configuration via environment
// variables with sensible defaults for local development.
type Config struct {
	ServerURL string // Server URL to connect to
}

// Default returns a Config with sensible defaults for local development.
func Default() *Config {
	return &Config{
		ServerURL: "http://localhost:8080",
	}
}

// Load loads configuration from environment variables, falling back to defaults.
// Environment Variables:
//   - FLIGHTPATH_SERVER_URL: Server URL (default: "http://localhost:8080")
func Load() *Config {
	cfg := Default()

	if serverURL := os.Getenv("FLIGHTPATH_SERVER_URL"); serverURL != "" {
		cfg.ServerURL = serverURL
	}

	return cfg
}
