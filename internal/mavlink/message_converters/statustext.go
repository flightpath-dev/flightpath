package message_converters

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// StatusTextToProtobuf
// Converts a MAVLink STATUSTEXT message to a protobuf StatusText message.
// The id and chunk_seq fields are ignored as they are only used for message reassembly.
func StatusTextToProtobuf(msg *common.MessageStatustext) *flightpath.StatusText {
	return &flightpath.StatusText{
		Severity: MavSeverityToProtobuf(msg.Severity),
		Text:     msg.Text,
	}
}

// MavSeverityToProtobuf
// Converts MAVLink MAV_SEVERITY to protobuf MavSeverity enum.
// Proto enum values are incremented by 1 to accommodate MAV_SEVERITY_UNSPECIFIED at 0.
// MAVLink 0 (EMERGENCY) maps to proto 1 (EMERGENCY), MAVLink 1 (ALERT) maps to proto 2, etc.
func MavSeverityToProtobuf(severity common.MAV_SEVERITY) flightpath.MavSeverity {
	// Add 1 to MAVLink value to account for UNSPECIFIED at 0 in proto
	return flightpath.MavSeverity(severity + 1)
}
