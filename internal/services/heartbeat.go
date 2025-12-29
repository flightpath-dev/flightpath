package services

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// HeartbeatService implements the HeartbeatService gRPC service
// and manages distribution of heartbeat messages to gRPC subscribers.
type HeartbeatService struct {
	flightpathconnect.UnimplementedHeartbeatServiceHandler
	ctx *ServiceContext

	// Stream management
	streams []*connect.ServerStream[flightpath.SubscribeHeartbeatResponse]
	mu      sync.RWMutex
}

// NewHeartbeatService creates a new HeartbeatService instance
// and registers it with the message dispatcher.
func NewHeartbeatService(ctx *ServiceContext) *HeartbeatService {
	service := &HeartbeatService{
		ctx:     ctx,
		streams: make([]*connect.ServerStream[flightpath.SubscribeHeartbeatResponse], 0),
	}

	// Register handler with dispatcher
	if ctx.Dispatcher != nil {
		ctx.Dispatcher.RegisterHandler("common.MessageHeartbeat", service)
	}

	return service
}

// SubscribeHeartbeat
// Streams HEARTBEAT messages from the MAVLink connection.
// Each message includes the heartbeat data with system/component IDs and enriched mode information.
func (s *HeartbeatService) SubscribeHeartbeat(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeHeartbeatRequest],
	stream *connect.ServerStream[flightpath.SubscribeHeartbeatResponse],
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
// Called by the dispatcher to distribute heartbeat messages to all registered streams.
func (s *HeartbeatService) OnMessage(systemID, componentID uint8, msg interface{}) {
	heartbeatMsg := msg.(*flightpath.Heartbeat)

	response := &flightpath.SubscribeHeartbeatResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:     uint32(systemID),
		ComponentId:  uint32(componentID),
		Heartbeat:    heartbeatMsg,
	}

	s.mu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeHeartbeatResponse], len(s.streams))
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
func (s *HeartbeatService) removeStream(stream *connect.ServerStream[flightpath.SubscribeHeartbeatResponse]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, st := range s.streams {
		if st == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			return
		}
	}
}

