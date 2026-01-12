package message_converters

import (
	"testing"

	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

func TestBaseModeToProtobuf(t *testing.T) {
	tests := []struct {
		name     string
		baseMode common.MAV_MODE_FLAG
		expected *flightpath.BaseMode
	}{
		{
			name:     "all flags off",
			baseMode: 0,
			expected: &flightpath.BaseMode{
				CustomModeEnabled:  false,
				TestEnabled:        false,
				AutoEnabled:        false,
				GuidedEnabled:      false,
				StabilizeEnabled:   false,
				HilEnabled:         false,
				ManualInputEnabled: false,
				SafetyArmed:        false,
			},
		},
		{
			name:     "all flags on",
			baseMode: 0xFF,
			expected: &flightpath.BaseMode{
				CustomModeEnabled:  true,
				TestEnabled:        true,
				AutoEnabled:        true,
				GuidedEnabled:      true,
				StabilizeEnabled:   true,
				HilEnabled:         true,
				ManualInputEnabled: true,
				SafetyArmed:        true,
			},
		},
		{
			name:     "only custom mode enabled (bit 0)",
			baseMode: 0x01,
			expected: &flightpath.BaseMode{
				CustomModeEnabled: true,
			},
		},
		{
			name:     "only safety armed (bit 7)",
			baseMode: 0x80,
			expected: &flightpath.BaseMode{
				SafetyArmed: true,
			},
		},
		{
			name:     "armed with custom mode (bits 0 and 7)",
			baseMode: 0x81,
			expected: &flightpath.BaseMode{
				CustomModeEnabled: true,
				SafetyArmed:       true,
			},
		},
		{
			name:     "typical armed auto mode",
			baseMode: 0x85, // custom_mode + auto + armed
			expected: &flightpath.BaseMode{
				CustomModeEnabled: true,
				AutoEnabled:       true,
				SafetyArmed:       true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BaseModeToProtobuf(tt.baseMode)
			if result.CustomModeEnabled != tt.expected.CustomModeEnabled {
				t.Errorf("CustomModeEnabled: got %v, want %v", result.CustomModeEnabled, tt.expected.CustomModeEnabled)
			}
			if result.TestEnabled != tt.expected.TestEnabled {
				t.Errorf("TestEnabled: got %v, want %v", result.TestEnabled, tt.expected.TestEnabled)
			}
			if result.AutoEnabled != tt.expected.AutoEnabled {
				t.Errorf("AutoEnabled: got %v, want %v", result.AutoEnabled, tt.expected.AutoEnabled)
			}
			if result.GuidedEnabled != tt.expected.GuidedEnabled {
				t.Errorf("GuidedEnabled: got %v, want %v", result.GuidedEnabled, tt.expected.GuidedEnabled)
			}
			if result.StabilizeEnabled != tt.expected.StabilizeEnabled {
				t.Errorf("StabilizeEnabled: got %v, want %v", result.StabilizeEnabled, tt.expected.StabilizeEnabled)
			}
			if result.HilEnabled != tt.expected.HilEnabled {
				t.Errorf("HilEnabled: got %v, want %v", result.HilEnabled, tt.expected.HilEnabled)
			}
			if result.ManualInputEnabled != tt.expected.ManualInputEnabled {
				t.Errorf("ManualInputEnabled: got %v, want %v", result.ManualInputEnabled, tt.expected.ManualInputEnabled)
			}
			if result.SafetyArmed != tt.expected.SafetyArmed {
				t.Errorf("SafetyArmed: got %v, want %v", result.SafetyArmed, tt.expected.SafetyArmed)
			}
		})
	}
}

func TestCustomModeToProtobuf_PX4(t *testing.T) {
	tests := []struct {
		name             string
		customMode       uint32
		expectedMainMode flightpath.MainMode
		expectedSubMode  flightpath.SubMode
	}{
		{
			name:             "manual mode",
			customMode:       0x00010000, // main_mode=1, sub_mode=0
			expectedMainMode: flightpath.MainMode_MAIN_MODE_MANUAL,
			expectedSubMode:  flightpath.SubMode_SUB_MODE_UNSPECIFIED,
		},
		{
			name:             "auto takeoff",
			customMode:       0x02040000, // main_mode=4 (AUTO), sub_mode=2 (TAKEOFF)
			expectedMainMode: flightpath.MainMode_MAIN_MODE_AUTO,
			expectedSubMode:  flightpath.SubMode_SUB_MODE_AUTO_TAKEOFF,
		},
		{
			name:             "auto loiter",
			customMode:       0x03040000, // main_mode=4 (AUTO), sub_mode=3 (LOITER)
			expectedMainMode: flightpath.MainMode_MAIN_MODE_AUTO,
			expectedSubMode:  flightpath.SubMode_SUB_MODE_AUTO_LOITER,
		},
		{
			name:             "auto mission",
			customMode:       0x04040000, // main_mode=4 (AUTO), sub_mode=4 (MISSION)
			expectedMainMode: flightpath.MainMode_MAIN_MODE_AUTO,
			expectedSubMode:  flightpath.SubMode_SUB_MODE_AUTO_MISSION,
		},
		{
			name:             "auto RTL",
			customMode:       0x05040000, // main_mode=4 (AUTO), sub_mode=5 (RTL)
			expectedMainMode: flightpath.MainMode_MAIN_MODE_AUTO,
			expectedSubMode:  flightpath.SubMode_SUB_MODE_AUTO_RTL,
		},
		{
			name:             "position control",
			customMode:       0x00030000, // main_mode=3 (POSCTL), sub_mode=0
			expectedMainMode: flightpath.MainMode_MAIN_MODE_POSCTL,
			expectedSubMode:  flightpath.SubMode_SUB_MODE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CustomModeToProtobuf(tt.customMode, common.MAV_AUTOPILOT_PX4)
			if result.MainMode != tt.expectedMainMode {
				t.Errorf("MainMode: got %v, want %v", result.MainMode, tt.expectedMainMode)
			}
			if result.SubMode != tt.expectedSubMode {
				t.Errorf("SubMode: got %v, want %v", result.SubMode, tt.expectedSubMode)
			}
		})
	}
}

