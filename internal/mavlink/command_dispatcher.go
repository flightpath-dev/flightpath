package mavlink

import (
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// ------------------------------------------------------------------------------------------------
// MAVLinkCommandDispatcher
// ------------------------------------------------------------------------------------------------
// Handles sending MAVLink commands to the drone.
// Encapsulates MAVLink protocol details so services don't need to deal with them directly.
//
//   - The main app creates it using `NewMAVLinkCommandDispatcher()`, passing it an initialized node.
//   - Commands are sent via `node.WriteMessageAll()` which sends to all active channels.
//   - This is appropriate for single-drone setups. If multi-drone support is needed later,
//     channel tracking can be added to use `WriteMessageTo(channel, cmd)` for precision.
//
// ------------------------------------------------------------------------------------------------
type MAVLinkCommandDispatcher struct {
	node *gomavlib.Node
}

// NewMAVLinkCommandDispatcher
// Creates a new command dispatcher using the provided node.
func NewMAVLinkCommandDispatcher(node *gomavlib.Node) *MAVLinkCommandDispatcher {
	return &MAVLinkCommandDispatcher{
		node: node,
	}
}

// SendCommandLong
// Sends a MAVLink COMMAND_LONG (76) message to the drone.
// Converts the protobuf request to a MAVLink MessageCommandLong and sends it.
// All parameters are floats.
func (d *MAVLinkCommandDispatcher) SendCommandLong(req *flightpath.SendCommandLongRequest) error {
	// Convert protobuf request to MAVLink MessageCommandLong
	msg := &common.MessageCommandLong{
		TargetSystem:    uint8(req.TargetSystemId),
		TargetComponent: uint8(req.TargetComponentId),
		Command:         common.MAV_CMD(req.Command),
		Confirmation:    0,
		Param1:          req.Param1,
		Param2:          req.Param2,
		Param3:          req.Param3,
		Param4:          req.Param4,
		Param5:          req.Param5,
		Param6:          req.Param6,
		Param7:          req.Param7,
	}

	// Send to all active channels
	return d.node.WriteMessageAll(msg)
}

// SendCommandInt
// Sends a MAVLink COMMAND_INT (75) message to the drone.
// Converts the protobuf request to a MAVLink MessageCommandInt and sends it.
// Parameters 5 and 6 (x, y) are int32 for higher precision (e.g., lat/lon scaled by 1E7).
func (d *MAVLinkCommandDispatcher) SendCommandInt(req *flightpath.SendCommandIntRequest) error {
	// Convert protobuf request to MAVLink MessageCommandInt
	msg := &common.MessageCommandInt{
		TargetSystem:    uint8(req.TargetSystemId),
		TargetComponent: uint8(req.TargetComponentId),
		Frame:           common.MAV_FRAME(req.Frame),
		Command:         common.MAV_CMD(req.Command),
		Current:         0, // Not used, set to 0
		Autocontinue:    0, // Not used, set to 0
		Param1:          req.Param1,
		Param2:          req.Param2,
		Param3:          req.Param3,
		Param4:          req.Param4,
		X:               req.X, // int32 for param5 (lat * 1E7 or local x * 1E4)
		Y:               req.Y, // int32 for param6 (lon * 1E7 or local y * 1E4)
		Z:               req.Z, // float for param7 (altitude in meters)
	}

	// Send to all active channels
	return d.node.WriteMessageAll(msg)
}
