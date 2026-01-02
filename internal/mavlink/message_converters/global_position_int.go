package message_converters

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// GlobalPositionIntToProtobuf
// Converts a MAVLink GLOBAL_POSITION_INT message to a protobuf GlobalPositionInt message.
func GlobalPositionIntToProtobuf(msg *common.MessageGlobalPositionInt) *flightpath.GlobalPositionInt {
	return &flightpath.GlobalPositionInt{
		TimeBootMs:  uint32(msg.TimeBootMs),
		Lat:         int32(msg.Lat),
		Lon:         int32(msg.Lon),
		Alt:         int32(msg.Alt),
		RelativeAlt: int32(msg.RelativeAlt),
		Vx:          int32(msg.Vx),
		Vy:          int32(msg.Vy),
		Vz:          int32(msg.Vz),
		Hdg:         uint32(msg.Hdg),
	}
}
