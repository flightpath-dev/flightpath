package services

import (
	"context"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath/flightpathconnect"
)

// MAVLinkService implements the gRPC service to distribute all MAVLink messages
// to gRPC subscribers.
type MAVLinkService struct {
	flightpathconnect.UnimplementedMAVLinkServiceHandler
	ctx *ServiceContext

	// Stream management for each message type
	heartbeatStreams         []*connect.ServerStream[flightpath.SubscribeHeartbeatResponse]
	gpsRawIntStreams         []*connect.ServerStream[flightpath.SubscribeGpsRawIntResponse]
	sysStatusStreams         []*connect.ServerStream[flightpath.SubscribeSysStatusResponse]
	extendedSysStateStreams  []*connect.ServerStream[flightpath.SubscribeExtendedSysStateResponse]
	statusTextStreams        []*connect.ServerStream[flightpath.SubscribeStatusTextResponse]
	radioStatusStreams       []*connect.ServerStream[flightpath.SubscribeRadioStatusResponse]
	globalPositionIntStreams []*connect.ServerStream[flightpath.SubscribeGlobalPositionIntResponse]
	vfrHudStreams            []*connect.ServerStream[flightpath.SubscribeVfrHudResponse]

	// Mutexes for each stream type
	heartbeatMu         sync.RWMutex
	gpsRawIntMu         sync.RWMutex
	sysStatusMu         sync.RWMutex
	extendedSysStateMu  sync.RWMutex
	statusTextMu        sync.RWMutex
	radioStatusMu       sync.RWMutex
	globalPositionIntMu sync.RWMutex
	vfrHudMu            sync.RWMutex
}

// NewMAVLinkService creates a new MAVLinkService instance
// and registers all message handlers with the message dispatcher.
func NewMAVLinkService(ctx *ServiceContext) *MAVLinkService {
	service := &MAVLinkService{
		ctx:                      ctx,
		heartbeatStreams:         make([]*connect.ServerStream[flightpath.SubscribeHeartbeatResponse], 0),
		gpsRawIntStreams:         make([]*connect.ServerStream[flightpath.SubscribeGpsRawIntResponse], 0),
		sysStatusStreams:         make([]*connect.ServerStream[flightpath.SubscribeSysStatusResponse], 0),
		extendedSysStateStreams:  make([]*connect.ServerStream[flightpath.SubscribeExtendedSysStateResponse], 0),
		statusTextStreams:        make([]*connect.ServerStream[flightpath.SubscribeStatusTextResponse], 0),
		radioStatusStreams:       make([]*connect.ServerStream[flightpath.SubscribeRadioStatusResponse], 0),
		globalPositionIntStreams: make([]*connect.ServerStream[flightpath.SubscribeGlobalPositionIntResponse], 0),
		vfrHudStreams:            make([]*connect.ServerStream[flightpath.SubscribeVfrHudResponse], 0),
	}

	// Register all handlers with dispatcher
	if ctx.Dispatcher != nil {
		ctx.Dispatcher.RegisterHandler("common.MessageHeartbeat", service)
		ctx.Dispatcher.RegisterHandler("common.MessageGpsRawInt", service)
		ctx.Dispatcher.RegisterHandler("common.MessageSysStatus", service)
		ctx.Dispatcher.RegisterHandler("common.MessageExtendedSysState", service)
		ctx.Dispatcher.RegisterHandler("common.MessageStatustext", service)
		ctx.Dispatcher.RegisterHandler("common.MessageRadioStatus", service)
		ctx.Dispatcher.RegisterHandler("common.MessageGlobalPositionInt", service)
		ctx.Dispatcher.RegisterHandler("common.MessageVfrHud", service)
	}

	return service
}

// SubscribeHeartbeat
// Streams HEARTBEAT messages from the MAVLink connection.
// Each message includes the heartbeat data with system/component IDs and enriched mode information.
func (s *MAVLinkService) SubscribeHeartbeat(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeHeartbeatRequest],
	stream *connect.ServerStream[flightpath.SubscribeHeartbeatResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Add stream to subscribers
	s.heartbeatMu.Lock()
	s.heartbeatStreams = append(s.heartbeatStreams, stream)
	s.heartbeatMu.Unlock()

	// Remove stream when context is cancelled
	go func() {
		<-ctx.Done()
		s.removeHeartbeatStream(stream)
	}()

	// Block until context is cancelled (stream closed)
	<-ctx.Done()
	return ctx.Err()
}

