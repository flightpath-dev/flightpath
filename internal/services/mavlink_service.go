package services

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// MAVLinkService implements the gRPC service to distribute all MAVLink messages
// to gRPC subscribers.
type MAVLinkService struct {
	flightpathconnect.UnimplementedMAVLinkServiceHandler
	ctx *ServiceContext

	// Stream management
	messagesStreams []*messagesStream
	messagesMu      sync.RWMutex
}

// messagesStream represents a MAVLink message stream with optional filtering
type messagesStream struct {
	stream       *connect.ServerStream[flightpath.SubscribeMessagesResponse]
	messageTypes map[flightpath.MavlinkMessageType]bool // Empty map means all types
}

// NewMAVLinkService creates a new MAVLinkService instance
// and registers all message handlers with the message dispatcher.
func NewMAVLinkService(ctx *ServiceContext) *MAVLinkService {
	service := &MAVLinkService{
		ctx:             ctx,
		messagesStreams: make([]*messagesStream, 0),
	}

	// Register all handlers with dispatcher
	if ctx.Dispatcher != nil {
		ctx.Dispatcher.RegisterHandler("common.MessageHeartbeat", service)
		ctx.Dispatcher.RegisterHandler("common.MessageGpsRawInt", service)
		ctx.Dispatcher.RegisterHandler("common.MessageSysStatus", service)
		ctx.Dispatcher.RegisterHandler("common.MessageExtendedSysState", service)
		ctx.Dispatcher.RegisterHandler("common.MessageStatustext", service)
		ctx.Dispatcher.RegisterHandler("common.MessageRadioStatus", service)
		ctx.Dispatcher.RegisterHandler("common.MessageGlobalPositionInt", service)
		ctx.Dispatcher.RegisterHandler("common.MessageVfrHud", service)
	}

	return service
}

// SubscribeMessages
// Streams all MAVLink messages (or filtered subset) from the MAVLink connection.
func (s *MAVLinkService) SubscribeMessages(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeMessagesRequest],
	stream *connect.ServerStream[flightpath.SubscribeMessagesResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Build filter map from request
	messageTypes := make(map[flightpath.MavlinkMessageType]bool)
	if len(req.Msg.MessageTypes) > 0 {
		for _, msgType := range req.Msg.MessageTypes {
			messageTypes[msgType] = true
		}
	}

	// Create stream wrapper
	ms := &messagesStream{
		stream:       stream,
		messageTypes: messageTypes,
	}

	// Add stream to subscribers
	s.messagesMu.Lock()
	s.messagesStreams = append(s.messagesStreams, ms)
	s.messagesMu.Unlock()

	// Remove stream when context is cancelled
	go func() {
		<-ctx.Done()
		s.removeMessagesStream(ms)
	}()

	// Block until context is cancelled (stream closed)
	<-ctx.Done()
	return ctx.Err()
}

// OnMessage
// Called by the dispatcher to distribute messages to all registered streams.
func (s *MAVLinkService) OnMessage(systemID, componentID uint8, msg interface{}) {
	switch m := msg.(type) {
	case *flightpath.Heartbeat:
		s.distributeMessage(systemID, componentID, flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_HEARTBEAT, &flightpath.SubscribeMessagesResponse_Heartbeat{Heartbeat: m})
	case *flightpath.GpsRawInt:
		s.distributeMessage(systemID, componentID, flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_GPS_RAW_INT, &flightpath.SubscribeMessagesResponse_GpsRawInt{GpsRawInt: m})
	case *flightpath.SysStatus:
		s.distributeMessage(systemID, componentID, flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_SYS_STATUS, &flightpath.SubscribeMessagesResponse_SysStatus{SysStatus: m})
	case *flightpath.ExtendedSysState:
		s.distributeMessage(systemID, componentID, flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_EXTENDED_SYS_STATE, &flightpath.SubscribeMessagesResponse_ExtendedSysState{ExtendedSysState: m})
	case *flightpath.StatusText:
		s.distributeMessage(systemID, componentID, flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_STATUSTEXT, &flightpath.SubscribeMessagesResponse_StatusText{StatusText: m})
	case *flightpath.RadioStatus:
		s.distributeMessage(systemID, componentID, flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_RADIO_STATUS, &flightpath.SubscribeMessagesResponse_RadioStatus{RadioStatus: m})
	case *flightpath.GlobalPositionInt:
		s.distributeMessage(systemID, componentID, flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_GLOBAL_POSITION_INT, &flightpath.SubscribeMessagesResponse_GlobalPositionInt{GlobalPositionInt: m})
	case *flightpath.VfrHud:
		s.distributeMessage(systemID, componentID, flightpath.MavlinkMessageType_MAVLINK_MESSAGE_TYPE_VFR_HUD, &flightpath.SubscribeMessagesResponse_VfrHud{VfrHud: m})
	}
}

// createSubscribeMessagesResponse
// Creates a SubscribeMessagesResponse with the given message data.
// The type switch is necessary because the Message field requires a concrete type
// that implements the unexported isSubscribeMessagesResponse_Message interface.
func createSubscribeMessagesResponse(
	systemID, componentID uint8,
	messageType flightpath.MavlinkMessageType,
	messageData interface{},
) *flightpath.SubscribeMessagesResponse {
	base := &flightpath.SubscribeMessagesResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		MessageType: messageType,
	}

	switch m := messageData.(type) {
	case *flightpath.SubscribeMessagesResponse_Heartbeat:
		base.Message = m
	case *flightpath.SubscribeMessagesResponse_SysStatus:
		base.Message = m
	case *flightpath.SubscribeMessagesResponse_GpsRawInt:
		base.Message = m
	case *flightpath.SubscribeMessagesResponse_GlobalPositionInt:
		base.Message = m
	case *flightpath.SubscribeMessagesResponse_VfrHud:
		base.Message = m
	case *flightpath.SubscribeMessagesResponse_RadioStatus:
		base.Message = m
	case *flightpath.SubscribeMessagesResponse_ExtendedSysState:
		base.Message = m
	case *flightpath.SubscribeMessagesResponse_StatusText:
		base.Message = m
	default:
		return nil // Unknown message type
	}

	return base
}

// distributeMessage
// Distributes messages to all subscribers, respecting filters.
func (s *MAVLinkService) distributeMessage(
	systemID, componentID uint8,
	messageType flightpath.MavlinkMessageType,
	messageData interface{},
) {
	response := createSubscribeMessagesResponse(systemID, componentID, messageType, messageData)
	if response == nil {
		return // Unknown message type, skip
	}

	s.messagesMu.RLock()
	streams := make([]*messagesStream, len(s.messagesStreams))
	copy(streams, s.messagesStreams)
	s.messagesMu.RUnlock()

	// Send to all streams that match the filter, removing dead ones
	for _, ms := range streams {
		// Check if this message type should be sent to this stream
		if len(ms.messageTypes) > 0 && !ms.messageTypes[messageType] {
			continue // Filtered out
		}

		if err := ms.stream.Send(response); err != nil {
			// Stream is dead, remove it
			s.removeMessagesStream(ms)
		}
	}
}

// removeMessagesStream
// Removes a stream from the subscribers list.
func (s *MAVLinkService) removeMessagesStream(stream *messagesStream) {
	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()

	for i, st := range s.messagesStreams {
		if st == stream {
			s.messagesStreams = append(s.messagesStreams[:i], s.messagesStreams[i+1:]...)
			return
		}
	}
}
