package services

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// GpsRawIntService implements the GpsRawIntService gRPC service
// and manages distribution of GPS_RAW_INT messages to gRPC subscribers.
type GpsRawIntService struct {
	flightpathconnect.UnimplementedGpsRawIntServiceHandler
	ctx *ServiceContext

	// Stream management
	streams []*connect.ServerStream[flightpath.SubscribeRawGpsResponse]
	mu      sync.RWMutex
}

// NewGpsRawIntService creates a new GpsRawIntService instance
// and registers it with the message dispatcher.
func NewGpsRawIntService(ctx *ServiceContext) *GpsRawIntService {
	service := &GpsRawIntService{
		ctx:     ctx,
		streams: make([]*connect.ServerStream[flightpath.SubscribeRawGpsResponse], 0),
	}

	// Register handler with dispatcher
	if ctx.Dispatcher != nil {
		ctx.Dispatcher.RegisterHandler("common.MessageGpsRawInt", service)
	}

	return service
}

// SubscribeRawGps
// Streams GPS_RAW_INT messages from the MAVLink connection.
// Each message includes the raw GPS data with system/component IDs.
func (s *GpsRawIntService) SubscribeRawGps(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeRawGpsRequest],
	stream *connect.ServerStream[flightpath.SubscribeRawGpsResponse],
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
// Called by the dispatcher to distribute GPS_RAW_INT messages to all registered streams.
func (s *GpsRawIntService) OnMessage(systemID, componentID uint8, msg interface{}) {
	gpsRawIntMsg := msg.(*flightpath.GpsRawInt)

	response := &flightpath.SubscribeRawGpsResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:     uint32(systemID),
		ComponentId:  uint32(componentID),
		GpsRawInt:    gpsRawIntMsg,
	}

	s.mu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeRawGpsResponse], len(s.streams))
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
func (s *GpsRawIntService) removeStream(stream *connect.ServerStream[flightpath.SubscribeRawGpsResponse]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, st := range s.streams {
		if st == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			return
		}
	}
}