// SubscribeGpsRawInt
// Streams GPS_RAW_INT messages from the MAVLink connection.
// Each message includes the raw GPS data with system/component IDs.
func (s *MAVLinkService) SubscribeGpsRawInt(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeGpsRawIntRequest],
	stream *connect.ServerStream[flightpath.SubscribeGpsRawIntResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Add stream to subscribers
	s.gpsRawIntMu.Lock()
	s.gpsRawIntStreams = append(s.gpsRawIntStreams, stream)
	s.gpsRawIntMu.Unlock()

	// Remove stream when context is cancelled
	go func() {
		<-ctx.Done()
		s.removeGpsRawIntStream(stream)
	}()

	// Block until context is cancelled (stream closed)
	<-ctx.Done()
	return ctx.Err()
}

// SubscribeSysStatus
// Streams SYS_STATUS messages from the MAVLink connection.
// Each message includes the system status data with sensor information and battery status.
func (s *MAVLinkService) SubscribeSysStatus(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeSysStatusRequest],
	stream *connect.ServerStream[flightpath.SubscribeSysStatusResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Add stream to subscribers
	s.sysStatusMu.Lock()
	s.sysStatusStreams = append(s.sysStatusStreams, stream)
	s.sysStatusMu.Unlock()

	// Remove stream when context is cancelled
	go func() {
		<-ctx.Done()
		s.removeSysStatusStream(stream)
	}()

	// Block until context is cancelled (stream closed)
	<-ctx.Done()
	return ctx.Err()
}

// SubscribeExtendedSysState
// Streams EXTENDED_SYS_STATE messages from the MAVLink connection.
// Each message includes the VTOL state and landed state information.
func (s *MAVLinkService) SubscribeExtendedSysState(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeExtendedSysStateRequest],
	stream *connect.ServerStream[flightpath.SubscribeExtendedSysStateResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Add stream to subscribers
	s.extendedSysStateMu.Lock()
	s.extendedSysStateStreams = append(s.extendedSysStateStreams, stream)
	s.extendedSysStateMu.Unlock()

	// Remove stream when context is cancelled
	go func() {
		<-ctx.Done()
		s.removeExtendedSysStateStream(stream)
	}()

	// Block until context is cancelled (stream closed)
	<-ctx.Done()
	return ctx.Err()
}

// SubscribeStatusText
// Streams STATUSTEXT messages from the MAVLink connection.
// Each message includes the severity level and status text.
func (s *MAVLinkService) SubscribeStatusText(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeStatusTextRequest],
	stream *connect.ServerStream[flightpath.SubscribeStatusTextResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Add stream to subscribers
	s.statusTextMu.Lock()
	s.statusTextStreams = append(s.statusTextStreams, stream)
	s.statusTextMu.Unlock()

	// Remove stream when context is cancelled
	go func() {
		<-ctx.Done()
		s.removeStatusTextStream(stream)
	}()

	// Block until context is cancelled (stream closed)
	<-ctx.Done()
	return ctx.Err()
}

// SubscribeRadioStatus
// Streams RADIO_STATUS messages from the MAVLink connection.
// Each message includes radio signal strength, noise levels, and error counts.
func (s *MAVLinkService) SubscribeRadioStatus(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeRadioStatusRequest],
	stream *connect.ServerStream[flightpath.SubscribeRadioStatusResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Add stream to subscribers
	s.radioStatusMu.Lock()
	s.radioStatusStreams = append(s.radioStatusStreams, stream)
	s.radioStatusMu.Unlock()

	// Remove stream when context is cancelled
	go func() {
		<-ctx.Done()
		s.removeRadioStatusStream(stream)
	}()

	// Block until context is cancelled (stream closed)
	<-ctx.Done()
	return ctx.Err()
}

// SubscribeGlobalPositionInt
// Streams GLOBAL_POSITION_INT messages from the MAVLink connection.
// Each message includes the filtered global position with latitude, longitude, altitude, velocity, and heading.
func (s *MAVLinkService) SubscribeGlobalPositionInt(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeGlobalPositionIntRequest],
	stream *connect.ServerStream[flightpath.SubscribeGlobalPositionIntResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Add stream to subscribers
	s.globalPositionIntMu.Lock()
	s.globalPositionIntStreams = append(s.globalPositionIntStreams, stream)
	s.globalPositionIntMu.Unlock()

	// Remove stream when context is cancelled
	go func() {
		<-ctx.Done()
		s.removeGlobalPositionIntStream(stream)
	}()

	// Block until context is cancelled (stream closed)
	<-ctx.Done()
	return ctx.Err()
}

