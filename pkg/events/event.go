package events

import (
	"time"
)

// EventType represents the type of event
type EventType string

const (
	EventTypeTicker    EventType = "ticker"
	EventTypeOrderbook EventType = "orderbook"
	EventTypeTrade     EventType = "trade"
	EventTypeOrder     EventType = "order"
	EventTypePosition  EventType = "position"
	EventTypeBalance   EventType = "balance"
	EventTypeError     EventType = "error"
	EventTypeStatus    EventType = "status"
)

// Event represents a generic event
type Event struct {
	Type      EventType
	Exchange  string
	Symbol    string
	Data      interface{}
	Timestamp time.Time
}

// NewEvent creates a new event
func NewEvent(eventType EventType, exchange, symbol string, data interface{}) *Event {
	return &Event{
		Type:      eventType,
		Exchange:  exchange,
		Symbol:    symbol,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// EventHandler defines the interface for event handlers
type EventHandler interface {
	Handle(event *Event) error
}

// EventBus manages event distribution
type EventBus struct {
	handlers map[EventType][]EventHandler
	ch       chan *Event
}

// NewEventBus creates a new event bus
func NewEventBus(bufferSize int) *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
		ch:       make(chan *Event, bufferSize),
	}
}

// Subscribe subscribes a handler to an event type
func (b *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	if _, ok := b.handlers[eventType]; !ok {
		b.handlers[eventType] = []EventHandler{}
	}
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish publishes an event to the bus
func (b *EventBus) Publish(event *Event) {
	select {
	case b.ch <- event:
	default:
		// Channel full, drop event or log warning
	}
}

// Start starts the event bus processing
func (b *EventBus) Start() {
	go func() {
		for event := range b.ch {
			if handlers, ok := b.handlers[event.Type]; ok {
				for _, handler := range handlers {
					go handler.Handle(event)
				}
			}
		}
	}()
}

// Stop stops the event bus
func (b *EventBus) Stop() {
	close(b.ch)
}
