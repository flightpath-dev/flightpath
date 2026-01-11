package mavlink

import (
	"context"
	"sync"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/internal/mavlink/message_converters"
)

// ------------------------------------------------------------------------------------------------
// MessageHandler
// ------------------------------------------------------------------------------------------------
// An interface for dispatching messages to a service.
// Services will register handlers with the MAVLinkMessageReceiver so that they can be called when
// a message is received.
// See RegisterHandler for more information.
// ------------------------------------------------------------------------------------------------
type MessageHandler interface {
	// OnMessage is called when a message of the handler's type is received.
	// The msg parameter will be the protobuf-converted message.
	OnMessage(systemID, componentID uint8, msg interface{})
}

// ------------------------------------------------------------------------------------------------
// MAVLinkMessageReceiver
// ------------------------------------------------------------------------------------------------
// Central receiver that reads incoming MAVLink messages from the drone and routes them to
// registered handlers.
//
//   - The main app creates it using `NewMAVLinkMessageReceiver()`, passing it an initialized node.
//   - `Node` is a MAVLink concept that is used to communicate with a drone on a configured endpoint
//     (serial, UDP, etc.).
//   - The receiver reads from the node's events and routes messages to the handlers registered
//     by various services.
//
// ------------------------------------------------------------------------------------------------
type MAVLinkMessageReceiver struct {
	node *gomavlib.Node

	// Registry: message type name -> handler
	handlers map[string]MessageHandler
	mu       sync.RWMutex

	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewMAVLinkMessageReceiver
// Creates a new message receiver using the provided node.
func NewMAVLinkMessageReceiver(node *gomavlib.Node) *MAVLinkMessageReceiver {
	ctx, cancel := context.WithCancel(context.Background())
	return &MAVLinkMessageReceiver{
		node:     node,
		handlers: make(map[string]MessageHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// RegisterHandler
// Registers a handler for a specific message type.
// The msgTypeName should be the fully qualified type name, e.g., "common.MessageHeartbeat".
func (r *MAVLinkMessageReceiver) RegisterHandler(msgTypeName string, handler MessageHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[msgTypeName] = handler
}

// Start
// Starts the receiver goroutine that reads from node.Events() and routes messages.
// This should be called once when the server starts.
func (r *MAVLinkMessageReceiver) Start() {
	r.wg.Add(1)
	go r.run()
}

// Stop
// Stops the receiver.
func (r *MAVLinkMessageReceiver) Stop() {
	r.cancel()
	r.wg.Wait()
}

// run
// Main receiver loop that reads from node.Events() and routes messages to handlers.
func (r *MAVLinkMessageReceiver) run() {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		// This is where the receiver reads from the node's events and routes messages to
		// the handlers registered by various services.
		case evt, ok := <-r.node.Events():
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
					r.dispatchHeartbeat(systemID, componentID, msg)
				case *common.MessageGpsRawInt:
					r.dispatchGpsRawInt(systemID, componentID, msg)
				case *common.MessageSysStatus:
					r.dispatchSysStatus(systemID, componentID, msg)
				case *common.MessageExtendedSysState:
					r.dispatchExtendedSysState(systemID, componentID, msg)
				case *common.MessageStatustext:
					r.dispatchStatusText(systemID, componentID, msg)
				case *common.MessageRadioStatus:
					r.dispatchRadioStatus(systemID, componentID, msg)
				case *common.MessageGlobalPositionInt:
					r.dispatchGlobalPositionInt(systemID, componentID, msg)
				case *common.MessageVfrHud:
					r.dispatchVfrHud(systemID, componentID, msg)
				}
			}
		}
	}
}

// dispatchHeartbeat
// Converts a HEARTBEAT message to protobuf and dispatches it to the registered handler.
func (r *MAVLinkMessageReceiver) dispatchHeartbeat(systemID, componentID uint8, msg *common.MessageHeartbeat) {
	pbHeartbeat := message_converters.HeartbeatToProtobuf(msg)

	r.mu.RLock()
	handler := r.handlers["common.MessageHeartbeat"]
	r.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbHeartbeat)
	}
}

// dispatchGpsRawInt
// Converts a GPS_RAW_INT message to protobuf and dispatches it to the registered handler.
func (r *MAVLinkMessageReceiver) dispatchGpsRawInt(systemID, componentID uint8, msg *common.MessageGpsRawInt) {
	pbGpsRawInt := message_converters.GpsRawIntToProtobuf(msg)

	r.mu.RLock()
	handler := r.handlers["common.MessageGpsRawInt"]
	r.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbGpsRawInt)
	}
}

// dispatchSysStatus
// Converts a SYS_STATUS message to protobuf and dispatches it to the registered handler.
func (r *MAVLinkMessageReceiver) dispatchSysStatus(systemID, componentID uint8, msg *common.MessageSysStatus) {
	pbSysStatus := message_converters.SysStatusToProtobuf(msg)

	r.mu.RLock()
	handler := r.handlers["common.MessageSysStatus"]
	r.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbSysStatus)
	}
}

// dispatchExtendedSysState
// Converts an EXTENDED_SYS_STATE message to protobuf and dispatches it to the registered handler.
func (r *MAVLinkMessageReceiver) dispatchExtendedSysState(systemID, componentID uint8, msg *common.MessageExtendedSysState) {
	pbExtendedSysState := message_converters.ExtendedSysStateToProtobuf(msg)

	r.mu.RLock()
	handler := r.handlers["common.MessageExtendedSysState"]
	r.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbExtendedSysState)
	}
}

// dispatchStatusText
// Converts a STATUSTEXT message to protobuf and dispatches it to the registered handler.
func (r *MAVLinkMessageReceiver) dispatchStatusText(systemID, componentID uint8, msg *common.MessageStatustext) {
	pbStatusText := message_converters.StatusTextToProtobuf(msg)

	r.mu.RLock()
	handler := r.handlers["common.MessageStatustext"]
	r.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbStatusText)
	}
}

// dispatchRadioStatus
// Converts a RADIO_STATUS message to protobuf and dispatches it to the registered handler.
func (r *MAVLinkMessageReceiver) dispatchRadioStatus(systemID, componentID uint8, msg *common.MessageRadioStatus) {
	pbRadioStatus := message_converters.RadioStatusToProtobuf(msg)

	r.mu.RLock()
	handler := r.handlers["common.MessageRadioStatus"]
	r.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbRadioStatus)
	}
}

// dispatchGlobalPositionInt
// Converts a GLOBAL_POSITION_INT message to protobuf and dispatches it to the registered handler.
func (r *MAVLinkMessageReceiver) dispatchGlobalPositionInt(systemID, componentID uint8, msg *common.MessageGlobalPositionInt) {
	pbGlobalPositionInt := message_converters.GlobalPositionIntToProtobuf(msg)

	r.mu.RLock()
	handler := r.handlers["common.MessageGlobalPositionInt"]
	r.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbGlobalPositionInt)
	}
}

// dispatchVfrHud
// Converts a VFR_HUD message to protobuf and dispatches it to the registered handler.
func (r *MAVLinkMessageReceiver) dispatchVfrHud(systemID, componentID uint8, msg *common.MessageVfrHud) {
	pbVfrHud := message_converters.VfrHudToProtobuf(msg)

	r.mu.RLock()
	handler := r.handlers["common.MessageVfrHud"]
	r.mu.RUnlock()

	if handler != nil {
		handler.OnMessage(systemID, componentID, pbVfrHud)
	}
}