// SubscribeVfrHud
// Streams VFR_HUD messages from the MAVLink connection.
// Each message includes key flight metrics typically displayed on a HUD for fixed wing aircraft.
func (s *MAVLinkService) SubscribeVfrHud(
	ctx context.Context,
	req *connect.Request[flightpath.SubscribeVfrHudRequest],
	stream *connect.ServerStream[flightpath.SubscribeVfrHudResponse],
) error {
	if s.ctx.Dispatcher == nil {
		return connect.NewError(connect.CodeFailedPrecondition, nil)
	}

	// Add stream to subscribers
	s.vfrHudMu.Lock()
	s.vfrHudStreams = append(s.vfrHudStreams, stream)
	s.vfrHudMu.Unlock()

	// Remove stream when context is cancelled
	go func() {
		<-ctx.Done()
		s.removeVfrHudStream(stream)
	}()

	// Block until context is cancelled (stream closed)
	<-ctx.Done()
	return ctx.Err()
}

// OnMessage
// Called by the dispatcher to distribute messages to all registered streams.
// Routes messages to the appropriate stream type based on the message type.
func (s *MAVLinkService) OnMessage(systemID, componentID uint8, msg interface{}) {
	switch m := msg.(type) {
	case *flightpath.Heartbeat:
		s.distributeHeartbeat(systemID, componentID, m)
	case *flightpath.GpsRawInt:
		s.distributeGpsRawInt(systemID, componentID, m)
	case *flightpath.SysStatus:
		s.distributeSysStatus(systemID, componentID, m)
	case *flightpath.ExtendedSysState:
		s.distributeExtendedSysState(systemID, componentID, m)
	case *flightpath.StatusText:
		s.distributeStatusText(systemID, componentID, m)
	case *flightpath.RadioStatus:
		s.distributeRadioStatus(systemID, componentID, m)
	case *flightpath.GlobalPositionInt:
		s.distributeGlobalPositionInt(systemID, componentID, m)
	case *flightpath.VfrHud:
		s.distributeVfrHud(systemID, componentID, m)
	}
}

// distributeHeartbeat
// Distributes HEARTBEAT messages to all registered streams.
func (s *MAVLinkService) distributeHeartbeat(systemID, componentID uint8, heartbeat *flightpath.Heartbeat) {
	response := &flightpath.SubscribeHeartbeatResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		Heartbeat:   heartbeat,
	}

	s.heartbeatMu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeHeartbeatResponse], len(s.heartbeatStreams))
	copy(streams, s.heartbeatStreams)
	s.heartbeatMu.RUnlock()

	// Send to all streams, removing dead ones
	for _, stream := range streams {
		if err := stream.Send(response); err != nil {
			// Stream is dead, remove it
			s.removeHeartbeatStream(stream)
		}
	}
}

// distributeGpsRawInt
// Distributes GPS_RAW_INT messages to all registered streams.
func (s *MAVLinkService) distributeGpsRawInt(systemID, componentID uint8, gpsRawInt *flightpath.GpsRawInt) {
	response := &flightpath.SubscribeGpsRawIntResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		GpsRawInt:   gpsRawInt,
	}

	s.gpsRawIntMu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeGpsRawIntResponse], len(s.gpsRawIntStreams))
	copy(streams, s.gpsRawIntStreams)
	s.gpsRawIntMu.RUnlock()

	// Send to all streams, removing dead ones
	for _, stream := range streams {
		if err := stream.Send(response); err != nil {
			// Stream is dead, remove it
			s.removeGpsRawIntStream(stream)
		}
	}
}

// distributeSysStatus
// Distributes SYS_STATUS messages to all registered streams.
func (s *MAVLinkService) distributeSysStatus(systemID, componentID uint8, sysStatus *flightpath.SysStatus) {
	response := &flightpath.SubscribeSysStatusResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		SysStatus:   sysStatus,
	}

	s.sysStatusMu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeSysStatusResponse], len(s.sysStatusStreams))
	copy(streams, s.sysStatusStreams)
	s.sysStatusMu.RUnlock()

	// Send to all streams, removing dead ones
	for _, stream := range streams {
		if err := stream.Send(response); err != nil {
			// Stream is dead, remove it
			s.removeSysStatusStream(stream)
		}
	}
}

