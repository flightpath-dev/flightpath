package message_converters

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// MissionCurrentToProtobuf
// Converts a MAVLink MISSION_CURRENT message to a protobuf MissionCurrent message.
// MAVLink enum values map directly to protobuf enum values (MAVLink 0/UNKNOWN maps to proto 0/UNSPECIFIED).
func MissionCurrentToProtobuf(msg *common.MessageMissionCurrent) *flightpath.MissionCurrent {
	return &flightpath.MissionCurrent{
		Seq:           uint32(msg.Seq),
		Total:         uint32(msg.Total),
		MissionState:  MissionStateToProtobuf(msg.MissionState),
		MissionMode:   uint32(msg.MissionMode),
		MissionId:     msg.MissionId,
		FenceId:       msg.FenceId,
		RallyPointsId: msg.RallyPointsId,
	}
}

// MissionStateToProtobuf
// Converts MAVLink MISSION_STATE to protobuf MissionState enum.
// MAVLink enum values map directly to protobuf enum values (MAVLink 0/UNKNOWN maps to proto 0/UNSPECIFIED).
func MissionStateToProtobuf(state common.MISSION_STATE) flightpath.MissionState {
	return flightpath.MissionState(state)
}
