package message_converters

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// SysStatusToProtobuf
// Converts a MAVLink SYS_STATUS message to a protobuf SysStatus message.
func SysStatusToProtobuf(msg *common.MessageSysStatus) *flightpath.SysStatus {
	return &flightpath.SysStatus{
		OnboardControlSensorsPresent: MavSysStatusSensorToProtobuf(uint32(msg.OnboardControlSensorsPresent)),
		OnboardControlSensorsEnabled: MavSysStatusSensorToProtobuf(uint32(msg.OnboardControlSensorsEnabled)),
		OnboardControlSensorsHealth:  MavSysStatusSensorToProtobuf(uint32(msg.OnboardControlSensorsHealth)),
		Load:                         uint32(msg.Load),
		VoltageBattery:               uint32(msg.VoltageBattery),
		CurrentBattery:               int32(msg.CurrentBattery),
		BatteryRemaining:             int32(msg.BatteryRemaining),
		DropRateComm:                 uint32(msg.DropRateComm),
		ErrorsComm:                   uint32(msg.ErrorsComm),
		ErrorsCount1:                 uint32(msg.ErrorsCount1),
		ErrorsCount2:                 uint32(msg.ErrorsCount2),
		ErrorsCount3:                 uint32(msg.ErrorsCount3),
		ErrorsCount4:                 uint32(msg.ErrorsCount4),
	}
}
