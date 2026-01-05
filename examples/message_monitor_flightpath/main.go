package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
	"github.com/flightpath-dev/flightpath/internal/config"
)

// ------------------------------------------------------------------------------------------------
// Message Monitor using Flightpath gRPC API
// ------------------------------------------------------------------------------------------------
// This example shows how to connect to the Flightpath gRPC server and stream heartbeat messages.
// It uses the Connect RPC client to communicate with the server and displays all received
// heartbeat messages with detailed information including system/component IDs, vehicle type,
// autopilot type, flight modes, and system status.
//
// Configuration is loaded from environment variables with sensible defaults:
//   - Default: http://localhost:8080 (standard Flightpath server address)
//   - See config.Load() function for all available environment variables
//
// To run this example:
//  1. Start a PX4 SITL (see docs/px4-sitl-setup.md)
//
//  2. Start the Flightpath server using the default configuration
//     (MAVLink running as a UDP server on port 14550 and gRPC running on http://localhost:8080)
//     go run cmd/server/main.go
//
//  3. Run this example using the default configuration (connecting to the gRPC server at http://localhost:8080)
//     go run examples/message_monitor_flightpath/main.go
//
// Once started, you should see PX4 heartbeat messages and message counts printed to the console.
// ------------------------------------------------------------------------------------------------

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}

	// Construct server URL from config
	serverURL := fmt.Sprintf("http://%s", cfg.ServerAddr())

	// Create service clients
	mavlinkService := createMAVLinkServiceClient(serverURL)

	// Setup graceful shutdown on Ctrl+C
	ctx := handleShutdown()

	// Data structures for tracking message counts and details
	var latestHeartbeat *flightpath.Heartbeat
	var latestGpsRawInt *flightpath.GpsRawInt
	messageCounts := make(map[string]int)
	var mu sync.Mutex

	// Subscribe to all MAVLink messages
	subscribeMessages(ctx, mavlinkService, serverURL, &latestHeartbeat, &latestGpsRawInt, messageCounts, &mu)
}

// createMAVLinkServiceClient creates the MAVLink service client
func createMAVLinkServiceClient(serverURL string) flightpathconnect.MAVLinkServiceClient {
	return flightpathconnect.NewMAVLinkServiceClient(
		&http.Client{},
		serverURL,
		connect.WithProtoJSON(),
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
		cancel() // Cancel the context, which signals SubscribeHeartbeat to stop
	}()

	return ctx
}

