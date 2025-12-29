package services

import (
	"context"
	"sync"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/internal/mavlink/message_converters"
)

// MessageHandler
// Interface for services that handle a specific message type.
type MessageHandler interface {
	// OnMessage is called when a message of the handler's type is received.
	// The msg parameter will be the protobuf-converted message.
	OnMessage(systemID, componentID uint8, msg interface{})
}

// MessageDispatcher
// Central dispatcher that reads from MAVLink node events and routes messages
// to registered handlers via interface method calls.
type MessageDispatcher struct {
	node *gomavlib.Node

	// Registry: message type name -> handler
	handlers map[string]MessageHandler
	mu       sync.RWMutex

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
		node:     node,
		handlers: make(map[string]MessageHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// RegisterHandler
// Registers a handler for a specific message type.
// The msgTypeName should be the fully qualified type name, e.g., "common.MessageHeartbeat".
func (d *MessageDispatcher) RegisterHandler(msgTypeName string, handler MessageHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[msgTypeName] = handler
}

// Start
// Starts the dispatcher goroutine that reads from node.Events() and routes messages.
// This should be called once when the server starts.
func (d *MessageDispatcher) Start() {
	d.wg.Add(1)
	go d.run()
}

// Stop
// Stops the dispatcher.
func (d *MessageDispatcher) Stop() {
	d.cancel()
	d.wg.Wait()
}

// run
// Main dispatcher loop that reads from node.Events() and routes messages to handlers.
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
					d.dispatchHeartbeat(systemID, componentID, msg)
				case *common.MessageGpsRawInt:
					d.dispatchGpsRawInt(systemID, componentID, msg)
				case *common.MessageSysStatus:
					d.dispatchSysStatus(systemID, componentID, msg)
				}
			}
		}
	}
}

// dispatchHeartbeat
// Converts a HEARTBEAT message to protobuf and dispatches it to the registered handler.
func (d *MessageDispatcher) dispatchHeartbeat(systemID, componentID uint8, msg *common.MessageHeartbeat) {
	pbHeartbeat := message_converters.HeartbeatToProtobuf(msg)

	d.mu.RLock()
	handler := d.handlers["common.MessageHeartbeat"]
	d.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbHeartbeat)
	}
}

// dispatchGpsRawInt
// Converts a GPS_RAW_INT message to protobuf and dispatches it to the registered handler.
func (d *MessageDispatcher) dispatchGpsRawInt(systemID, componentID uint8, msg *common.MessageGpsRawInt) {
	pbGpsRawInt := message_converters.GpsRawIntToProtobuf(msg)

	d.mu.RLock()
	handler := d.handlers["common.MessageGpsRawInt"]
	d.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbGpsRawInt)
	}
}

// dispatchSysStatus
// Converts a SYS_STATUS message to protobuf and dispatches it to the registered handler.
func (d *MessageDispatcher) dispatchSysStatus(systemID, componentID uint8, msg *common.MessageSysStatus) {
	pbSysStatus := message_converters.SysStatusToProtobuf(msg)

	d.mu.RLock()
	handler := d.handlers["common.MessageSysStatus"]
	d.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbSysStatus)
	}
}
