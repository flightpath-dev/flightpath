package mavlink

import (
	"encoding/json"
	"fmt"

	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Constants based on https://github.com/PX4/PX4-Autopilot/blob/main/src/modules/commander/px4_custom_mode.h

type PX4_CUSTOM_MAIN_MODE uint8
type PX4_CUSTOM_SUB_MODE_AUTO uint8
type PX4_CUSTOM_SUB_MODE_POSCTL uint8

// PX4_CUSTOM_MAIN_MODE represents the main flight mode
const (
	PX4_CUSTOM_MAIN_MODE_MANUAL PX4_CUSTOM_MAIN_MODE = iota + 1
	PX4_CUSTOM_MAIN_MODE_ALTCTL
	PX4_CUSTOM_MAIN_MODE_POSCTL
	PX4_CUSTOM_MAIN_MODE_AUTO
	PX4_CUSTOM_MAIN_MODE_ACRO
	PX4_CUSTOM_MAIN_MODE_OFFBOARD
	PX4_CUSTOM_MAIN_MODE_STABILIZED
	PX4_CUSTOM_MAIN_MODE_RATTITUDE_LEGACY
	PX4_CUSTOM_MAIN_MODE_SIMPLE
	PX4_CUSTOM_MAIN_MODE_TERMINATION
	PX4_CUSTOM_MAIN_MODE_ALTITUDE_CRUISE
)

// PX4_CUSTOM_SUB_MODE_AUTO represents sub-modes for AUTO main mode
const (
	PX4_CUSTOM_SUB_MODE_AUTO_READY PX4_CUSTOM_SUB_MODE_AUTO = iota + 1
	PX4_CUSTOM_SUB_MODE_AUTO_TAKEOFF
	PX4_CUSTOM_SUB_MODE_AUTO_LOITER
	PX4_CUSTOM_SUB_MODE_AUTO_MISSION
	PX4_CUSTOM_SUB_MODE_AUTO_RTL
	PX4_CUSTOM_SUB_MODE_AUTO_LAND
	PX4_CUSTOM_SUB_MODE_AUTO_RESERVED_DO_NOT_USE
	PX4_CUSTOM_SUB_MODE_AUTO_FOLLOW_TARGET
	PX4_CUSTOM_SUB_MODE_AUTO_PRECLAND
	PX4_CUSTOM_SUB_MODE_AUTO_VTOL_TAKEOFF
	PX4_CUSTOM_SUB_MODE_EXTERNAL1
	PX4_CUSTOM_SUB_MODE_EXTERNAL2
	PX4_CUSTOM_SUB_MODE_EXTERNAL3
	PX4_CUSTOM_SUB_MODE_EXTERNAL4
	PX4_CUSTOM_SUB_MODE_EXTERNAL5
	PX4_CUSTOM_SUB_MODE_EXTERNAL6
	PX4_CUSTOM_SUB_MODE_EXTERNAL7
	PX4_CUSTOM_SUB_MODE_EXTERNAL8
)

// PX4_CUSTOM_SUB_MODE_POSCTL represents sub-modes for POSCTL main mode
const (
	PX4_CUSTOM_SUB_MODE_POSCTL_POSCTL PX4_CUSTOM_SUB_MODE_POSCTL = iota
	PX4_CUSTOM_SUB_MODE_POSCTL_ORBIT
	PX4_CUSTOM_SUB_MODE_POSCTL_SLOW
)

var mainModeNames = map[PX4_CUSTOM_MAIN_MODE]string{
	PX4_CUSTOM_MAIN_MODE_MANUAL:           "MANUAL",
	PX4_CUSTOM_MAIN_MODE_ALTCTL:           "ALTCTL",
	PX4_CUSTOM_MAIN_MODE_POSCTL:           "POSCTL",
	PX4_CUSTOM_MAIN_MODE_AUTO:             "AUTO",
	PX4_CUSTOM_MAIN_MODE_ACRO:             "ACRO",
	PX4_CUSTOM_MAIN_MODE_OFFBOARD:         "OFFBOARD",
	PX4_CUSTOM_MAIN_MODE_STABILIZED:       "STABILIZED",
	PX4_CUSTOM_MAIN_MODE_RATTITUDE_LEGACY: "RATTITUDE_LEGACY",
	PX4_CUSTOM_MAIN_MODE_SIMPLE:           "SIMPLE",
	PX4_CUSTOM_MAIN_MODE_TERMINATION:      "TERMINATION",
	PX4_CUSTOM_MAIN_MODE_ALTITUDE_CRUISE:  "ALTITUDE_CRUISE",
}

var autoSubModeNames = map[PX4_CUSTOM_SUB_MODE_AUTO]string{
	PX4_CUSTOM_SUB_MODE_AUTO_READY:         "READY",
	PX4_CUSTOM_SUB_MODE_AUTO_TAKEOFF:       "TAKEOFF",
	PX4_CUSTOM_SUB_MODE_AUTO_LOITER:        "LOITER",
	PX4_CUSTOM_SUB_MODE_AUTO_MISSION:       "MISSION",
	PX4_CUSTOM_SUB_MODE_AUTO_RTL:           "RTL",
	PX4_CUSTOM_SUB_MODE_AUTO_LAND:          "LAND",
	PX4_CUSTOM_SUB_MODE_AUTO_FOLLOW_TARGET: "FOLLOW_TARGET",
	PX4_CUSTOM_SUB_MODE_AUTO_PRECLAND:      "PRECLAND",
	PX4_CUSTOM_SUB_MODE_AUTO_VTOL_TAKEOFF:  "VTOL_TAKEOFF",
	PX4_CUSTOM_SUB_MODE_EXTERNAL1:          "EXTERNAL1",
	PX4_CUSTOM_SUB_MODE_EXTERNAL2:          "EXTERNAL2",
	PX4_CUSTOM_SUB_MODE_EXTERNAL3:          "EXTERNAL3",
	PX4_CUSTOM_SUB_MODE_EXTERNAL4:          "EXTERNAL4",
	PX4_CUSTOM_SUB_MODE_EXTERNAL5:          "EXTERNAL5",
	PX4_CUSTOM_SUB_MODE_EXTERNAL6:          "EXTERNAL6",
	PX4_CUSTOM_SUB_MODE_EXTERNAL7:          "EXTERNAL7",
	PX4_CUSTOM_SUB_MODE_EXTERNAL8:          "EXTERNAL8",
}

var posctlSubModeNames = map[PX4_CUSTOM_SUB_MODE_POSCTL]string{
	PX4_CUSTOM_SUB_MODE_POSCTL_POSCTL: "POSCTL",
	PX4_CUSTOM_SUB_MODE_POSCTL_ORBIT:  "ORBIT",
	PX4_CUSTOM_SUB_MODE_POSCTL_SLOW:   "SLOW",
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
	mainMode := PX4_CUSTOM_MAIN_MODE((customMode >> 16) & 0xFF)
	subMode := uint8((customMode >> 24) & 0xFF)

	result := map[string]interface{}{
		"raw":           fmt.Sprintf("0x%08X", customMode),
		"main_mode":     fmt.Sprintf("0x%02X", uint8(mainMode)),
		"main_mode_str": mainModeNames[mainMode],
		"sub_mode":      fmt.Sprintf("0x%02X", subMode),
	}

	// Add sub mode name based on main mode
	switch mainMode {
	case PX4_CUSTOM_MAIN_MODE_AUTO:
		if name, ok := autoSubModeNames[PX4_CUSTOM_SUB_MODE_AUTO(subMode)]; ok {
			result["sub_mode_str"] = name
		} else {
			result["sub_mode_str"] = fmt.Sprintf("UNKNOWN(%d)", subMode)
		}
	case PX4_CUSTOM_MAIN_MODE_POSCTL:
		if name, ok := posctlSubModeNames[PX4_CUSTOM_SUB_MODE_POSCTL(subMode)]; ok {
			result["sub_mode_str"] = name
		} else {
			result["sub_mode_str"] = fmt.Sprintf("UNKNOWN(%d)", subMode)
		}
	default:
		result["sub_mode_str"] = "N/A"
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
