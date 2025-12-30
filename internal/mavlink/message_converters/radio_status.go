package message_converters

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// RadioStatusToProtobuf
// Converts a MAVLink RADIO_STATUS message to a protobuf RadioStatus message.
func RadioStatusToProtobuf(msg *common.MessageRadioStatus) *flightpath.RadioStatus {
	return &flightpath.RadioStatus{
		Rssi:     uint32(msg.Rssi),
		Remrssi:  uint32(msg.Remrssi),
		Txbuf:    uint32(msg.Txbuf),
		Noise:    uint32(msg.Noise),
		Remnoise: uint32(msg.Remnoise),
		Rxerrors: uint32(msg.Rxerrors),
		Fixed:    uint32(msg.Fixed),
	}
}

