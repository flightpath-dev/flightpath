package services

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
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
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Subscribe to heartbeat events from the centralized dispatcher
	heartbeatChan := s.ctx.Dispatcher.SubscribeHeartbeat(ctx)

	// Stream heartbeat messages to client
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-heartbeatChan:
			if !ok {
				// Channel closed, dispatcher might have stopped
				return nil
			}

			// Heartbeat is already converted to protobuf by the dispatcher
			response := &flightpath.SubscribeHeartbeatResponse{
				TimestampMs: time.Now().UnixMilli(),
				SystemId:    uint32(event.SystemID),
				ComponentId: uint32(event.ComponentID),
				Heartbeat:   event.Heartbeat,
			}

			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}
