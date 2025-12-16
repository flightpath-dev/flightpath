package message_converters

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// GpsRawIntToProtobuf
// Converts a MAVLink GPS_RAW_INT message to a protobuf GpsRawInt message.
func GpsRawIntToProtobuf(msg *common.MessageGpsRawInt) *flightpath.GpsRawInt {
	return &flightpath.GpsRawInt{
		TimeUsec:          msg.TimeUsec,
		FixType:           GpsFixTypeToProtobuf(msg.FixType),
		Lat:               msg.Lat,
		Lon:               msg.Lon,
		Alt:               msg.Alt,
		Eph:               uint32(msg.Eph),
		Epv:               uint32(msg.Epv),
		Vel:               uint32(msg.Vel),
		Cog:               uint32(msg.Cog),
		SatellitesVisible: uint32(msg.SatellitesVisible),
		AltEllipsoid:      msg.AltEllipsoid,
		HAcc:              msg.HAcc,
		VAcc:              msg.VAcc,
		VelAcc:            msg.VelAcc,
		HdgAcc:            msg.HdgAcc,
		Yaw:               uint32(msg.Yaw),
	}
}