// distributeExtendedSysState
// Distributes EXTENDED_SYS_STATE messages to all registered streams.
func (s *MAVLinkService) distributeExtendedSysState(systemID, componentID uint8, extendedSysState *flightpath.ExtendedSysState) {
	response := &flightpath.SubscribeExtendedSysStateResponse{
		TimestampMs:      time.Now().UnixMilli(),
		SystemId:         uint32(systemID),
		ComponentId:      uint32(componentID),
		ExtendedSysState: extendedSysState,
	}

	s.extendedSysStateMu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeExtendedSysStateResponse], len(s.extendedSysStateStreams))
	copy(streams, s.extendedSysStateStreams)
	s.extendedSysStateMu.RUnlock()

	// Send to all streams, removing dead ones
	for _, stream := range streams {
		if err := stream.Send(response); err != nil {
			// Stream is dead, remove it
			s.removeExtendedSysStateStream(stream)
		}
	}
}

// distributeStatusText
// Distributes STATUSTEXT messages to all registered streams.
func (s *MAVLinkService) distributeStatusText(systemID, componentID uint8, statusText *flightpath.StatusText) {
	response := &flightpath.SubscribeStatusTextResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		StatusText:  statusText,
	}

	s.statusTextMu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeStatusTextResponse], len(s.statusTextStreams))
	copy(streams, s.statusTextStreams)
	s.statusTextMu.RUnlock()

	// Send to all streams, removing dead ones
	for _, stream := range streams {
		if err := stream.Send(response); err != nil {
			// Stream is dead, remove it
			s.removeStatusTextStream(stream)
		}
	}
}

// distributeRadioStatus
// Distributes RADIO_STATUS messages to all registered streams.
func (s *MAVLinkService) distributeRadioStatus(systemID, componentID uint8, radioStatus *flightpath.RadioStatus) {
	response := &flightpath.SubscribeRadioStatusResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		RadioStatus: radioStatus,
	}

	s.radioStatusMu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeRadioStatusResponse], len(s.radioStatusStreams))
	copy(streams, s.radioStatusStreams)
	s.radioStatusMu.RUnlock()

	// Send to all streams, removing dead ones
	for _, stream := range streams {
		if err := stream.Send(response); err != nil {
			// Stream is dead, remove it
			s.removeRadioStatusStream(stream)
		}
	}
}

// distributeGlobalPositionInt
// Distributes GLOBAL_POSITION_INT messages to all registered streams.
func (s *MAVLinkService) distributeGlobalPositionInt(systemID, componentID uint8, globalPositionInt *flightpath.GlobalPositionInt) {
	response := &flightpath.SubscribeGlobalPositionIntResponse{
		TimestampMs:       time.Now().UnixMilli(),
		SystemId:          uint32(systemID),
		ComponentId:       uint32(componentID),
		GlobalPositionInt: globalPositionInt,
	}

	s.globalPositionIntMu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeGlobalPositionIntResponse], len(s.globalPositionIntStreams))
	copy(streams, s.globalPositionIntStreams)
	s.globalPositionIntMu.RUnlock()

	// Send to all streams, removing dead ones
	for _, stream := range streams {
		if err := stream.Send(response); err != nil {
			// Stream is dead, remove it
			s.removeGlobalPositionIntStream(stream)
		}
	}
}

// distributeVfrHud
// Distributes VFR_HUD messages to all registered streams.
func (s *MAVLinkService) distributeVfrHud(systemID, componentID uint8, vfrHud *flightpath.VfrHud) {
	response := &flightpath.SubscribeVfrHudResponse{
		TimestampMs: time.Now().UnixMilli(),
		SystemId:    uint32(systemID),
		ComponentId: uint32(componentID),
		VfrHud:      vfrHud,
	}

	s.vfrHudMu.RLock()
	streams := make([]*connect.ServerStream[flightpath.SubscribeVfrHudResponse], len(s.vfrHudStreams))
	copy(streams, s.vfrHudStreams)
	s.vfrHudMu.RUnlock()

	// Send to all streams, removing dead ones
	for _, stream := range streams {
		if err := stream.Send(response); err != nil {
			// Stream is dead, remove it
			s.removeVfrHudStream(stream)
		}
	}
}