func TestCustomModeToProtobuf_NonPX4(t *testing.T) {
	// Non-PX4 autopilots should return unspecified modes
	autopilots := []common.MAV_AUTOPILOT{
		common.MAV_AUTOPILOT_ARDUPILOTMEGA,
		common.MAV_AUTOPILOT_GENERIC,
		common.MAV_AUTOPILOT_INVALID,
	}

	for _, autopilot := range autopilots {
		t.Run(autopilot.String(), func(t *testing.T) {
			result := CustomModeToProtobuf(0x02040000, autopilot)
			if result.MainMode != flightpath.MainMode_MAIN_MODE_UNSPECIFIED {
				t.Errorf("MainMode: got %v, want UNSPECIFIED", result.MainMode)
			}
			if result.SubMode != flightpath.SubMode_SUB_MODE_UNSPECIFIED {
				t.Errorf("SubMode: got %v, want UNSPECIFIED", result.SubMode)
			}
		})
	}
}

func TestGpsFixTypeToProtobuf(t *testing.T) {
	tests := []struct {
		name     string
		fixType  common.GPS_FIX_TYPE
		expected flightpath.GpsFixType
	}{
		{
			name:     "no GPS",
			fixType:  common.GPS_FIX_TYPE_NO_GPS,
			expected: flightpath.GpsFixType_GPS_FIX_TYPE_NO_GPS,
		},
		{
			name:     "no fix",
			fixType:  common.GPS_FIX_TYPE_NO_FIX,
			expected: flightpath.GpsFixType_GPS_FIX_TYPE_NO_FIX,
		},
		{
			name:     "2D fix",
			fixType:  common.GPS_FIX_TYPE_2D_FIX,
			expected: flightpath.GpsFixType_GPS_FIX_TYPE_2D_FIX,
		},
		{
			name:     "3D fix",
			fixType:  common.GPS_FIX_TYPE_3D_FIX,
			expected: flightpath.GpsFixType_GPS_FIX_TYPE_3D_FIX,
		},
		{
			name:     "DGPS",
			fixType:  common.GPS_FIX_TYPE_DGPS,
			expected: flightpath.GpsFixType_GPS_FIX_TYPE_DGPS,
		},
		{
			name:     "RTK float",
			fixType:  common.GPS_FIX_TYPE_RTK_FLOAT,
			expected: flightpath.GpsFixType_GPS_FIX_TYPE_RTK_FLOAT,
		},
		{
			name:     "RTK fixed",
			fixType:  common.GPS_FIX_TYPE_RTK_FIXED,
			expected: flightpath.GpsFixType_GPS_FIX_TYPE_RTK_FIXED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GpsFixTypeToProtobuf(tt.fixType)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMavAutopilotToProtobuf(t *testing.T) {
	tests := []struct {
		name      string
		autopilot common.MAV_AUTOPILOT
		expected  flightpath.MavAutopilot
	}{
		{
			name:      "generic",
			autopilot: common.MAV_AUTOPILOT_GENERIC,
			expected:  flightpath.MavAutopilot_MAV_AUTOPILOT_UNSPECIFIED,
		},
		{
			name:      "PX4",
			autopilot: common.MAV_AUTOPILOT_PX4,
			expected:  flightpath.MavAutopilot_MAV_AUTOPILOT_PX4,
		},
		{
			name:      "ArduPilot",
			autopilot: common.MAV_AUTOPILOT_ARDUPILOTMEGA,
			expected:  flightpath.MavAutopilot_MAV_AUTOPILOT_ARDUPILOTMEGA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MavAutopilotToProtobuf(tt.autopilot)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMavStateToProtobuf(t *testing.T) {
	tests := []struct {
		name     string
		state    common.MAV_STATE
		expected flightpath.MavState
	}{
		{
			name:     "uninit",
			state:    common.MAV_STATE_UNINIT,
			expected: flightpath.MavState_MAV_STATE_UNSPECIFIED,
		},
		{
			name:     "boot",
			state:    common.MAV_STATE_BOOT,
			expected: flightpath.MavState_MAV_STATE_BOOT,
		},
		{
			name:     "standby",
			state:    common.MAV_STATE_STANDBY,
			expected: flightpath.MavState_MAV_STATE_STANDBY,
		},
		{
			name:     "active",
			state:    common.MAV_STATE_ACTIVE,
			expected: flightpath.MavState_MAV_STATE_ACTIVE,
		},
		{
			name:     "critical",
			state:    common.MAV_STATE_CRITICAL,
			expected: flightpath.MavState_MAV_STATE_CRITICAL,
		},
		{
			name:     "emergency",
			state:    common.MAV_STATE_EMERGENCY,
			expected: flightpath.MavState_MAV_STATE_EMERGENCY,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MavStateToProtobuf(tt.state)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMavTypeToProtobuf(t *testing.T) {
	tests := []struct {
		name     string
		mavType  common.MAV_TYPE
		expected flightpath.MavType
	}{
		{
			name:     "generic",
			mavType:  common.MAV_TYPE_GENERIC,
			expected: flightpath.MavType_MAV_TYPE_UNSPECIFIED,
		},
		{
			name:     "fixed wing",
			mavType:  common.MAV_TYPE_FIXED_WING,
			expected: flightpath.MavType_MAV_TYPE_FIXED_WING,
		},
		{
			name:     "quadrotor",
			mavType:  common.MAV_TYPE_QUADROTOR,
			expected: flightpath.MavType_MAV_TYPE_QUADROTOR,
		},
		{
			name:     "hexarotor",
			mavType:  common.MAV_TYPE_HEXAROTOR,
			expected: flightpath.MavType_MAV_TYPE_HEXAROTOR,
		},
		{
			name:     "GCS",
			mavType:  common.MAV_TYPE_GCS,
			expected: flightpath.MavType_MAV_TYPE_GCS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MavTypeToProtobuf(tt.mavType)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMavSysStatusSensorToProtobuf(t *testing.T) {
	t.Run("all sensors off", func(t *testing.T) {
		result := MavSysStatusSensorToProtobuf(0)
		if result.Sensor_3DGyro || result.SensorGps || result.SensorBattery {
			t.Error("expected all sensors to be off")
		}
	})

	t.Run("gyro enabled (bit 0)", func(t *testing.T) {
		result := MavSysStatusSensorToProtobuf(0x01)
		if !result.Sensor_3DGyro {
			t.Error("expected 3D gyro to be enabled")
		}
		if result.Sensor_3DAccel {
			t.Error("expected 3D accel to be disabled")
		}
	})

	t.Run("GPS enabled (bit 5)", func(t *testing.T) {
		result := MavSysStatusSensorToProtobuf(0x20)
		if !result.SensorGps {
			t.Error("expected GPS to be enabled")
		}
	})

	t.Run("battery enabled (bit 25)", func(t *testing.T) {
		result := MavSysStatusSensorToProtobuf(0x2000000)
		if !result.SensorBattery {
			t.Error("expected battery sensor to be enabled")
		}
	})

	t.Run("multiple sensors enabled", func(t *testing.T) {
		// gyro (0x01) + accel (0x02) + mag (0x04) + GPS (0x20)
		result := MavSysStatusSensorToProtobuf(0x27)
		if !result.Sensor_3DGyro {
			t.Error("expected 3D gyro to be enabled")
		}
		if !result.Sensor_3DAccel {
			t.Error("expected 3D accel to be enabled")
		}
		if !result.Sensor_3DMag {
			t.Error("expected 3D mag to be enabled")
		}
		if !result.SensorGps {
			t.Error("expected GPS to be enabled")
		}
		if result.SensorBattery {
			t.Error("expected battery to be disabled")
		}
	})

	t.Run("all sensors on", func(t *testing.T) {
		result := MavSysStatusSensorToProtobuf(0xFFFFFFFF)
		if !result.Sensor_3DGyro || !result.SensorGps || !result.SensorBattery || !result.ExtensionUsed {
			t.Error("expected all sensors to be on")
		}
	})
}

func TestFlightMainModeToProtobuf(t *testing.T) {
	tests := []struct {
		mode     uint8
		expected flightpath.MainMode
	}{
		{0, flightpath.MainMode_MAIN_MODE_UNSPECIFIED},
		{1, flightpath.MainMode_MAIN_MODE_MANUAL},
		{2, flightpath.MainMode_MAIN_MODE_ALTCTL},
		{3, flightpath.MainMode_MAIN_MODE_POSCTL},
		{4, flightpath.MainMode_MAIN_MODE_AUTO},
		{5, flightpath.MainMode_MAIN_MODE_ACRO},
		{6, flightpath.MainMode_MAIN_MODE_OFFBOARD},
	}

	for _, tt := range tests {
		result := FlightMainModeToProtobuf(tt.mode)
		if result != tt.expected {
			t.Errorf("FlightMainModeToProtobuf(%d) = %v, want %v", tt.mode, result, tt.expected)
		}
	}
}

func TestFlightSubModeToProtobuf(t *testing.T) {
	tests := []struct {
		subMode  uint8
		expected flightpath.SubMode
	}{
		{0, flightpath.SubMode_SUB_MODE_UNSPECIFIED},
		{1, flightpath.SubMode_SUB_MODE_AUTO_READY},
		{2, flightpath.SubMode_SUB_MODE_AUTO_TAKEOFF},
		{3, flightpath.SubMode_SUB_MODE_AUTO_LOITER},
		{4, flightpath.SubMode_SUB_MODE_AUTO_MISSION},
		{5, flightpath.SubMode_SUB_MODE_AUTO_RTL},
		{6, flightpath.SubMode_SUB_MODE_AUTO_LAND},
	}

	for _, tt := range tests {
		result := FlightSubModeToProtobuf(tt.subMode)
		if result != tt.expected {
			t.Errorf("FlightSubModeToProtobuf(%d) = %v, want %v", tt.subMode, result, tt.expected)
		}
	}
}

func TestDecodePX4CustomMode(t *testing.T) {
	result := DecodePX4CustomMode(0x02040000) // AUTO/TAKEOFF

	if result["main_mode"] != "0x04" {
		t.Errorf("main_mode: got %v, want 0x04", result["main_mode"])
	}
	if result["sub_mode"] != "0x02" {
		t.Errorf("sub_mode: got %v, want 0x02", result["sub_mode"])
	}
	if result["main_mode_str"] != "AUTO" {
		t.Errorf("main_mode_str: got %v, want AUTO", result["main_mode_str"])
	}
	if result["sub_mode_str"] != "TAKEOFF" {
		t.Errorf("sub_mode_str: got %v, want TAKEOFF", result["sub_mode_str"])
	}
}
