package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
	"github.com/flightpath-dev/flightpath/internal/config"
)

// ------------------------------------------------------------------------------------------------
// Send Command Example using Flightpath gRPC API
// ------------------------------------------------------------------------------------------------
// This example demonstrates how to send MAVLink commands to a drone using the Flightpath gRPC API.
//
// It performs a complete flight sequence:
//   1. Connect to Flightpath gRPC server and discover system/component IDs from heartbeat
//   2. Arm the drone (MAV_CMD_COMPONENT_ARM_DISARM)
//   3. Send takeoff command (MAV_CMD_DO_SET_MODE to AUTO/TAKEOFF)
//   4. Wait for ACTIVE + AUTO/TAKEOFF
//   5. Wait for ACTIVE + AUTO/LOITER
//   6. Wait 10 seconds in loiter
//   7. Send RTL command (MAV_CMD_DO_SET_MODE to AUTO/RTL)
//   8. Wait for ACTIVE + AUTO/RTL
//   9. Wait for STANDBY + AUTO/LOITER (landed and disarmed)
//  10. Exit
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
//     go run examples/send_command_flightpath/main.go
// ------------------------------------------------------------------------------------------------

const (
	// Command timeout
	CommandTimeout = 10 * time.Second
	// Heartbeat wait timeout
	HeartbeatTimeout = 60 * time.Second
)

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

	// Channel to signal when heartbeat conditions are met
	heartbeatConditionMet := make(chan bool, 1)
	var mu sync.Mutex
	var latestHeartbeat *flightpath.Heartbeat
	var targetSystemID uint32
	var targetComponentID uint32

	// Channels to receive system_id and component_id from first heartbeat
	systemIDChan := make(chan uint32, 1)
	componentIDChan := make(chan uint32, 1)

	// Start goroutine to monitor heartbeat messages
	go monitorHeartbeats(ctx, mavlinkService, serverURL, &latestHeartbeat, heartbeatConditionMet, &mu, systemIDChan, componentIDChan)

	// Wait for first heartbeat to get system and component IDs
	fmt.Println("Waiting for first heartbeat to discover system and component IDs...")
	select {
	case targetSystemID = <-systemIDChan:
		targetComponentID = <-componentIDChan
		fmt.Printf("Discovered system ID: %d, component ID: %d\n", targetSystemID, targetComponentID)
	case <-ctx.Done():
		fmt.Println("Context canceled, exiting...")
		return
	case <-time.After(HeartbeatTimeout):
		fmt.Fprintf(os.Stderr, "Timeout waiting for first heartbeat\n")
		os.Exit(1)
	}

	// Step 2: Arm the drone
	fmt.Println("Arming the drone...")
	err = sendArmCommand(ctx, mavlinkService, targetSystemID, targetComponentID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending arm command: %v\n", err)
		os.Exit(1)
	}

	// Step 3: Send takeoff command (MAV_CMD_DO_SET_MODE to AUTO/TAKEOFF)
	fmt.Println("Sending takeoff command...")
	err = sendTakeoffCommand(ctx, mavlinkService, targetSystemID, targetComponentID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending takeoff command: %v\n", err)
		os.Exit(1)
	}

	// Step 4: Wait for ACTIVE + AUTO/TAKEOFF
	fmt.Println("Waiting for ACTIVE + AUTO/TAKEOFF...")
	if !waitForHeartbeatCondition(ctx, &latestHeartbeat, heartbeatConditionMet, &mu,
		flightpath.MavState_MAV_STATE_ACTIVE,
		flightpath.MainMode_MAIN_MODE_AUTO,
		flightpath.SubMode_SUB_MODE_AUTO_TAKEOFF,
		HeartbeatTimeout) {
		fmt.Fprintf(os.Stderr, "Failed to reach ACTIVE + AUTO/TAKEOFF\n")
		os.Exit(1)
	}

	// Step 5: Wait for ACTIVE + AUTO/LOITER
	fmt.Println("Waiting for ACTIVE + AUTO/LOITER...")
	if !waitForHeartbeatCondition(ctx, &latestHeartbeat, heartbeatConditionMet, &mu,
		flightpath.MavState_MAV_STATE_ACTIVE,
		flightpath.MainMode_MAIN_MODE_AUTO,
		flightpath.SubMode_SUB_MODE_AUTO_LOITER,
		HeartbeatTimeout) {
		fmt.Fprintf(os.Stderr, "Failed to reach ACTIVE + AUTO/LOITER\n")
		os.Exit(1)
	}

	// Step 6: Wait 10 seconds in loiter
	fmt.Println("Waiting 10 seconds in loiter mode...")
	time.Sleep(10 * time.Second)
	fmt.Println("Finished waiting in loiter mode")

	// Step 7: Send RTL command (MAV_CMD_DO_SET_MODE to AUTO/RTL)
	fmt.Println("Sending RTL command...")
	err = sendRTLCommand(ctx, mavlinkService, targetSystemID, targetComponentID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending RTL command: %v\n", err)
		os.Exit(1)
	}

	// Step 8: Wait for ACTIVE + AUTO/RTL
	fmt.Println("Waiting for ACTIVE + AUTO/RTL...")
	if !waitForHeartbeatCondition(ctx, &latestHeartbeat, heartbeatConditionMet, &mu,
		flightpath.MavState_MAV_STATE_ACTIVE,
		flightpath.MainMode_MAIN_MODE_AUTO,
		flightpath.SubMode_SUB_MODE_AUTO_RTL,
		HeartbeatTimeout) {
		fmt.Fprintf(os.Stderr, "Failed to reach ACTIVE + AUTO/RTL\n")
		os.Exit(1)
	}

	// Step 9: Wait for STANDBY + AUTO/LOITER (landed and disarmed)
	fmt.Println("Waiting for STANDBY + AUTO/LOITER (landed and disarmed)...")
	if !waitForHeartbeatCondition(ctx, &latestHeartbeat, heartbeatConditionMet, &mu,
		flightpath.MavState_MAV_STATE_STANDBY,
		flightpath.MainMode_MAIN_MODE_AUTO,
		flightpath.SubMode_SUB_MODE_AUTO_LOITER,
		HeartbeatTimeout) {
		fmt.Fprintf(os.Stderr, "Failed to reach STANDBY + AUTO/LOITER (landed)\n")
		os.Exit(1)
	}

	// Step 10: Exit
	fmt.Println("Example completed successfully")
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
		cancel() // Cancel the context, which signals operations to stop
	}()

	return ctx
}

