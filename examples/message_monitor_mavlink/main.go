package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/internal/config"
	mavcommon "github.com/flightpath-dev/flightpath/internal/mavlink/dialects/common"
	"github.com/flightpath-dev/flightpath/internal/mavlink/message_converters"
)

// ------------------------------------------------------------------------------------------------
// Message Monitor using MAVLink
// ------------------------------------------------------------------------------------------------
// This example shows how to act as a GCS by listening to the PX4 autopilot's broadcast messages.
// It uses gomavlib to connect to MAVLink endpoints (Serial, UDP, or TCP) and displays
// all received messages with counts.
//
// Configuration is loaded from environment variables with sensible defaults:
//   - Default: UDP server on port 14550 (standard PX4 SITL port)
//   - See config.Load() function for all available environment variables
//
// To run this example:
//  1. Start a PX4 SITL (see docs/px4-sitl-setup.md)
//
//  2. Run this example using the default configuration (MAVLink running as a UDP server on port 14550)
//     go run examples/message_monitor_mavlink/main.go
//
//  3. Or configure a serial connection via environment variables:
//     export FLIGHTPATH_MAVLINK_ENDPOINT_TYPE=serial
//     export FLIGHTPATH_MAVLINK_SERIAL_DEVICE=/dev/cu.usbserial-D30JAXGS
//     export FLIGHTPATH_MAVLINK_SERIAL_BAUD=57600
//
//     go run examples/message_monitor_mavlink/main.go
//
//  4. Or configure a UDP server connection via environment variables:
//     export FLIGHTPATH_MAVLINK_ENDPOINT_TYPE=udp-server
//     export FLIGHTPATH_MAVLINK_UDP_ADDRESS=0.0.0.0:14550
//
//     go run examples/message_monitor_mavlink/main.go
//
// Once started, you should see the PX4's broadcast message types and counts printed to the console.
// Additionally, the latest HEARTBEAT message is printed in detail.
// ------------------------------------------------------------------------------------------------

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}

	// Create a node which acts as a GCS, communicating with the configured endpoint.
	// We use system ID 254 to coexist with QGroundControl (which uses 255).
	node := &gomavlib.Node{
		Endpoints:   []gomavlib.EndpointConf{cfg.MAVLink.Endpoint},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 254,
	}
	err = node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// Data structures for tracking message counts and details
	var latestHeartbeat *flightpath.SubscribeHeartbeatResponse
	messageCounts := make(map[string]int)

	// Process incoming messages
	for evt := range node.Events() {
		// Process only frame events
		if eventFrame, ok := evt.(*gomavlib.EventFrame); ok {
			msg := eventFrame.Message()
			msgID := mavcommon.MavMessageId(msg.GetID())

			// Special handling for HEARTBEAT messages
			if heartbeat, ok := msg.(*common.MessageHeartbeat); ok {
				// Convert heartbeat to SubscribeHeartbeatResponse
				pbHeartbeat := message_converters.HeartbeatToProtobuf(heartbeat)
				latestHeartbeat = &flightpath.SubscribeHeartbeatResponse{
					TimestampMs: time.Now().UnixMilli(),
					SystemId:    uint32(eventFrame.SystemID()),
					ComponentId: uint32(eventFrame.ComponentID()),
					Heartbeat:   pbHeartbeat,
				}
				messageCounts["HEARTBEAT"]++
			} else {
				// For all other messages, just increment the count using the message ID as string
				messageCounts[msgID.String()]++
			}

			// Render dashboard after processing each message
			renderDashboard(latestHeartbeat, messageCounts)
		}
	}
}

// renderDashboard
// Renders a dashboard showing message counts and latest heartbeat information.
// Clears the screen and displays all information in a single update to minimize flicker.
func renderDashboard(latestHeartbeat *flightpath.SubscribeHeartbeatResponse, messageCounts map[string]int) {
	var buf strings.Builder

	// Clear screen and move cursor to top
	buf.WriteString("\033[2J\033[H")

	// Header
	buf.WriteString("=== MAVLink Message Monitor ===\n\n")

	// Latest HEARTBEAT message
	if latestHeartbeat != nil {
		buf.WriteString("Latest HEARTBEAT:\n")
		buf.WriteString("----------------\n")

		// Convert the timestamp to a human-readable format
		timestamp := time.Unix(0, latestHeartbeat.TimestampMs*int64(time.Millisecond))
		buf.WriteString(fmt.Sprintf("Timestamp: %s (%d ms)\n", timestamp.Format("2006-01-02 15:04:05.000"), latestHeartbeat.TimestampMs))

		// Print system and component IDs
		buf.WriteString(fmt.Sprintf("System ID: %d, Component ID: %d\n", latestHeartbeat.SystemId, latestHeartbeat.ComponentId))

		if latestHeartbeat.Heartbeat != nil {
			hb := latestHeartbeat.Heartbeat
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
