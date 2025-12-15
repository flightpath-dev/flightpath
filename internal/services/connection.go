package services

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
	"github.com/flightpath-dev/flightpath/internal/mavlink"
)

// ConnectionService implements the ConnectionService gRPC service
type ConnectionService struct {
	flightpathconnect.UnimplementedConnectionServiceHandler
	ctx *ServiceContext
}

// NewConnectionService creates a new ConnectionService instance
func NewConnectionService(ctx *ServiceContext) *ConnectionService {
	return &ConnectionService{
		ctx: ctx,
	}
}

// SubscribeHeartbeat
// Streams HEARTBEAT messages from the MAVLink connection.
// Each message includes the heartbeat data with system/component IDs and enriched mode information.
func (s *ConnectionService) SubscribeHeartbeat(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeHeartbeatRequest],
	stream *connect.ServerStream[flightpath.SubscribeHeartbeatResponse],
) error {
	if s.ctx.Node == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// HeartbeatEvent contains a heartbeat message with its system/component IDs
	type HeartbeatEvent struct {
		Heartbeat   *common.MessageHeartbeat
		SystemID    uint8
		ComponentID uint8
	}

	// Create a channel to receive heartbeat events
	heartbeatChan := make(chan HeartbeatEvent, 10)

	// Start goroutine to listen to MAVLink events
	go func() {
		defer close(heartbeatChan)

		for evt := range s.ctx.Node.Events() {
			// Process only frame events
			if eventFrame, ok := evt.(*gomavlib.EventFrame); ok {
				msg := eventFrame.Message()

				// Filter for HEARTBEAT messages
				if heartbeat, ok := msg.(*common.MessageHeartbeat); ok {
					event := HeartbeatEvent{
						Heartbeat:   heartbeat,
						SystemID:    eventFrame.SystemID(),
						ComponentID: eventFrame.ComponentID(),
					}

					select {
					case heartbeatChan <- event:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	// Stream heartbeat messages to client
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-heartbeatChan:
			if !ok {
				// Channel closed, node might have disconnected
				return nil
			}

			// Convert MAVLink heartbeat to protobuf
			pbHeartbeat := convertHeartbeatToProtobuf(event.Heartbeat)

			response := &flightpath.SubscribeHeartbeatResponse{
				TimestampMs: time.Now().UnixMilli(),
				SystemId:    uint32(event.SystemID),
				ComponentId: uint32(event.ComponentID),
				Heartbeat:   pbHeartbeat,
			}

			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

// convertHeartbeatToProtobuf
// Converts a MAVLink HEARTBEAT message to a protobuf Heartbeat message.
func convertHeartbeatToProtobuf(msg *common.MessageHeartbeat) *flightpath.Heartbeat {
	// Convert base_mode bitfield to structured BaseMode
	// MAVLink MAV_MODE_FLAG bit positions (from MAVLink spec):
	// Bit 7 (128): MAV_MODE_FLAG_SAFETY_ARMED
	// Bit 6 (64):  MAV_MODE_FLAG_MANUAL_INPUT_ENABLED
	// Bit 5 (32):  MAV_MODE_FLAG_HIL_ENABLED
	// Bit 4 (16):  MAV_MODE_FLAG_STABILIZE_ENABLED
	// Bit 3 (8):   MAV_MODE_FLAG_GUIDED_ENABLED
	// Bit 2 (4):   MAV_MODE_FLAG_AUTO_ENABLED
	// Bit 1 (2):   MAV_MODE_FLAG_TEST_ENABLED
	// Bit 0 (1):   MAV_MODE_FLAG_CUSTOM_MODE_ENABLED
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

	baseMode := &flightpath.BaseMode{
		CustomModeEnabled:  (msg.BaseMode & MAV_MODE_FLAG_CUSTOM_MODE_ENABLED) != 0,
		TestEnabled:        (msg.BaseMode & MAV_MODE_FLAG_TEST_ENABLED) != 0,
		AutoEnabled:        (msg.BaseMode & MAV_MODE_FLAG_AUTO_ENABLED) != 0,
		GuidedEnabled:      (msg.BaseMode & MAV_MODE_FLAG_GUIDED_ENABLED) != 0,
		StabilizeEnabled:   (msg.BaseMode & MAV_MODE_FLAG_STABILIZE_ENABLED) != 0,
		HilEnabled:         (msg.BaseMode & MAV_MODE_FLAG_HIL_ENABLED) != 0,
		ManualInputEnabled: (msg.BaseMode & MAV_MODE_FLAG_MANUAL_INPUT_ENABLED) != 0,
		SafetyArmed:        (msg.BaseMode & MAV_MODE_FLAG_SAFETY_ARMED) != 0,
	}

	// Convert custom_mode (decode if PX4)
	var customMode *flightpath.CustomMode
	if msg.Autopilot == common.MAV_AUTOPILOT_PX4 {
		// Extract main_mode and sub_mode from PX4 custom_mode uint32
		px4MainMode := uint8((msg.CustomMode >> 16) & 0xFF)
		px4SubMode := uint8((msg.CustomMode >> 24) & 0xFF)

		mainMode := mavlink.Px4MainModeToProtobuf(px4MainMode)
		var subMode flightpath.SubMode
		switch mainMode {
		case flightpath.MainMode_MAIN_MODE_AUTO:
			subMode = mavlink.Px4AutoSubModeToProtobuf(px4SubMode)
		case flightpath.MainMode_MAIN_MODE_POSCTL:
			subMode = mavlink.Px4PosctlSubModeToProtobuf(px4SubMode)
		default:
			subMode = flightpath.SubMode_SUB_MODE_UNSPECIFIED
		}

		customMode = &flightpath.CustomMode{
			MainMode: mainMode,
			SubMode:  subMode,
		}
	} else {
		// For non-PX4 autopilots, set to unspecified
		customMode = &flightpath.CustomMode{
			MainMode: flightpath.MainMode_MAIN_MODE_UNSPECIFIED,
			SubMode:  flightpath.SubMode_SUB_MODE_UNSPECIFIED,
		}
	}

	return &flightpath.Heartbeat{
		Type:           convertMavTypeToProtobuf(msg.Type),
		Autopilot:      convertMavAutopilotToProtobuf(msg.Autopilot),
		BaseMode:       baseMode,
		CustomMode:     customMode,
		SystemStatus:   convertMavStateToProtobuf(msg.SystemStatus),
		MavlinkVersion: uint32(msg.MavlinkVersion),
	}
}

// convertMavTypeToProtobuf
// Converts MAVLink MAV_TYPE to protobuf MavType enum.
// Note: Direct cast assumes MAVLink enum values match protobuf enum values.
// If values don't match, explicit mapping may be needed.
func convertMavTypeToProtobuf(mavType common.MAV_TYPE) flightpath.MavType {
	// Direct mapping: MAVLink enum values should match protobuf enum values
	// MAV_TYPE_QUADROTOR = 3 in both MAVLink and protobuf
	return flightpath.MavType(mavType)
}

// convertMavAutopilotToProtobuf
// Converts MAVLink MAV_AUTOPILOT to protobuf MavAutopilot enum.
// Note: Direct cast assumes MAVLink enum values match protobuf enum values.
// If values don't match, explicit mapping may be needed.
func convertMavAutopilotToProtobuf(autopilot common.MAV_AUTOPILOT) flightpath.MavAutopilot {
	// Direct mapping: MAVLink enum values should match protobuf enum values
	// MAV_AUTOPILOT_PX4 = 13 in both MAVLink and protobuf
	return flightpath.MavAutopilot(autopilot)
}

// convertMavStateToProtobuf
// Converts MAVLink MAV_STATE to protobuf MavState enum.
// Note: MAV_STATE_UNINIT (0) maps to MAV_STATE_UNSPECIFIED (0) in protobuf
func convertMavStateToProtobuf(state common.MAV_STATE) flightpath.MavState {
	// Direct mapping: MAVLink enum values match protobuf enum values
	// UNINIT (0) = UNSPECIFIED (0)
	return flightpath.MavState(state)
}
