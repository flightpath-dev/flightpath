package message_converters

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// VfrHudToProtobuf
// Converts a MAVLink VFR_HUD message to a protobuf VfrHud message.
func VfrHudToProtobuf(msg *common.MessageVfrHud) *flightpath.VfrHud {
	return &flightpath.VfrHud{
		Airspeed:    msg.Airspeed,
		Groundspeed: msg.Groundspeed,
		Heading:     int32(msg.Heading),
		Throttle:    uint32(msg.Throttle),
		Alt:         msg.Alt,
		Climb:       msg.Climb,
	}
}

