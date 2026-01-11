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

// SendCommand
// Sends a MAVLink command to the drone.
// Converts the protobuf request to a MAVLink MessageCommandLong and sends it.
func (d *MAVLinkCommandDispatcher) SendCommand(req *flightpath.SendCommandRequest) error {
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
