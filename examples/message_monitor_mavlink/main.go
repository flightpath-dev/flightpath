package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/internal/mavlink"
	mavcommon "github.com/flightpath-dev/flightpath/internal/mavlink/dialects/common"
)

// ------------------------------------------------------------------------------------------------
// Message Monitor using MAVLink
// ------------------------------------------------------------------------------------------------
// This example shows how to act as a GCS by listening to the PX4 autopilot's broadcast messages.
// It uses gomavlib to connect to UDP port 14550 where the PX4 broadcasts its messages using MAVLink.
//
// To run this example:
//  1. Start a PX4 SITL (see docs/px4-sitl-setup.md)
//  2. Run this example:
//     go run examples/message_monitor_mavlink/main.go
//
// You should see the PX4's broadcast message types and counts printed to the console.
// Additionally, the latest HEARTBEAT message is printed in detail.
// ------------------------------------------------------------------------------------------------

func main() {
	// Create a node which acts as a GCS, communicating with a UDP endpoint in server mode.
	// We use port 14550 because that's where the PX4 broadcasts its messages.
	// We use system ID 254 to coexist with QGroundControl (which uses 255).
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointUDPServer{Address: "0.0.0.0:14550"},
		},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 254,
	}
	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// Data structures for tracking message counts and details
	messageCounts := make(map[mavcommon.MavMessageId]int)
	var latestHeartbeat string

	// Process incoming messages
	for evt := range node.Events() {
		// Process only frame events
		if eventFrame, ok := evt.(*gomavlib.EventFrame); ok {
			msg := eventFrame.Message()
			msgID := mavcommon.MavMessageId(msg.GetID())

			// Special handling for HEARTBEAT messages
			if heartbeat, ok := msg.(*common.MessageHeartbeat); ok {
				latestHeartbeat = processHeartbeatMessage(heartbeat, eventFrame.SystemID(), eventFrame.ComponentID(), messageCounts)
			} else {
				// For all other messages, just increment the count using the message ID
				messageCounts[msgID]++
			}

			// Render dashboard after processing each message
			renderDashboard(messageCounts, latestHeartbeat)
		}
	}
}

// processHeartbeatMessage
// Processes a HEARTBEAT message by converting to map, marshaling to JSON, and formatting.
// Increments the message count and returns a formatted string on success, or empty string on error.
func processHeartbeatMessage(msg *common.MessageHeartbeat, systemID uint8, componentID uint8, messageCounts map[mavcommon.MavMessageId]int) string {
	// Converts the message to a map with decoded fields for better readability
	msgMap, err := mavlink.HeartbeatMessageToMap(msg)
	if err != nil {
		messageCounts[mavcommon.MavMessageIdHeartbeat]++
		return ""
	}

	// Pretty print message map as JSON
	msgJSON, err := json.MarshalIndent(msgMap, "", "  ")
	if err != nil {
		messageCounts[mavcommon.MavMessageIdHeartbeat]++
		return ""
	}

	// Format with system and component IDs
	formatted := fmt.Sprintf("System ID: %d, Component ID: %d\n%s", systemID, componentID, string(msgJSON))
	messageCounts[mavcommon.MavMessageIdHeartbeat]++
	return formatted
}

// renderDashboard
// Renders a dashboard showing message counts and latest heartbeat information.
// Clears the screen and displays all information in a single update to minimize flicker.
func renderDashboard(messageCounts map[mavcommon.MavMessageId]int, latestHeartbeat string) {
	var buf strings.Builder

	// Clear screen and move cursor to top
	buf.WriteString("\033[2J\033[H")

	// Header
	buf.WriteString("=== MAVLink Message Monitor ===\n\n")

	// Latest HEARTBEAT message
	if latestHeartbeat != "" {
		buf.WriteString("Latest HEARTBEAT:\n")
		buf.WriteString(latestHeartbeat)
		buf.WriteString("\n\n")
	}

	// Message counts table
	buf.WriteString("Message Counts:\n")
	buf.WriteString("---------------\n")

	// Sort message IDs by name for consistent display
	messageIDs := make([]mavcommon.MavMessageId, 0, len(messageCounts))
	for id := range messageCounts {
		messageIDs = append(messageIDs, id)
	}
	sort.Slice(messageIDs, func(i, j int) bool { return messageIDs[i].String() < messageIDs[j].String() })

	// Print message counts with IDs in parentheses
	for _, id := range messageIDs {
		displayName := fmt.Sprintf("%s (%d)", id.String(), uint32(id))
		buf.WriteString(fmt.Sprintf("  %-30s %d\n", displayName, messageCounts[id]))
	}

	buf.WriteString("\n")

	// Write everything at once to minimize flicker
	fmt.Fprint(os.Stdout, buf.String())
}
