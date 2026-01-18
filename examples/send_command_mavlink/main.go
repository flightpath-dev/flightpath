package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/internal/config"
	"github.com/flightpath-dev/flightpath/internal/mavlink/message_converters"
)

// ------------------------------------------------------------------------------------------------
// Send Command Example using MAVLink directly
// ------------------------------------------------------------------------------------------------
// This example demonstrates how to send MAVLink commands directly using gomavlib.
//
// It performs a complete flight sequence:
//   1. Initialize MAVLink node
//   2. Wait for heartbeat to discover channel and system/component IDs
//   3. Wait 10 seconds for drone to recognize that GCS is connected
//   4. Arm the drone (MAV_CMD_COMPONENT_ARM_DISARM)
//   5. Send takeoff command (MAV_CMD_DO_SET_MODE to AUTO/TAKEOFF)
//   6. Wait for ACTIVE + AUTO/TAKEOFF
//   7. Wait for ACTIVE + AUTO/LOITER
//   8. Wait 10 seconds in loiter
//   9. Send RTL command (MAV_CMD_DO_SET_MODE to AUTO/RTL)
//  10. Wait for ACTIVE + AUTO/RTL
//  11. Wait for STANDBY + AUTO/LOITER (landed and disarmed)
//  12. Exit
//
// Configuration is loaded from environment variables with sensible defaults:
//   - Default: UDP server on port 14550 (standard PX4 SITL port)
//   - See config.Load() function for all available environment variables
//
// To run this example:
//  1. Start a PX4 SITL (see docs/px4-sitl-setup.md)
//
//  2. Run this example using the default configuration (MAVLink running as a UDP server on port 14550)
//     go run examples/send_command_mavlink/main.go
//
//  3. Or configure a serial connection via environment variables:
//     export FLIGHTPATH_MAVLINK_ENDPOINT_TYPE=serial
//     export FLIGHTPATH_MAVLINK_SERIAL_DEVICE=/dev/cu.usbserial-D30JAXGS
//     export FLIGHTPATH_MAVLINK_SERIAL_BAUD=57600
//
//     go run examples/send_command_mavlink/main.go
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
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Step 1: Initialize MAVLink node
	// Create a node which acts as a GCS, communicating with the configured endpoint.
	// We use system ID 254 to coexist with QGroundControl (which uses 255).
	log.Println("Initializing MAVLink node...")
	node := &gomavlib.Node{
		Endpoints:   []gomavlib.EndpointConf{cfg.MAVLink.Endpoint},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 254,
	}
	err = node.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize node: %v", err)
	}
	defer node.Close()
	log.Println("Node initialized successfully")

	// Step 2: Wait for heartbeat to discover channel and system/component IDs
	heartbeatChan, heartbeatSystemID, heartbeatComponentID := detectChannel(node)

	// Step 3: Wait 10 seconds for drone to recognize that GCS is connected
	log.Println("Waiting 10 seconds for drone to recognize GCS connection...")
	sleepWhileListening(node, 10*time.Second)
	log.Println("GCS connection established")

	// Step 4: Arm the drone
	log.Println("Arming the drone...")
	err = sendArmCommand(node, heartbeatChan, heartbeatSystemID, heartbeatComponentID)
	if err != nil {
		log.Fatalf("Arm command failed: %v", err)
	}

	// Step 5: Send takeoff command (MAV_CMD_DO_SET_MODE to AUTO/TAKEOFF)
	log.Println("Sending takeoff command...")
	err = sendTakeoffCommand(node, heartbeatChan, heartbeatSystemID, heartbeatComponentID)
	if err != nil {
		log.Fatalf("Failed to set mode to AUTO/TAKEOFF: %v", err)
	}

	// Step 6: Wait for ACTIVE + AUTO/TAKEOFF
	log.Println("Waiting for ACTIVE + AUTO/TAKEOFF...")
	_, err = waitForHeartbeatCondition(
		node,
		heartbeatSystemID,
		heartbeatComponentID,
		flightpath.MavState_MAV_STATE_ACTIVE,
		flightpath.MainMode_MAIN_MODE_AUTO,
		flightpath.SubMode_SUB_MODE_AUTO_TAKEOFF,
		HeartbeatTimeout,
	)
	if err != nil {
		log.Fatalf("Failed to reach ACTIVE + AUTO/TAKEOFF: %v", err)
	}

	// Step 7: Wait for ACTIVE + AUTO/LOITER
	log.Println("Waiting for ACTIVE + AUTO/LOITER...")
	_, err = waitForHeartbeatCondition(
		node,
		heartbeatSystemID,
		heartbeatComponentID,
		flightpath.MavState_MAV_STATE_ACTIVE,
		flightpath.MainMode_MAIN_MODE_AUTO,
		flightpath.SubMode_SUB_MODE_AUTO_LOITER,
		HeartbeatTimeout,
	)
	if err != nil {
		log.Fatalf("Failed to reach ACTIVE + AUTO/LOITER: %v", err)
	}

	// Step 8: Wait for 10 seconds
	log.Println("Waiting 10 seconds in loiter mode...")
	sleepWhileListening(node, 10*time.Second)
	log.Println("Finished waiting in loiter mode")

	// Step 9: Send RTL command (MAV_CMD_DO_SET_MODE to AUTO/RTL)
	log.Println("Sending RTL command...")
	err = sendRTLCommand(node, heartbeatChan, heartbeatSystemID, heartbeatComponentID)
	if err != nil {
		log.Fatalf("Failed to set mode to AUTO/RTL: %v", err)
	}

	// Step 10: Wait for ACTIVE + AUTO/RTL
	log.Println("Waiting for ACTIVE + AUTO/RTL...")
	_, err = waitForHeartbeatCondition(
		node,
		heartbeatSystemID,
		heartbeatComponentID,
		flightpath.MavState_MAV_STATE_ACTIVE,
		flightpath.MainMode_MAIN_MODE_AUTO,
		flightpath.SubMode_SUB_MODE_AUTO_RTL,
		HeartbeatTimeout,
	)
	if err != nil {
		log.Fatalf("Failed to reach ACTIVE + AUTO/RTL: %v", err)
	}

	// Step 11: Wait for STANDBY + AUTO/LOITER (landed and disarmed)
	log.Println("Waiting for STANDBY + AUTO/LOITER (landed and disarmed)...")
	_, err = waitForHeartbeatCondition(
		node,
		heartbeatSystemID,
		heartbeatComponentID,
		flightpath.MavState_MAV_STATE_STANDBY,
		flightpath.MainMode_MAIN_MODE_AUTO,
		flightpath.SubMode_SUB_MODE_AUTO_LOITER,
		HeartbeatTimeout,
	)
	if err != nil {
		log.Fatalf("Failed to reach STANDBY + AUTO/LOITER (landed): %v", err)
	}

	// Step 12: Exit
	log.Println("Example completed successfully")
}

