package services

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// RadioStatusService implements the RadioStatusService gRPC service
// and manages distribution of RADIO_STATUS messages to gRPC subscribers.
type RadioStatusService struct {
	flightpathconnect.UnimplementedRadioStatusServiceHandler
	ctx *ServiceContext

	// Stream management
	streams []*connect.ServerStream[flightpath.SubscribeRadioStatusResponse]
	mu      sync.RWMutex
}

// NewRadioStatusService creates a new RadioStatusService instance
// and registers it with the message dispatcher.
func NewRadioStatusService(ctx *ServiceContext) *RadioStatusService {
	service := &RadioStatusService{
		ctx:     ctx,
		streams: make([]*connect.ServerStream[flightpath.SubscribeRadioStatusResponse], 0),
	}

	// Register handler with dispatcher
	if ctx.Dispatcher != nil {
		ctx.Dispatcher.RegisterHandler("common.MessageRadioStatus", service)
	}

	return service
}

// SubscribeRadioStatus
// Streams RADIO_STATUS messages from the MAVLink connection.
// Each message includes radio signal strength, noise levels, and error counts.
func (s *RadioStatusService) SubscribeRadioStatus(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeRadioStatusRequest],
	stream *connect.ServerStream[flightpath.SubscribeRadioStatusResponse],
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
// Called by the dispatcher to distribute RADIO_STATUS messages to all registered streams.
func (s *RadioStatusService) OnMessage(systemID, componentID uint8, msg interface{}) {
	radioStatusMsg := msg.(*flightpath.RadioStatus)

	response := &flightpath.SubscribeRadioStatusResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		RadioStatus: radioStatusMsg,
	}

	s.mu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeRadioStatusResponse], len(s.streams))
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
func (s *RadioStatusService) removeStream(stream *connect.ServerStream[flightpath.SubscribeRadioStatusResponse]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, st := range s.streams {
		if st == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			return
		}
	}
}

