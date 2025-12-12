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

// StreamHeartbeats streams heartbeat messages every 2 seconds
func (s *ConnectionService) StreamHeartbeats(
	ctx context.Context,
	req *connect.Request[flightpath.StreamHeartbeatsRequest],
	stream *connect.ServerStream[flightpath.StreamHeartbeatsResponse],
) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			response := &flightpath.StreamHeartbeatsResponse{
				TimestampMs: time.Now().UnixMilli(),
			}
			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}
