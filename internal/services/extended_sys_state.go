package services

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// ExtendedSysStateService implements the ExtendedSysStateService gRPC service
// and manages distribution of EXTENDED_SYS_STATE messages to gRPC subscribers.
type ExtendedSysStateService struct {
	flightpathconnect.UnimplementedExtendedSysStateServiceHandler
	ctx *ServiceContext

	// Stream management
	streams []*connect.ServerStream[flightpath.SubscribeExtendedSysStateResponse]
	mu      sync.RWMutex
}

// NewExtendedSysStateService creates a new ExtendedSysStateService instance
// and registers it with the message dispatcher.
func NewExtendedSysStateService(ctx *ServiceContext) *ExtendedSysStateService {
	service := &ExtendedSysStateService{
		ctx:     ctx,
		streams: make([]*connect.ServerStream[flightpath.SubscribeExtendedSysStateResponse], 0),
	}

	// Register handler with dispatcher
	if ctx.Dispatcher != nil {
		ctx.Dispatcher.RegisterHandler("common.MessageExtendedSysState", service)
	}

	return service
}

// SubscribeExtendedSysState
// Streams EXTENDED_SYS_STATE messages from the MAVLink connection.
// Each message includes the VTOL state and landed state information.
func (s *ExtendedSysStateService) SubscribeExtendedSysState(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeExtendedSysStateRequest],
	stream *connect.ServerStream[flightpath.SubscribeExtendedSysStateResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Add stream to subscribers
	s.mu.Lock()
	s.streams = append(s.streams, stream)
	s.mu.Unlock()

	// Remove stream when context is cancelled
	go func() {
		<-ctx.Done()
		s.removeStream(stream)
	}()

	// Block until context is cancelled (stream closed)
	<-ctx.Done()
	return ctx.Err()
}

// OnMessage
// Called by the dispatcher to distribute EXTENDED_SYS_STATE messages to all registered streams.
func (s *ExtendedSysStateService) OnMessage(systemID, componentID uint8, msg interface{}) {
	extendedSysStateMsg := msg.(*flightpath.ExtendedSysState)

	response := &flightpath.SubscribeExtendedSysStateResponse{
		TimestampMs:      time.Now().UnixMilli(),
		SystemId:         uint32(systemID),
		ComponentId:      uint32(componentID),
		ExtendedSysState: extendedSysStateMsg,
	}

	s.mu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeExtendedSysStateResponse], len(s.streams))
	copy(streams, s.streams)
	s.mu.RUnlock()

	// Send to all streams, removing dead ones
	for _, stream := range streams {
		if err := stream.Send(response); err != nil {
			// Stream is dead, remove it
			s.removeStream(stream)
		}
	}
}

// removeStream
// Removes a stream from the subscribers list.
func (s *ExtendedSysStateService) removeStream(stream *connect.ServerStream[flightpath.SubscribeExtendedSysStateResponse]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, st := range s.streams {
		if st == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			return
		}
	}
}
