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
	Heartbeat   *flightpath.Heartbeat
	SystemID    uint8
	ComponentID uint8
}

// MessageDispatcher
// Central dispatcher that reads from MAVLink node events and routes messages
// to topic-specific channels. Supports multiple subscribers per message type.
type MessageDispatcher struct {
	node *gomavlib.Node

	// Heartbeat subscribers
	heartbeatSubscribers []chan HeartbeatEvent
	heartbeatMu          sync.RWMutex

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
				msg := eventFrame.Message()

				// Route HEARTBEAT messages
				if heartbeat, ok := msg.(*common.MessageHeartbeat); ok {
					// Convert to protobuf once, then broadcast to all subscribers
					pbHeartbeat := message_converters.HeartbeatToProtobuf(heartbeat)
					d.broadcastHeartbeat(HeartbeatEvent{
						Heartbeat:   pbHeartbeat,
						SystemID:    eventFrame.SystemID(),
						ComponentID: eventFrame.ComponentID(),
					})
				}

				// Future: Add routing for other message types here
				// e.g., ATTITUDE, GLOBAL_POSITION_INT, etc.
			}
		}
	}
}

// broadcastHeartbeat
// Broadcasts a heartbeat event to all subscribers.
func (d *MessageDispatcher) broadcastHeartbeat(event HeartbeatEvent) {
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
