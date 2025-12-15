package mavlink

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

// Px4MainModeToProtobuf
// Converts PX4 main mode (uint8, 1-based) to protobuf MainMode enum.
func Px4MainModeToProtobuf(px4Mode uint8) flightpath.MainMode {
	// PX4 modes are 1-based, protobuf enums are 1-based (with 0 = UNSPECIFIED)
	// Direct mapping: 1→MANUAL, 2→ALTCTL, etc.
	if px4Mode == 0 {
		return flightpath.MainMode_MAIN_MODE_UNSPECIFIED
	}
	return flightpath.MainMode(px4Mode)
}

// Px4AutoSubModeToProtobuf
// Converts PX4 AUTO sub-mode (uint8, 1-based) to protobuf SubMode enum.
func Px4AutoSubModeToProtobuf(px4SubMode uint8) flightpath.SubMode {
	// PX4 AUTO sub-modes are 1-based, protobuf SubMode AUTO sub-modes are 1-based
	// Direct mapping: 1→READY, 2→TAKEOFF, etc.
	if px4SubMode == 0 {
		return flightpath.SubMode_SUB_MODE_UNSPECIFIED
	}
	return flightpath.SubMode(px4SubMode)
}

// Px4PosctlSubModeToProtobuf
// Converts PX4 POSCTL sub-mode (uint8, 0-based) to protobuf SubMode enum.
func Px4PosctlSubModeToProtobuf(px4SubMode uint8) flightpath.SubMode {
	// PX4 POSCTL sub-modes are 0-based (0=POSCTL, 1=ORBIT, 2=SLOW)
	// Protobuf POSCTL sub-modes start at 10 (10=POSCTL, 11=ORBIT, 12=SLOW)
	switch px4SubMode {
	case 0:
		return flightpath.SubMode_SUB_MODE_POSCTL
	case 1:
		return flightpath.SubMode_SUB_MODE_ORBIT
	case 2:
		return flightpath.SubMode_SUB_MODE_SLOW
	default:
		return flightpath.SubMode_SUB_MODE_UNSPECIFIED
	}
}

// Helper function to get clean string name from MainMode enum (strips "MAIN_MODE_" prefix)
func mainModeString(mode flightpath.MainMode) string {
	name := mode.String()
	// Remove "MAIN_MODE_" prefix if present
	if strings.HasPrefix(name, "MAIN_MODE_") {
		return strings.TrimPrefix(name, "MAIN_MODE_")
	}
	return name
}

// Helper function to get clean string name from SubMode enum (strips "SUB_MODE_" prefix)
func subModeString(mode flightpath.SubMode) string {
	name := mode.String()
	// Remove "SUB_MODE_" prefix if present
	if strings.HasPrefix(name, "SUB_MODE_") {
		return strings.TrimPrefix(name, "SUB_MODE_")
	}
	return name
}

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
	mainMode := Px4MainModeToProtobuf(px4MainMode)
	var subMode flightpath.SubMode

	// Convert sub mode based on main mode
	switch mainMode {
	case flightpath.MainMode_MAIN_MODE_AUTO:
		subMode = Px4AutoSubModeToProtobuf(px4SubMode)
	case flightpath.MainMode_MAIN_MODE_POSCTL:
		subMode = Px4PosctlSubModeToProtobuf(px4SubMode)
	default:
		subMode = flightpath.SubMode_SUB_MODE_UNSPECIFIED
	}

	result := map[string]interface{}{
		"raw":           fmt.Sprintf("0x%08X", customMode),
		"main_mode":     fmt.Sprintf("0x%02X", px4MainMode),
		"main_mode_str": mainModeString(mainMode),
		"sub_mode":      fmt.Sprintf("0x%02X", px4SubMode),
		"sub_mode_str":  subModeString(subMode),
	}

	return result
}

// HeartbeatMessageToMap
// Converts a HEARTBEAT message to a map with decoded fields for better readability.
// For example, the PX4 CustomMode is decoded into a human-readable format.
func HeartbeatMessageToMap(msg *common.MessageHeartbeat) (map[string]interface{}, error) {
	// Convert message to JSON first
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var msgMap map[string]interface{}
	if err := json.Unmarshal(msgJSON, &msgMap); err != nil {
		return nil, err
	}

	// Decode CustomMode if this is a PX4 autopilot
	if msg.Autopilot == common.MAV_AUTOPILOT_PX4 {
		decoded := DecodePX4CustomMode(msg.CustomMode)
		msgMap["CustomModeDecoded"] = decoded
	}

	return msgMap, nil
}
