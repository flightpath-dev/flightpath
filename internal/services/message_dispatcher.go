package services

import (
	"context"
	"sync"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
	"github.com/flightpath-dev/flightpath/internal/mavlink/message_converters"
)

// HeartbeatEvent contains a converted protobuf heartbeat message with its system/component IDs
type HeartbeatEvent struct {
	SystemID    uint8
	ComponentID uint8
	Heartbeat   *flightpath.Heartbeat
}

// GpsRawIntEvent contains a converted protobuf GPS_RAW_INT message with its system/component IDs
type GpsRawIntEvent struct {
	SystemID    uint8
	ComponentID uint8
	GpsRawInt   *flightpath.GpsRawInt
}

// MessageDispatcher
// Central dispatcher that reads from MAVLink node events and routes messages
// to topic-specific channels. Supports multiple subscribers per message type.
type MessageDispatcher struct {
	node *gomavlib.Node

	// Heartbeat subscribers
	heartbeatSubscribers []chan HeartbeatEvent
	heartbeatMu          sync.RWMutex

	// GPS_RAW_INT subscribers
	gpsRawIntSubscribers []chan GpsRawIntEvent
	gpsRawIntMu          sync.RWMutex

	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewMessageDispatcher
// Creates a new message dispatcher that will start processing events from the node.
func NewMessageDispatcher(node *gomavlib.Node) *MessageDispatcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &MessageDispatcher{
		node:                 node,
		heartbeatSubscribers: make([]chan HeartbeatEvent, 0),
		gpsRawIntSubscribers: make([]chan GpsRawIntEvent, 0),
		ctx:                  ctx,
		cancel:               cancel,
	}
}

// Start
// Starts the dispatcher goroutine that reads from node.Events() and routes messages.
// This should be called once when the server starts.
func (d *MessageDispatcher) Start() {
	d.wg.Add(1)
	go d.run()
}

// Stop
// Stops the dispatcher and closes all subscriber channels.
func (d *MessageDispatcher) Stop() {
	d.cancel()
	d.wg.Wait()

	// Close all subscriber channels
	d.heartbeatMu.Lock()
	for _, ch := range d.heartbeatSubscribers {
		close(ch)
	}
	d.heartbeatSubscribers = nil
	d.heartbeatMu.Unlock()

	d.gpsRawIntMu.Lock()
	for _, ch := range d.gpsRawIntSubscribers {
		close(ch)
	}
	d.gpsRawIntSubscribers = nil
	d.gpsRawIntMu.Unlock()
}

// SubscribeHeartbeat
// Subscribes to heartbeat messages. Returns a channel that will receive heartbeat events.
// The channel will be closed when the dispatcher stops or when UnsubscribeHeartbeat is called.
// The caller should handle context cancellation to unsubscribe.
func (d *MessageDispatcher) SubscribeHeartbeat(ctx context.Context) <-chan HeartbeatEvent {
	ch := make(chan HeartbeatEvent, 10)

	d.heartbeatMu.Lock()
	d.heartbeatSubscribers = append(d.heartbeatSubscribers, ch)
	d.heartbeatMu.Unlock()

	// Unsubscribe when context is cancelled
	go func() {
		<-ctx.Done()
		d.UnsubscribeHeartbeat(ch)
	}()

	return ch
}

// UnsubscribeHeartbeat
// Removes a heartbeat subscriber channel.
func (d *MessageDispatcher) UnsubscribeHeartbeat(ch chan HeartbeatEvent) {
	d.heartbeatMu.Lock()
	defer d.heartbeatMu.Unlock()

	for i, subscriber := range d.heartbeatSubscribers {
		if subscriber == ch {
			// Remove from slice
			d.heartbeatSubscribers = append(d.heartbeatSubscribers[:i], d.heartbeatSubscribers[i+1:]...)
			close(ch)
			return
		}
	}
}

