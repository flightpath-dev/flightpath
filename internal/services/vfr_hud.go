package services

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// VfrHudService implements the VfrHudService gRPC service
// and manages distribution of VFR_HUD messages to gRPC subscribers.
type VfrHudService struct {
	flightpathconnect.UnimplementedVfrHudServiceHandler
	ctx *ServiceContext

	// Stream management
	streams []*connect.ServerStream[flightpath.SubscribeVfrHudResponse]
	mu      sync.RWMutex
}

// NewVfrHudService creates a new VfrHudService instance
// and registers it with the message dispatcher.
func NewVfrHudService(ctx *ServiceContext) *VfrHudService {
	service := &VfrHudService{
		ctx:     ctx,
		streams: make([]*connect.ServerStream[flightpath.SubscribeVfrHudResponse], 0),
	}

	// Register handler with dispatcher
	if ctx.Dispatcher != nil {
		ctx.Dispatcher.RegisterHandler("common.MessageVfrHud", service)
	}

	return service
}

// SubscribeVfrHud
// Streams VFR_HUD messages from the MAVLink connection.
// Each message includes key flight metrics typically displayed on a HUD for fixed wing aircraft.
func (s *VfrHudService) SubscribeVfrHud(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeVfrHudRequest],
	stream *connect.ServerStream[flightpath.SubscribeVfrHudResponse],
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
// Called by the dispatcher to distribute VFR_HUD messages to all registered streams.
func (s *VfrHudService) OnMessage(systemID, componentID uint8, msg interface{}) {
	vfrHudMsg := msg.(*flightpath.VfrHud)

	response := &flightpath.SubscribeVfrHudResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		VfrHud:      vfrHudMsg,
	}

	s.mu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeVfrHudResponse], len(s.streams))
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
func (s *VfrHudService) removeStream(stream *connect.ServerStream[flightpath.SubscribeVfrHudResponse]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, st := range s.streams {
		if st == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			return
		}
	}
}

