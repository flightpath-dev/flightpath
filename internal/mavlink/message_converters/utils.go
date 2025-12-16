package message_converters

import (
	"fmt"
	"strings"

	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// BaseModeToProtobuf
// Converts MAVLink base_mode bitfield to protobuf BaseMode structured message.
// MAVLink MAV_MODE_FLAG bit positions (from MAVLink spec):
// Bit 7 (128): MAV_MODE_FLAG_SAFETY_ARMED
// Bit 6 (64):  MAV_MODE_FLAG_MANUAL_INPUT_ENABLED
// Bit 5 (32):  MAV_MODE_FLAG_HIL_ENABLED
// Bit 4 (16):  MAV_MODE_FLAG_STABILIZE_ENABLED
// Bit 3 (8):   MAV_MODE_FLAG_GUIDED_ENABLED
// Bit 2 (4):   MAV_MODE_FLAG_AUTO_ENABLED
// Bit 1 (2):   MAV_MODE_FLAG_TEST_ENABLED
// Bit 0 (1):   MAV_MODE_FLAG_CUSTOM_MODE_ENABLED
func BaseModeToProtobuf(baseMode common.MAV_MODE_FLAG) *flightpath.BaseMode {
	const (
		MAV_MODE_FLAG_SAFETY_ARMED         = 128 // 0x80
		MAV_MODE_FLAG_MANUAL_INPUT_ENABLED = 64  // 0x40
		MAV_MODE_FLAG_HIL_ENABLED          = 32  // 0x20
		MAV_MODE_FLAG_STABILIZE_ENABLED    = 16  // 0x10
		MAV_MODE_FLAG_GUIDED_ENABLED       = 8   // 0x08
		MAV_MODE_FLAG_AUTO_ENABLED         = 4   // 0x04
		MAV_MODE_FLAG_TEST_ENABLED         = 2   // 0x02
		MAV_MODE_FLAG_CUSTOM_MODE_ENABLED  = 1   // 0x01
	)

	baseModeUint8 := uint8(baseMode)
	return &flightpath.BaseMode{
		CustomModeEnabled:  (baseModeUint8 & MAV_MODE_FLAG_CUSTOM_MODE_ENABLED) != 0,
		TestEnabled:        (baseModeUint8 & MAV_MODE_FLAG_TEST_ENABLED) != 0,
		AutoEnabled:        (baseModeUint8 & MAV_MODE_FLAG_AUTO_ENABLED) != 0,
		GuidedEnabled:      (baseModeUint8 & MAV_MODE_FLAG_GUIDED_ENABLED) != 0,
		StabilizeEnabled:   (baseModeUint8 & MAV_MODE_FLAG_STABILIZE_ENABLED) != 0,
		HilEnabled:         (baseModeUint8 & MAV_MODE_FLAG_HIL_ENABLED) != 0,
		ManualInputEnabled: (baseModeUint8 & MAV_MODE_FLAG_MANUAL_INPUT_ENABLED) != 0,
		SafetyArmed:        (baseModeUint8 & MAV_MODE_FLAG_SAFETY_ARMED) != 0,
	}
}

// CustomModeToProtobuf
// Converts MAVLink custom_mode uint32 to protobuf CustomMode message.
// For PX4 autopilots, decodes the custom_mode into main_mode and sub_mode.
// For other autopilots, returns unspecified values.
func CustomModeToProtobuf(customMode uint32, autopilot common.MAV_AUTOPILOT) *flightpath.CustomMode {
	if autopilot == common.MAV_AUTOPILOT_PX4 {
		// Extract main_mode and sub_mode from PX4 custom_mode uint32
		px4MainMode := uint8((customMode >> 16) & 0xFF)
		px4SubMode := uint8((customMode >> 24) & 0xFF)

		mainMode := FlightMainModeToProtobuf(px4MainMode)
		subMode := FlightSubModeToProtobuf(px4SubMode)

		return &flightpath.CustomMode{
			MainMode: mainMode,
			SubMode:  subMode,
		}
	}

	// For non-PX4 autopilots, set to unspecified
	return &flightpath.CustomMode{
		MainMode: flightpath.MainMode_MAIN_MODE_UNSPECIFIED,
		SubMode:  flightpath.SubMode_SUB_MODE_UNSPECIFIED,
	}
}

// FlightMainModeToProtobuf
// Converts flight main mode (uint8, 1-based) to protobuf MainMode enum.
// Direct mapping: 0→UNSPECIFIED, 1→MANUAL, 2→ALTCTL, etc.
func FlightMainModeToProtobuf(mode uint8) flightpath.MainMode {
	return flightpath.MainMode(mode)
}

// FlightSubModeToProtobuf
// Converts flight sub-mode (uint8, 1-based) to protobuf SubMode enum.
// Direct mapping: 0→UNSPECIFIED, 1→AUTO_READY, 2→AUTO_TAKEOFF, etc.
func FlightSubModeToProtobuf(subMode uint8) flightpath.SubMode {
	return flightpath.SubMode(subMode)
}

// GpsFixTypeToProtobuf
// Converts MAVLink GPS_FIX_TYPE to protobuf GpsFixType enum.
// Proto enum values are incremented by 1 to accommodate GPS_FIX_TYPE_UNSPECIFIED at 0.
// MAVLink 0 (NO_GPS) maps to proto 1 (NO_GPS), MAVLink 1 (NO_FIX) maps to proto 2 (NO_FIX), etc.
func GpsFixTypeToProtobuf(fixType common.GPS_FIX_TYPE) flightpath.GpsFixType {
	// Add 1 to MAVLink value to account for UNSPECIFIED at 0 in proto
	return flightpath.GpsFixType(fixType + 1)
}

// MavAutopilotToProtobuf
// Converts MAVLink MAV_AUTOPILOT to protobuf MavAutopilot enum.
// Note: Direct cast assumes MAVLink enum values match protobuf enum values.
// If values don't match, explicit mapping may be needed.
func MavAutopilotToProtobuf(autopilot common.MAV_AUTOPILOT) flightpath.MavAutopilot {
	// Direct mapping: MAVLink enum values should match protobuf enum values
	// MAV_AUTOPILOT_PX4 = 13 in both MAVLink and protobuf
	return flightpath.MavAutopilot(autopilot)
}

// MavStateToProtobuf
// Converts MAVLink MAV_STATE to protobuf MavState enum.
// Note: MAV_STATE_UNINIT (0) maps to MAV_STATE_UNSPECIFIED (0) in protobuf
func MavStateToProtobuf(state common.MAV_STATE) flightpath.MavState {
	// Direct mapping: MAVLink enum values match protobuf enum values
	// UNINIT (0) = UNSPECIFIED (0)
	return flightpath.MavState(state)
}

// MavTypeToProtobuf
// Converts MAVLink MAV_TYPE to protobuf MavType enum.
// Note: Direct cast assumes MAVLink enum values match protobuf enum values.
// If values don't match, explicit mapping may be needed.
func MavTypeToProtobuf(mavType common.MAV_TYPE) flightpath.MavType {
	// Direct mapping: MAVLink enum values should match protobuf enum values
	// MAV_TYPE_QUADROTOR = 3 in both MAVLink and protobuf
	return flightpath.MavType(mavType)
}

// DecodePX4CustomMode
// Decodes PX4 CustomMode uint32 into human-readable format.
// Based on: https://github.com/PX4/PX4-Autopilot/blob/main/src/modules/commander/px4_custom_mode.h
//
// Structure:
//
//	SM MM 00 00
//	^  ^  -----
//	|  |  Reserved
//	|  |
//	|  Bits 16-23: main mode
//	|
//	Bits 24-31: sub mode
func DecodePX4CustomMode(customMode uint32) map[string]interface{} {
	// Extract main_mode (bits 16-23) and sub_mode (bits 24-31)
	px4MainMode := uint8((customMode >> 16) & 0xFF)
	px4SubMode := uint8((customMode >> 24) & 0xFF)

	// Convert to protobuf enums
	mainMode := FlightMainModeToProtobuf(px4MainMode)
	subMode := FlightSubModeToProtobuf(px4SubMode)

	result := map[string]interface{}{
		"raw":           fmt.Sprintf("0x%08X", customMode),
		"main_mode":     fmt.Sprintf("0x%02X", px4MainMode),
		"main_mode_str": mainModeString(mainMode),
		"sub_mode":      fmt.Sprintf("0x%02X", px4SubMode),
		"sub_mode_str":  subModeString(subMode),
	}

	return result
}

// mainModeString
// Helper function to get clean string name from MainMode enum (strips "MAIN_MODE_" prefix).
func mainModeString(mode flightpath.MainMode) string {
	name := mode.String()
	// Remove "MAIN_MODE_" prefix if present
	if strings.HasPrefix(name, "MAIN_MODE_") {
		return strings.TrimPrefix(name, "MAIN_MODE_")
	}
	return name
}

// subModeString
// Helper function to get clean string name from SubMode enum (strips "SUB_MODE_" prefix).
func subModeString(mode flightpath.SubMode) string {
	name := mode.String()
	// Remove "SUB_MODE_AUTO_" prefix for AUTO sub-modes
	if strings.HasPrefix(name, "SUB_MODE_AUTO_") {
		return strings.TrimPrefix(name, "SUB_MODE_AUTO_")
	}
	// Remove "SUB_MODE_POSCTL_" prefix for POSCTL sub-modes
	if strings.HasPrefix(name, "SUB_MODE_POSCTL_") {
		return strings.TrimPrefix(name, "SUB_MODE_POSCTL_")
	}
	// Remove "SUB_MODE_" prefix for other sub-modes
	if strings.HasPrefix(name, "SUB_MODE_") {
		return strings.TrimPrefix(name, "SUB_MODE_")
	}
	return name
}