// SubscribeGpsRawInt
// Subscribes to GPS_RAW_INT messages. Returns a channel that will receive GPS_RAW_INT events.
// The channel will be closed when the dispatcher stops or when UnsubscribeGpsRawInt is called.
// The caller should handle context cancellation to unsubscribe.
func (d *MessageDispatcher) SubscribeGpsRawInt(ctx context.Context) <-chan GpsRawIntEvent {
	ch := make(chan GpsRawIntEvent, 10)

	d.gpsRawIntMu.Lock()
	d.gpsRawIntSubscribers = append(d.gpsRawIntSubscribers, ch)
	d.gpsRawIntMu.Unlock()

	// Unsubscribe when context is cancelled
	go func() {
		<-ctx.Done()
		d.UnsubscribeGpsRawInt(ch)
	}()

	return ch
}

// UnsubscribeGpsRawInt
// Removes a GPS_RAW_INT subscriber channel.
func (d *MessageDispatcher) UnsubscribeGpsRawInt(ch chan GpsRawIntEvent) {
	d.gpsRawIntMu.Lock()
	defer d.gpsRawIntMu.Unlock()

	for i, subscriber := range d.gpsRawIntSubscribers {
		if subscriber == ch {
			// Remove from slice
			d.gpsRawIntSubscribers = append(d.gpsRawIntSubscribers[:i], d.gpsRawIntSubscribers[i+1:]...)
			close(ch)
			return
		}
	}
}

// run
// Main dispatcher loop that reads from node.Events() and routes messages to subscribers.
func (d *MessageDispatcher) run() {
	defer d.wg.Done()

	for {
		select {
		case <-d.ctx.Done():
			return
		case evt, ok := <-d.node.Events():
			if !ok {
				// Node events channel closed
				return
			}

			// Process only frame events
			if eventFrame, ok := evt.(*gomavlib.EventFrame); ok {
				systemID := eventFrame.SystemID()
				componentID := eventFrame.ComponentID()
				msg := eventFrame.Message()

				// Route messages based on type
				switch msg := msg.(type) {
				case *common.MessageHeartbeat:
					d.broadcastHeartbeat(systemID, componentID, msg)
				case *common.MessageGpsRawInt:
					d.broadcastGpsRawInt(systemID, componentID, msg)
				}
			}
		}
	}
}

// broadcastHeartbeat
// Converts a HEARTBEAT message to protobuf and broadcasts it to all subscribers.
func (d *MessageDispatcher) broadcastHeartbeat(systemID, componentID uint8, msg *common.MessageHeartbeat) {
	pbHeartbeat := message_converters.HeartbeatToProtobuf(msg)
	event := HeartbeatEvent{
		SystemID:    systemID,
		ComponentID: componentID,
		Heartbeat:   pbHeartbeat,
	}

	d.heartbeatMu.RLock()
	subscribers := make([]chan HeartbeatEvent, len(d.heartbeatSubscribers))
	copy(subscribers, d.heartbeatSubscribers)
	d.heartbeatMu.RUnlock()

	// Send to all subscribers (non-blocking)
	for _, ch := range subscribers {
		select {
		case ch <- event:
		default:
			// Channel full, skip this subscriber to avoid blocking
		}
	}
}

// broadcastGpsRawInt
// Converts a GPS_RAW_INT message to protobuf and broadcasts it to all subscribers.
func (d *MessageDispatcher) broadcastGpsRawInt(systemID, componentID uint8, msg *common.MessageGpsRawInt) {
	pbGpsRawInt := message_converters.GpsRawIntToProtobuf(msg)
	event := GpsRawIntEvent{
		SystemID:    systemID,
		ComponentID: componentID,
		GpsRawInt:   pbGpsRawInt,
	}

	d.gpsRawIntMu.RLock()
	subscribers := make([]chan GpsRawIntEvent, len(d.gpsRawIntSubscribers))
	copy(subscribers, d.gpsRawIntSubscribers)
	d.gpsRawIntMu.RUnlock()

	// Send to all subscribers (non-blocking)
	for _, ch := range subscribers {
		select {
		case ch <- event:
		default:
			// Channel full, skip this subscriber to avoid blocking
		}
	}
}
