package message_converters

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// HeartbeatToProtobuf
// Converts a MAVLink HEARTBEAT message to a protobuf Heartbeat message.
func HeartbeatToProtobuf(msg *common.MessageHeartbeat) *flightpath.Heartbeat {
	return &flightpath.Heartbeat{
		Type:           MavTypeToProtobuf(msg.Type),
		Autopilot:      MavAutopilotToProtobuf(msg.Autopilot),
		BaseMode:       BaseModeToProtobuf(msg.BaseMode),
		CustomMode:     CustomModeToProtobuf(msg.CustomMode, msg.Autopilot),
		SystemStatus:   MavStateToProtobuf(msg.SystemStatus),
		MavlinkVersion: uint32(msg.MavlinkVersion),
	}
}