// removeHeartbeatStream
// Removes a stream from the heartbeat subscribers list.
func (s *MAVLinkService) removeHeartbeatStream(stream *connect.ServerStream[flightpath.SubscribeHeartbeatResponse]) {
	s.heartbeatMu.Lock()
	defer s.heartbeatMu.Unlock()

	for i, st := range s.heartbeatStreams {
		if st == stream {
			s.heartbeatStreams = append(s.heartbeatStreams[:i], s.heartbeatStreams[i+1:]...)
			return
		}
	}
}

// removeGpsRawIntStream
// Removes a stream from the GPS_RAW_INT subscribers list.
func (s *MAVLinkService) removeGpsRawIntStream(stream *connect.ServerStream[flightpath.SubscribeGpsRawIntResponse]) {
	s.gpsRawIntMu.Lock()
	defer s.gpsRawIntMu.Unlock()

	for i, st := range s.gpsRawIntStreams {
		if st == stream {
			s.gpsRawIntStreams = append(s.gpsRawIntStreams[:i], s.gpsRawIntStreams[i+1:]...)
			return
		}
	}
}

// removeSysStatusStream
// Removes a stream from the SYS_STATUS subscribers list.
func (s *MAVLinkService) removeSysStatusStream(stream *connect.ServerStream[flightpath.SubscribeSysStatusResponse]) {
	s.sysStatusMu.Lock()
	defer s.sysStatusMu.Unlock()

	for i, st := range s.sysStatusStreams {
		if st == stream {
			s.sysStatusStreams = append(s.sysStatusStreams[:i], s.sysStatusStreams[i+1:]...)
			return
		}
	}
}

// removeExtendedSysStateStream
// Removes a stream from the EXTENDED_SYS_STATE subscribers list.
func (s *MAVLinkService) removeExtendedSysStateStream(stream *connect.ServerStream[flightpath.SubscribeExtendedSysStateResponse]) {
	s.extendedSysStateMu.Lock()
	defer s.extendedSysStateMu.Unlock()

	for i, st := range s.extendedSysStateStreams {
		if st == stream {
			s.extendedSysStateStreams = append(s.extendedSysStateStreams[:i], s.extendedSysStateStreams[i+1:]...)
			return
		}
	}
}

// removeStatusTextStream
// Removes a stream from the STATUSTEXT subscribers list.
func (s *MAVLinkService) removeStatusTextStream(stream *connect.ServerStream[flightpath.SubscribeStatusTextResponse]) {
	s.statusTextMu.Lock()
	defer s.statusTextMu.Unlock()

	for i, st := range s.statusTextStreams {
		if st == stream {
			s.statusTextStreams = append(s.statusTextStreams[:i], s.statusTextStreams[i+1:]...)
			return
		}
	}
}

// removeRadioStatusStream
// Removes a stream from the RADIO_STATUS subscribers list.
func (s *MAVLinkService) removeRadioStatusStream(stream *connect.ServerStream[flightpath.SubscribeRadioStatusResponse]) {
	s.radioStatusMu.Lock()
	defer s.radioStatusMu.Unlock()

	for i, st := range s.radioStatusStreams {
		if st == stream {
			s.radioStatusStreams = append(s.radioStatusStreams[:i], s.radioStatusStreams[i+1:]...)
			return
		}
	}
}

// removeGlobalPositionIntStream
// Removes a stream from the GLOBAL_POSITION_INT subscribers list.
func (s *MAVLinkService) removeGlobalPositionIntStream(stream *connect.ServerStream[flightpath.SubscribeGlobalPositionIntResponse]) {
	s.globalPositionIntMu.Lock()
	defer s.globalPositionIntMu.Unlock()

	for i, st := range s.globalPositionIntStreams {
		if st == stream {
			s.globalPositionIntStreams = append(s.globalPositionIntStreams[:i], s.globalPositionIntStreams[i+1:]...)
			return
		}
	}
}

// removeVfrHudStream
// Removes a stream from the VFR_HUD subscribers list.
func (s *MAVLinkService) removeVfrHudStream(stream *connect.ServerStream[flightpath.SubscribeVfrHudResponse]) {
	s.vfrHudMu.Lock()
	defer s.vfrHudMu.Unlock()

	for i, st := range s.vfrHudStreams {
		if st == stream {
			s.vfrHudStreams = append(s.vfrHudStreams[:i], s.vfrHudStreams[i+1:]...)
			return
		}
	}
}
