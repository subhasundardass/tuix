// events/events.go
package events

import (
	"fmt"
	"sync"
)

// EventHandler is a function that handles an event
type EventHandler func(data interface{})

// EventBus manages event subscriptions and publishing
type EventBus struct {
	mu           sync.RWMutex
	handlers     map[string][]EventHandler
	onceHandlers map[string][]EventHandler
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers:     make(map[string][]EventHandler),
		onceHandlers: make(map[string][]EventHandler),
	}
}

// Subscribe registers a handler for an event
func (e *EventBus) Subscribe(event string, handler EventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[event] = append(e.handlers[event], handler)
}

// SubscribeOnce registers a handler that runs only once
func (e *EventBus) SubscribeOnce(event string, handler EventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.onceHandlers[event] = append(e.onceHandlers[event], handler)
}

// Publish sends an event with data to all subscribers
func (e *EventBus) Publish(event string, data interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Call regular handlers
	if handlers, ok := e.handlers[event]; ok {
		for _, handler := range handlers {
			go func(h EventHandler) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("⚠️ Event handler panic: %v\n", r)
					}
				}()
				h(data)
			}(handler)
		}
	}

	// Call once handlers and remove them
	if onceHandlers, ok := e.onceHandlers[event]; ok {
		// We need to unlock and relock to modify the map
		e.mu.RUnlock()
		e.mu.Lock()
		for _, handler := range onceHandlers {
			go func(h EventHandler) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("⚠️ Once handler panic: %v\n", r)
					}
				}()
				h(data)
			}(handler)
		}
		delete(e.onceHandlers, event)
		e.mu.Unlock()
		e.mu.RLock()
	}
}

// Unsubscribe removes all handlers for an event
func (e *EventBus) Unsubscribe(event string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.handlers, event)
	delete(e.onceHandlers, event)
}
