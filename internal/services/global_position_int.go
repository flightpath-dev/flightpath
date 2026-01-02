package services

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// GlobalPositionIntService implements the GlobalPositionIntService gRPC service
// and manages distribution of GLOBAL_POSITION_INT messages to gRPC subscribers.
type GlobalPositionIntService struct {
	flightpathconnect.UnimplementedGlobalPositionIntServiceHandler
	ctx *ServiceContext

	// Stream management
	streams []*connect.ServerStream[flightpath.SubscribeGlobalPositionIntResponse]
	mu      sync.RWMutex
}

// NewGlobalPositionIntService creates a new GlobalPositionIntService instance
// and registers it with the message dispatcher.
func NewGlobalPositionIntService(ctx *ServiceContext) *GlobalPositionIntService {
	service := &GlobalPositionIntService{
		ctx:     ctx,
		streams: make([]*connect.ServerStream[flightpath.SubscribeGlobalPositionIntResponse], 0),
	}

	// Register handler with dispatcher
	if ctx.Dispatcher != nil {
		ctx.Dispatcher.RegisterHandler("common.MessageGlobalPositionInt", service)
	}

	return service
}

// SubscribeGlobalPositionInt
// Streams GLOBAL_POSITION_INT messages from the MAVLink connection.
// Each message includes the filtered global position with latitude, longitude, altitude, velocity, and heading.
func (s *GlobalPositionIntService) SubscribeGlobalPositionInt(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeGlobalPositionIntRequest],
	stream *connect.ServerStream[flightpath.SubscribeGlobalPositionIntResponse],
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
// Called by the dispatcher to distribute GLOBAL_POSITION_INT messages to all registered streams.
func (s *GlobalPositionIntService) OnMessage(systemID, componentID uint8, msg interface{}) {
	globalPositionIntMsg := msg.(*flightpath.GlobalPositionInt)

	response := &flightpath.SubscribeGlobalPositionIntResponse{
		TimestampMs:       time.Now().UnixMilli(),
		SystemId:          uint32(systemID),
		ComponentId:       uint32(componentID),
		GlobalPositionInt: globalPositionIntMsg,
	}

	s.mu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeGlobalPositionIntResponse], len(s.streams))
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
func (s *GlobalPositionIntService) removeStream(stream *connect.ServerStream[flightpath.SubscribeGlobalPositionIntResponse]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, st := range s.streams {
		if st == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			return
		}
	}
}