// detectChannel
// Detects the channel by waiting for a heartbeat message and returns the channel, system ID, and component ID.
func detectChannel(node *gomavlib.Node) (*gomavlib.Channel, uint8, uint8) {
	log.Println("Detecting channel...")
	for {
		evt := <-node.Events()

		if evt, ok := evt.(*gomavlib.EventFrame); ok {
			if _, ok := evt.Message().(*common.MessageHeartbeat); ok {
				channel := evt.Channel
				systemID := evt.SystemID()
				componentID := evt.ComponentID()
				log.Printf("Received heartbeat from channel %v on system %d, component %d\n", channel, systemID, componentID)
				return channel, systemID, componentID
			}
		}
	}
}

// waitForHeartbeatCondition
// Waits for a heartbeat message that matches the specified system status and custom mode conditions.
// Returns the heartbeat message if found, or an error if timeout occurs.
func waitForHeartbeatCondition(
	node *gomavlib.Node,
	targetSystemID, targetComponentID uint8,
	expectedSystemStatus flightpath.MavState,
	expectedMainMode flightpath.MainMode,
	expectedSubMode flightpath.SubMode,
	timeout time.Duration,
) (*common.MessageHeartbeat, error) {
	t := time.NewTimer(timeout)
	defer t.Stop()

	log.Printf("Wait heartbeat: system_status=%v, main_mode=%v, sub_mode=%v\n",
		expectedSystemStatus, expectedMainMode, expectedSubMode)

	for {
		select {
		case evt := <-node.Events():
			if evt, ok := evt.(*gomavlib.EventFrame); ok {
				if heartbeat, ok := evt.Message().(*common.MessageHeartbeat); ok {
					if evt.SystemID() == targetSystemID && evt.ComponentID() == targetComponentID {
						// Convert to protobuf to check conditions
						pbHeartbeat := message_converters.HeartbeatToProtobuf(heartbeat)

						if pbHeartbeat.SystemStatus == expectedSystemStatus &&
							pbHeartbeat.CustomMode != nil &&
							pbHeartbeat.CustomMode.MainMode == expectedMainMode &&
							pbHeartbeat.CustomMode.SubMode == expectedSubMode {
							log.Printf("Got heartbeat:  system_status=%v, main_mode=%v, sub_mode=%v\n",
								pbHeartbeat.SystemStatus, pbHeartbeat.CustomMode.MainMode, pbHeartbeat.CustomMode.SubMode)
							return heartbeat, nil
						}
					}
				}
			}

		case <-t.C:
			return nil, fmt.Errorf("timeout waiting for heartbeat condition: system_status=%v, main_mode=%v, sub_mode=%v",
				expectedSystemStatus, expectedMainMode, expectedSubMode)
		}
	}
}