// monitorHeartbeats
// Subscribes to heartbeat messages and monitors for various conditions.
// Signals on heartbeatConditionMet channel when any heartbeat condition is met.
// Also sends system_id and component_id from the first heartbeat message.
func monitorHeartbeats(
	ctx context.Context,
	mavlinkService flightpathconnect.MAVLinkServiceClient,
	serverURL string,
	latestHeartbeat **flightpath.Heartbeat,
	heartbeatConditionMet chan<- bool,
	mu *sync.Mutex,
	systemIDChan chan<- uint32,
	componentIDChan chan<- uint32,
) {
	fmt.Printf("Connecting to SubscribeMessages endpoint: %s\n", serverURL)

	// Create SubscribeMessages request - only subscribe to heartbeat messages
	req := connect.NewRequest(&flightpath.SubscribeMessagesRequest{
		MessageTypes: []flightpath.MavlinkMessageType{
			flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_HEARTBEAT,
		},
	})

	// Call SubscribeMessages to start the stream
	stream, err := mavlinkService.SubscribeMessages(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling SubscribeMessages: %v\n", err)
		os.Exit(1)
	}

	firstHeartbeat := true

	// Receive messages from the stream
	for stream.Receive() {
		// Get the message from the stream
		msg := stream.Msg()

		// Process heartbeat messages
		if msg.MessageType == flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_HEARTBEAT {
			if heartbeat := msg.GetHeartbeat(); heartbeat != nil {
				mu.Lock()
				*latestHeartbeat = heartbeat
				mu.Unlock()

				// Extract system_id and component_id from the first heartbeat
				if firstHeartbeat {
					select {
					case systemIDChan <- msg.SystemId:
					default:
					}
					select {
					case componentIDChan <- msg.ComponentId:
					default:
					}
					firstHeartbeat = false
				}

				// Signal that a new heartbeat was received (non-blocking)
				// The waitForHeartbeatCondition function will check if it matches the desired condition
				select {
				case heartbeatConditionMet <- true:
				default:
				}
			}
		}
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

// waitForHeartbeatCondition
// Waits for a heartbeat message that matches the specified system status and custom mode conditions.
// Returns true if the condition is met, false if timeout occurs.
func waitForHeartbeatCondition(
	ctx context.Context,
	latestHeartbeat **flightpath.Heartbeat,
	heartbeatConditionMet <-chan bool,
	mu *sync.Mutex,
	expectedSystemStatus flightpath.MavState,
	expectedMainMode flightpath.MainMode,
	expectedSubMode flightpath.SubMode,
	timeout time.Duration,
) bool {
	t := time.NewTimer(timeout)
	defer t.Stop()

	for {
		select {
		case <-heartbeatConditionMet:
			mu.Lock()
			heartbeat := *latestHeartbeat
			mu.Unlock()

			if heartbeat != nil &&
				heartbeat.SystemStatus == expectedSystemStatus &&
				heartbeat.CustomMode != nil &&
				heartbeat.CustomMode.MainMode == expectedMainMode &&
				heartbeat.CustomMode.SubMode == expectedSubMode {
				fmt.Printf("Got heartbeat: system_status=%v, main_mode=%v, sub_mode=%v\n",
					heartbeat.SystemStatus, heartbeat.CustomMode.MainMode, heartbeat.CustomMode.SubMode)
				return true
			}

		case <-t.C:
			return false

		case <-ctx.Done():
			return false
		}
	}
}

// sendArmCommand
// Sends an arm command to the drone.
func sendArmCommand(ctx context.Context, mavlinkService flightpathconnect.MAVLinkServiceClient, targetSystemID, targetComponentID uint32) error {
	req := connect.NewRequest(&flightpath.SendCommandRequest{
		TargetSystemId:    targetSystemID,
		TargetComponentId: targetComponentID,
		Command:           uint32(common.MAV_CMD_COMPONENT_ARM_DISARM),
		Param1:            1.0, // 1 to arm, 0 to disarm
		Param2:            0.0, // 0 = normal arming (not force)
		Param3:            0.0, // Unused
		Param4:            0.0, // Unused
		Param5:            0.0, // Unused
		Param6:            0.0, // Unused
		Param7:            0.0, // Unused
	})

	resp, err := mavlinkService.SendCommand(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to call SendCommand: %w", err)
	}

	if !resp.Msg.Success {
		return fmt.Errorf("command failed: %s", resp.Msg.ErrorMessage)
	}

	return nil
}

// sendTakeoffCommand
// Sends a takeoff command by setting mode to AUTO/TAKEOFF using MAV_CMD_DO_SET_MODE.
func sendTakeoffCommand(ctx context.Context, mavlinkService flightpathconnect.MAVLinkServiceClient, targetSystemID, targetComponentID uint32) error {
	req := connect.NewRequest(&flightpath.SendCommandRequest{
		TargetSystemId:    targetSystemID,
		TargetComponentId: targetComponentID,
		Command:           uint32(common.MAV_CMD_DO_SET_MODE),
		Param1:            129.0, // MAV_MODE_FLAG_SAFETY_ARMED (128) | MAV_MODE_FLAG_CUSTOM_MODE_ENABLED (1)
		Param2:            float32(flightpath.MainMode_MAIN_MODE_AUTO),
		Param3:            float32(flightpath.SubMode_SUB_MODE_AUTO_TAKEOFF),
		Param4:            0.0,
		Param5:            0.0,
		Param6:            0.0,
		Param7:            0.0,
	})

	resp, err := mavlinkService.SendCommand(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to call SendCommand: %w", err)
	}

	if !resp.Msg.Success {
		return fmt.Errorf("command failed: %s", resp.Msg.ErrorMessage)
	}

	return nil
}

// sendRTLCommand
// Sends a return-to-launch (RTL) command by setting mode to AUTO/RTL using MAV_CMD_DO_SET_MODE.
func sendRTLCommand(ctx context.Context, mavlinkService flightpathconnect.MAVLinkServiceClient, targetSystemID, targetComponentID uint32) error {
	req := connect.NewRequest(&flightpath.SendCommandRequest{
		TargetSystemId:    targetSystemID,
		TargetComponentId: targetComponentID,
		Command:           uint32(common.MAV_CMD_DO_SET_MODE),
		Param1:            129.0, // MAV_MODE_FLAG_SAFETY_ARMED (128) | MAV_MODE_FLAG_CUSTOM_MODE_ENABLED (1)
		Param2:            float32(flightpath.MainMode_MAIN_MODE_AUTO),
		Param3:            float32(flightpath.SubMode_SUB_MODE_AUTO_RTL),
		Param4:            0.0,
		Param5:            0.0,
		Param6:            0.0,
		Param7:            0.0,
	})

	resp, err := mavlinkService.SendCommand(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to call SendCommand: %w", err)
	}

	if !resp.Msg.Success {
		return fmt.Errorf("command failed: %s", resp.Msg.ErrorMessage)
	}

	return nil
}
