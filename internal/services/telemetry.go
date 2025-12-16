package services

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// TelemetryService implements the TelemetryService gRPC service
type TelemetryService struct {
	flightpathconnect.UnimplementedTelemetryServiceHandler
	ctx *ServiceContext
}

// NewTelemetryService creates a new TelemetryService instance
func NewTelemetryService(ctx *ServiceContext) *TelemetryService {
	return &TelemetryService{
		ctx: ctx,
	}
}

// SubscribeRawGps
// Streams GPS_RAW_INT messages from the MAVLink connection.
// Each message includes the raw GPS data with system/component IDs.
func (s *TelemetryService) SubscribeRawGps(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeRawGpsRequest],
	stream *connect.ServerStream[flightpath.SubscribeRawGpsResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Subscribe to GPS_RAW_INT events from the centralized dispatcher
	gpsRawIntChan := s.ctx.Dispatcher.SubscribeGpsRawInt(ctx)

	// Stream GPS_RAW_INT messages to client
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-gpsRawIntChan:
			if !ok {
				// Channel closed, dispatcher might have stopped
				return nil
			}

			// GPS_RAW_INT is already converted to protobuf by the dispatcher
			response := &flightpath.SubscribeRawGpsResponse{
				TimestampMs: time.Now().UnixMilli(),
				SystemId:     uint32(event.SystemID),
				ComponentId:  uint32(event.ComponentID),
				GpsRawInt:    event.GpsRawInt,
			}

			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

