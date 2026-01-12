package message_converters

import (
	"testing"

	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

func TestHeartbeatToProtobuf(t *testing.T) {
	t.Run("PX4 quadrotor in auto takeoff", func(t *testing.T) {
		msg := &common.MessageHeartbeat{
			Type:           common.MAV_TYPE_QUADROTOR,
			Autopilot:      common.MAV_AUTOPILOT_PX4,
			BaseMode:       0x85,       // custom_mode + auto + armed
			CustomMode:     0x02040000, // AUTO/TAKEOFF
			SystemStatus:   common.MAV_STATE_ACTIVE,
			MavlinkVersion: 3,
		}

		result := HeartbeatToProtobuf(msg)

		if result.Type != flightpath.MavType_MAV_TYPE_QUADROTOR {
			t.Errorf("Type: got %v, want MAV_TYPE_QUADROTOR", result.Type)
		}
		if result.Autopilot != flightpath.MavAutopilot_MAV_AUTOPILOT_PX4 {
			t.Errorf("Autopilot: got %v, want MAV_AUTOPILOT_PX4", result.Autopilot)
		}
		if result.SystemStatus != flightpath.MavState_MAV_STATE_ACTIVE {
			t.Errorf("SystemStatus: got %v, want MAV_STATE_ACTIVE", result.SystemStatus)
		}
		if result.MavlinkVersion != 3 {
			t.Errorf("MavlinkVersion: got %v, want 3", result.MavlinkVersion)
		}

		// Check base mode
		if result.BaseMode == nil {
			t.Fatal("BaseMode is nil")
		}
		if !result.BaseMode.CustomModeEnabled {
			t.Error("BaseMode.CustomModeEnabled should be true")
		}
		if !result.BaseMode.AutoEnabled {
			t.Error("BaseMode.AutoEnabled should be true")
		}
		if !result.BaseMode.SafetyArmed {
			t.Error("BaseMode.SafetyArmed should be true")
		}

		// Check custom mode
		if result.CustomMode == nil {
			t.Fatal("CustomMode is nil")
		}
		if result.CustomMode.MainMode != flightpath.MainMode_MAIN_MODE_AUTO {
			t.Errorf("CustomMode.MainMode: got %v, want MAIN_MODE_AUTO", result.CustomMode.MainMode)
		}
		if result.CustomMode.SubMode != flightpath.SubMode_SUB_MODE_AUTO_TAKEOFF {
			t.Errorf("CustomMode.SubMode: got %v, want SUB_MODE_AUTO_TAKEOFF", result.CustomMode.SubMode)
		}
	})

	t.Run("disarmed standby", func(t *testing.T) {
		msg := &common.MessageHeartbeat{
			Type:           common.MAV_TYPE_QUADROTOR,
			Autopilot:      common.MAV_AUTOPILOT_PX4,
			BaseMode:       0x01,       // only custom_mode enabled, not armed
			CustomMode:     0x03040000, // AUTO/LOITER
			SystemStatus:   common.MAV_STATE_STANDBY,
			MavlinkVersion: 3,
		}

		result := HeartbeatToProtobuf(msg)

		if result.SystemStatus != flightpath.MavState_MAV_STATE_STANDBY {
			t.Errorf("SystemStatus: got %v, want MAV_STATE_STANDBY", result.SystemStatus)
		}
		if result.BaseMode.SafetyArmed {
			t.Error("BaseMode.SafetyArmed should be false")
		}
	})

	t.Run("GCS heartbeat", func(t *testing.T) {
		msg := &common.MessageHeartbeat{
			Type:           common.MAV_TYPE_GCS,
			Autopilot:      common.MAV_AUTOPILOT_INVALID,
			BaseMode:       0,
			CustomMode:     0,
			SystemStatus:   common.MAV_STATE_ACTIVE,
			MavlinkVersion: 3,
		}

		result := HeartbeatToProtobuf(msg)

		if result.Type != flightpath.MavType_MAV_TYPE_GCS {
			t.Errorf("Type: got %v, want MAV_TYPE_GCS", result.Type)
		}
	})

	t.Run("fixed wing", func(t *testing.T) {
		msg := &common.MessageHeartbeat{
			Type:           common.MAV_TYPE_FIXED_WING,
			Autopilot:      common.MAV_AUTOPILOT_PX4,
			BaseMode:       0x85,
			CustomMode:     0x04040000, // AUTO/MISSION
			SystemStatus:   common.MAV_STATE_ACTIVE,
			MavlinkVersion: 3,
		}

		result := HeartbeatToProtobuf(msg)

		if result.Type != flightpath.MavType_MAV_TYPE_FIXED_WING {
			t.Errorf("Type: got %v, want MAV_TYPE_FIXED_WING", result.Type)
		}
		if result.CustomMode.SubMode != flightpath.SubMode_SUB_MODE_AUTO_MISSION {
			t.Errorf("CustomMode.SubMode: got %v, want SUB_MODE_AUTO_MISSION", result.CustomMode.SubMode)
		}
	})
}
