package message_converters

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// MissionItemReachedToProtobuf
// Converts a MAVLink MISSION_ITEM_REACHED message to a protobuf MissionItemReached message.
func MissionItemReachedToProtobuf(msg *common.MessageMissionItemReached) *flightpath.MissionItemReached {
	return &flightpath.MissionItemReached{
		Seq: uint32(msg.Seq),
	}
}