// subscribeMessages connects to the server and streams all MAVLink messages
func subscribeMessages(
	ctx context.Context,
	mavlinkService flightpathconnect.MAVLinkServiceClient,
	serverURL string,
	latestHeartbeat **flightpath.Heartbeat,
	latestGpsRawInt **flightpath.GpsRawInt,
	messageCounts map[string]int,
	mu *sync.Mutex,
) {
	fmt.Printf("Connecting to SubscribeMessages endpoint: %s\n", serverURL)

	// Create SubscribeMessages request (empty = all message types)
	req := connect.NewRequest(&flightpath.SubscribeMessagesRequest{})

	// Call SubscribeMessages to start the stream, pass ctx for cancellation when user presses Ctrl+C
	stream, err := mavlinkService.SubscribeMessages(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling SubscribeMessages: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Press Ctrl+C to stop")
	fmt.Println("")

	// Receive messages from the stream (stream.Receive() is a blocking call)
	for stream.Receive() {
		// Get the message from the stream
		msg := stream.Msg()

		// Process message based on type
		mu.Lock()
		switch msg.MessageType {
		case flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_HEARTBEAT:
			if heartbeat := msg.GetHeartbeat(); heartbeat != nil {
				messageCounts["HEARTBEAT"]++
				*latestHeartbeat = heartbeat
			}
		case flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_GPS_RAW_INT:
			if gps := msg.GetGpsRawInt(); gps != nil {
				messageCounts["GPS_RAW_INT"]++
				*latestGpsRawInt = gps
			}
		case flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_SYS_STATUS:
			messageCounts["SYS_STATUS"]++
		case flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_EXTENDED_SYS_STATE:
			messageCounts["EXTENDED_SYS_STATE"]++
		case flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_STATUSTEXT:
			messageCounts["STATUSTEXT"]++
		case flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_RADIO_STATUS:
			messageCounts["RADIO_STATUS"]++
		case flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_GLOBAL_POSITION_INT:
			messageCounts["GLOBAL_POSITION_INT"]++
		case flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_VFR_HUD:
			messageCounts["VFR_HUD"]++
		}
		mu.Unlock()

		// Render dashboard after processing each message
		mu.Lock()
		renderDashboard(*latestHeartbeat, *latestGpsRawInt, messageCounts)
		mu.Unlock()
	}

	// Receive loop exited, check if there was an error from the stream
	if err := stream.Err(); err != nil {
		// Check if the error is due to context cancellation (user pressed Ctrl+C)
		if err == context.Canceled {
			fmt.Println("\nStream canceled by user")
			return
		}
		fmt.Fprintf(os.Stderr, "Stream error: %v\n", err)
		os.Exit(1)
	}
}

// renderDashboard
// Renders a dashboard showing message counts and latest message information.
// Clears the screen and displays all information in a single update to minimize flicker.
func renderDashboard(latestHeartbeat *flightpath.Heartbeat, latestGpsRawInt *flightpath.GpsRawInt, messageCounts map[string]int) {
	var buf strings.Builder

	// Clear screen and move cursor to top
	buf.WriteString("\033[2J\033[H")

	// Header
	buf.WriteString("=== Flightpath Message Monitor ===\n\n")

	// Latest HEARTBEAT message
	if latestHeartbeat != nil {
		buf.WriteString("Latest HEARTBEAT:\n")
		buf.WriteString("----------------\n")

		hb := latestHeartbeat
		buf.WriteString(fmt.Sprintf("Vehicle Type: %s\n", hb.Type.String()))
		buf.WriteString(fmt.Sprintf("Autopilot: %s\n", hb.Autopilot.String()))
		buf.WriteString(fmt.Sprintf("System Status: %s\n", hb.SystemStatus.String()))
		buf.WriteString(fmt.Sprintf("MAVLink Version: %d\n", hb.MavlinkVersion))

		if hb.BaseMode != nil {
			bm := hb.BaseMode
			buf.WriteString(fmt.Sprintf(
				"Base Mode: custom_mode=%v, test=%v, auto=%v, guided=%v, stabilize=%v, hil=%v, manual=%v, safety=%v\n",
				bm.CustomModeEnabled, bm.TestEnabled, bm.AutoEnabled, bm.GuidedEnabled,
				bm.StabilizeEnabled, bm.HilEnabled, bm.ManualInputEnabled, bm.SafetyArmed))
		}

		if hb.CustomMode != nil {
			cm := hb.CustomMode
			buf.WriteString(fmt.Sprintf("Custom Mode: %s / %s\n", cm.MainMode.String(), cm.SubMode.String()))
		}

		buf.WriteString("\n")
	}

	// Latest GPS_RAW_INT message
	if latestGpsRawInt != nil {
		buf.WriteString("Latest GPS_RAW_INT:\n")
		buf.WriteString("------------------\n")

		gps := latestGpsRawInt
		buf.WriteString(fmt.Sprintf("Fix Type: %s\n", gps.FixType.String()))
		buf.WriteString(fmt.Sprintf("Latitude: %.7f° (raw: %d)\n", float64(gps.Lat)/1e7, gps.Lat))
		buf.WriteString(fmt.Sprintf("Longitude: %.7f° (raw: %d)\n", float64(gps.Lon)/1e7, gps.Lon))
		buf.WriteString(fmt.Sprintf("Altitude (MSL): %.3f m (raw: %d mm)\n", float64(gps.Alt)/1000, gps.Alt))
		if gps.AltEllipsoid != 0 {
			buf.WriteString(fmt.Sprintf("Altitude (Ellipsoid): %.3f m (raw: %d mm)\n", float64(gps.AltEllipsoid)/1000, gps.AltEllipsoid))
		}
		buf.WriteString(fmt.Sprintf("HDOP: %.2f (raw: %d)\n", float64(gps.Eph)/100, gps.Eph))
		buf.WriteString(fmt.Sprintf("VDOP: %.2f (raw: %d)\n", float64(gps.Epv)/100, gps.Epv))
		if gps.Vel != 65535 {
			buf.WriteString(fmt.Sprintf("Ground Speed: %.2f m/s (raw: %d cm/s)\n", float64(gps.Vel)/100, gps.Vel))
		}
		if gps.Cog != 65535 {
			buf.WriteString(fmt.Sprintf("Course over Ground: %.2f° (raw: %d)\n", float64(gps.Cog)/100, gps.Cog))
		}
		if gps.SatellitesVisible != 255 {
			buf.WriteString(fmt.Sprintf("Satellites Visible: %d\n", gps.SatellitesVisible))
		}
		if gps.HAcc != 0 {
			buf.WriteString(fmt.Sprintf("Horizontal Accuracy: %.3f m (raw: %d mm)\n", float64(gps.HAcc)/1000, gps.HAcc))
		}
		if gps.VAcc != 0 {
			buf.WriteString(fmt.Sprintf("Vertical Accuracy: %.3f m (raw: %d mm)\n", float64(gps.VAcc)/1000, gps.VAcc))
		}
		if gps.VelAcc != 0 {
			buf.WriteString(fmt.Sprintf("Speed Accuracy: %.3f m/s (raw: %d mm/s)\n", float64(gps.VelAcc)/1000, gps.VelAcc))
		}

		buf.WriteString("\n")
	}

	// Message counts table
	buf.WriteString("Message Counts:\n")
	buf.WriteString("---------------\n")

	// Sort message types by name for consistent display
	messageTypes := make([]string, 0, len(messageCounts))
	for msgType := range messageCounts {
		messageTypes = append(messageTypes, msgType)
	}
	sort.Strings(messageTypes)

	// Print message counts
	for _, msgType := range messageTypes {
		buf.WriteString(fmt.Sprintf("  %-30s %d\n", msgType, messageCounts[msgType]))
	}

	buf.WriteString("\n")

	// Write everything at once to minimize flicker
	fmt.Fprint(os.Stdout, buf.String())
}
