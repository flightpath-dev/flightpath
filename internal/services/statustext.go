package services

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// StatusTextService implements the StatusTextService gRPC service
// and manages distribution of STATUSTEXT messages to gRPC subscribers.
type StatusTextService struct {
	flightpathconnect.UnimplementedStatusTextServiceHandler
	ctx *ServiceContext

	// Stream management
	streams []*connect.ServerStream[flightpath.SubscribeStatusTextResponse]
	mu      sync.RWMutex
}

// NewStatusTextService creates a new StatusTextService instance
// and registers it with the message dispatcher.
func NewStatusTextService(ctx *ServiceContext) *StatusTextService {
	service := &StatusTextService{
		ctx:     ctx,
		streams: make([]*connect.ServerStream[flightpath.SubscribeStatusTextResponse], 0),
	}

	// Register handler with dispatcher
	if ctx.Dispatcher != nil {
		ctx.Dispatcher.RegisterHandler("common.MessageStatustext", service)
	}

	return service
}

// SubscribeStatusText
// Streams STATUSTEXT messages from the MAVLink connection.
// Each message includes the severity level and status text.
func (s *StatusTextService) SubscribeStatusText(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeStatusTextRequest],
	stream *connect.ServerStream[flightpath.SubscribeStatusTextResponse],
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
// Called by the dispatcher to distribute STATUSTEXT messages to all registered streams.
func (s *StatusTextService) OnMessage(systemID, componentID uint8, msg interface{}) {
	statusTextMsg := msg.(*flightpath.StatusText)

	response := &flightpath.SubscribeStatusTextResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		StatusText:  statusTextMsg,
	}

	s.mu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeStatusTextResponse], len(s.streams))
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
func (s *StatusTextService) removeStream(stream *connect.ServerStream[flightpath.SubscribeStatusTextResponse]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, st := range s.streams {
		if st == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			return
		}
	}
}