// writeAndWaitCommandLong
// Sends a COMMAND_LONG message and waits for COMMAND_ACK response.
func writeAndWaitCommandLong(
	node *gomavlib.Node,
	channel *gomavlib.Channel,
	cmd *common.MessageCommandLong,
	timeout time.Duration,
) error {
	log.Printf("Sending command %d to system %d, component %d...", cmd.Command, cmd.TargetSystem, cmd.TargetComponent)
	err := node.WriteMessageTo(channel, cmd)
	if err != nil {
		return fmt.Errorf("failed to write command: %w", err)
	}

	t := time.NewTimer(timeout)
	defer t.Stop()

	for {
		select {
		case evt := <-node.Events():
			if evt, ok := evt.(*gomavlib.EventFrame); ok {
				if ack, ok2 := evt.Message().(*common.MessageCommandAck); ok2 {
					if ack.Command == cmd.Command &&
						evt.SystemID() == cmd.TargetSystem &&
						evt.ComponentID() == cmd.TargetComponent {
						switch {
						case ack.Result == common.MAV_RESULT_IN_PROGRESS:
							log.Printf("Command progress: %d%%\n", ack.Progress)

						case ack.Result != common.MAV_RESULT_ACCEPTED:
							return fmt.Errorf("command failed with result %v", ack.Result)

						default:
							log.Printf("Command succeeded (result: %v)\n", ack.Result)
							return nil
						}
					}
				}
			}

		case <-t.C:
			return fmt.Errorf("command timed out after %v", timeout)
		}
	}
}

// writeAndWaitCommandInt
// Sends a COMMAND_INT message and waits for COMMAND_ACK response.
func writeAndWaitCommandInt(
	node *gomavlib.Node,
	channel *gomavlib.Channel,
	cmd *common.MessageCommandInt,
	timeout time.Duration,
) error {
	log.Printf("Sending command %d to system %d, component %d...", cmd.Command, cmd.TargetSystem, cmd.TargetComponent)
	err := node.WriteMessageTo(channel, cmd)
	if err != nil {
		return fmt.Errorf("failed to write command: %w", err)
	}

	t := time.NewTimer(timeout)
	defer t.Stop()

	for {
		select {
		case evt := <-node.Events():
			if evt, ok := evt.(*gomavlib.EventFrame); ok {
				if ack, ok2 := evt.Message().(*common.MessageCommandAck); ok2 {
					if ack.Command == cmd.Command &&
						evt.SystemID() == cmd.TargetSystem &&
						evt.ComponentID() == cmd.TargetComponent {
						switch {
						case ack.Result == common.MAV_RESULT_IN_PROGRESS:
							log.Printf("Command progress: %d%%\n", ack.Progress)

						case ack.Result != common.MAV_RESULT_ACCEPTED:
							return fmt.Errorf("command failed with result %v", ack.Result)

						default:
							log.Printf("Command succeeded (result: %v)\n", ack.Result)
							return nil
						}
					}
				}
			}

		case <-t.C:
			return fmt.Errorf("command timed out after %v", timeout)
		}
	}
}

// sleepWhileListening
// Sleeps for the specified duration while continuing to process node events.
func sleepWhileListening(node *gomavlib.Node, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()

	for {
		select {
		case <-node.Events():
		case <-t.C:
			return
		}
	}
}

// sendArmCommand
// Sends an arm command to the drone using MAV_CMD_COMPONENT_ARM_DISARM.
func sendArmCommand(
	node *gomavlib.Node,
	channel *gomavlib.Channel,
	targetSystemID, targetComponentID uint8,
) error {
	return writeAndWaitCommandLong(node,
		channel,
		&common.MessageCommandLong{
			TargetSystem:    targetSystemID,
			TargetComponent: targetComponentID,
			Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
			Param1:          1, // 1 to arm, 0 to disarm
			Param2:          0, // 0 = normal arming (not force)
			Param3:          0, // Unused
			Param4:          0, // Unused
			Param5:          0, // Unused
			Param6:          0, // Unused
			Param7:          0, // Unused
		},
		5*time.Second,
	)
}

