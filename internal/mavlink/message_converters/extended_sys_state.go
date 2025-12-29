package message_converters

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// ExtendedSysStateToProtobuf
// Converts a MAVLink EXTENDED_SYS_STATE message to a protobuf ExtendedSysState message.
// MAVLink enum values map directly to protobuf enum values (MAVLink 0/UNDEFINED maps to proto 0/UNSPECIFIED).
func ExtendedSysStateToProtobuf(msg *common.MessageExtendedSysState) *flightpath.ExtendedSysState {
	return &flightpath.ExtendedSysState{
		VtolState:   flightpath.MavVtolState(msg.VtolState),
		LandedState: flightpath.MavLandedState(msg.LandedState),
	}
}
