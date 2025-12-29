package services

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// SysStatusService implements the SysStatusService gRPC service
// and manages distribution of SYS_STATUS messages to gRPC subscribers.
type SysStatusService struct {
	flightpathconnect.UnimplementedSysStatusServiceHandler
	ctx *ServiceContext

	// Stream management
	streams []*connect.ServerStream[flightpath.SubscribeSysStatusResponse]
	mu      sync.RWMutex
}

// NewSysStatusService creates a new SysStatusService instance
// and registers it with the message dispatcher.
func NewSysStatusService(ctx *ServiceContext) *SysStatusService {
	service := &SysStatusService{
		ctx:     ctx,
		streams: make([]*connect.ServerStream[flightpath.SubscribeSysStatusResponse], 0),
	}

	// Register handler with dispatcher
	if ctx.Dispatcher != nil {
		ctx.Dispatcher.RegisterHandler("common.MessageSysStatus", service)
	}

	return service
}

// SubscribeSysStatus
// Streams SYS_STATUS messages from the MAVLink connection.
// Each message includes the system status data with sensor information and battery status.
func (s *SysStatusService) SubscribeSysStatus(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeSysStatusRequest],
	stream *connect.ServerStream[flightpath.SubscribeSysStatusResponse],
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
// Called by the dispatcher to distribute SYS_STATUS messages to all registered streams.
func (s *SysStatusService) OnMessage(systemID, componentID uint8, msg interface{}) {
	sysStatusMsg := msg.(*flightpath.SysStatus)

	response := &flightpath.SubscribeSysStatusResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		SysStatus:   sysStatusMsg,
	}

	s.mu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeSysStatusResponse], len(s.streams))
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
func (s *SysStatusService) removeStream(stream *connect.ServerStream[flightpath.SubscribeSysStatusResponse]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, st := range s.streams {
		if st == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			return
		}
	}
}