// requestGlobalPositionInt
// Requests GLOBAL_POSITION_INT message using MAV_CMD_REQUEST_MESSAGE and waits for the response.
func requestGlobalPositionInt(
	node *gomavlib.Node,
	channel *gomavlib.Channel,
	targetSystemID, targetComponentID uint8,
	timeout time.Duration,
) (*common.MessageGlobalPositionInt, error) {
	// Request GLOBAL_POSITION_INT (message ID 33) using MAV_CMD_REQUEST_MESSAGE
	log.Println("Requesting GLOBAL_POSITION_INT message...")
	err := node.WriteMessageTo(channel, &common.MessageCommandLong{
		TargetSystem:    targetSystemID,
		TargetComponent: targetComponentID,
		Command:         common.MAV_CMD_REQUEST_MESSAGE,
		Param1:          33, // Message ID for GLOBAL_POSITION_INT
		Param2:          0,
		Param3:          0,
		Param4:          0,
		Param5:          0,
		Param6:          0,
		Param7:          0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to request GLOBAL_POSITION_INT: %w", err)
	}

	t := time.NewTimer(timeout)
	defer t.Stop()

	for {
		select {
		case evt := <-node.Events():
			if evt, ok := evt.(*gomavlib.EventFrame); ok {
				if pos, ok2 := evt.Message().(*common.MessageGlobalPositionInt); ok2 {
					if evt.SystemID() == targetSystemID && evt.ComponentID() == targetComponentID {
						log.Printf("Received GLOBAL_POSITION_INT: lat=%d, lon=%d, alt=%d mm\n",
							pos.Lat, pos.Lon, pos.Alt)
						return pos, nil
					}
				}
			}

		case <-t.C:
			return nil, fmt.Errorf("timeout waiting for GLOBAL_POSITION_INT after %v", timeout)
		}
	}
}

// sendTakeoffCommand
// Sends a takeoff command using MAV_CMD_NAV_TAKEOFF with COMMAND_INT for high precision lat/lon.
// First requests the current position, then calculates takeoff altitude as current MSL altitude + 10 meters.
// Uses MAV_FRAME_GLOBAL_INT frame with absolute MSL altitude.
func sendTakeoffCommand(
	node *gomavlib.Node,
	channel *gomavlib.Channel,
	targetSystemID, targetComponentID uint8,
) error {
	// Step 1: Request current global position
	pos, err := requestGlobalPositionInt(node, channel, targetSystemID, targetComponentID, CommandTimeout)
	if err != nil {
		return fmt.Errorf("failed to get current position: %w", err)
	}

	// Step 2: Extract position data
	// GLOBAL_POSITION_INT uses:
	// - lat/lon: int32 in degrees * 1E7
	// - alt: int32 in mm (MSL)
	latitude := pos.Lat     // Already in degrees * 1E7
	longitude := pos.Lon    // Already in degrees * 1E7
	currentAltMm := pos.Alt // MSL altitude in mm

	// Step 3: Calculate takeoff altitude (current MSL + 10 meters)
	// 10 meters = 10000 mm
	takeoffAltMm := currentAltMm + 10000
	// Convert to meters for COMMAND_INT Z parameter (float)
	takeoffAltM := float32(takeoffAltMm) / 1000.0

	log.Printf("Current position: lat=%.7f, lon=%.7f, alt=%.2f m MSL\n",
		float32(latitude)/1e7, float32(longitude)/1e7, float32(currentAltMm)/1000.0)
	log.Printf("Takeoff altitude: %.2f m MSL\n", takeoffAltM)

	// Step 4: Send takeoff command with MAV_FRAME_GLOBAL_INT
	return writeAndWaitCommandInt(node,
		channel,
		&common.MessageCommandInt{
			TargetSystem:    targetSystemID,
			TargetComponent: targetComponentID,
			Frame:           common.MAV_FRAME_GLOBAL_INT,
			Command:         common.MAV_CMD_NAV_TAKEOFF,
			Current:         0,           // Not used
			Autocontinue:    0,           // Not used
			Param1:          -1,          // Minimum pitch (-1 = undefined, use default)
			Param2:          0,           // Empty
			Param3:          0,           // Empty
			Param4:          0,           // Yaw angle (0 = undefined, use current heading)
			X:               latitude,    // int32: latitude in degrees * 1E7
			Y:               longitude,   // int32: longitude in degrees * 1E7
			Z:               takeoffAltM, // float: altitude in meters (MSL, absolute)
		},
		CommandTimeout,
	)
}

// sendRTLCommand
// Sends a return-to-launch (RTL) command using MAV_CMD_NAV_RETURN_TO_LAUNCH.
func sendRTLCommand(
	node *gomavlib.Node,
	channel *gomavlib.Channel,
	targetSystemID, targetComponentID uint8,
) error {
	return writeAndWaitCommandLong(node,
		channel,
		&common.MessageCommandLong{
			TargetSystem:    targetSystemID,
			TargetComponent: targetComponentID,
			Command:         common.MAV_CMD_NAV_RETURN_TO_LAUNCH,
			Param1:          0, // Unused
			Param2:          0, // Unused
			Param3:          0, // Unused
			Param4:          0, // Unused
			Param5:          0, // Unused
			Param6:          0, // Unused
			Param7:          0, // Unused
		},
		CommandTimeout,
	)
}
