package event

import (
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if handlers, ok := ed.handlers[eventName]; ok {
		for _, h := range handlers {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}

	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Dispatch(event EventInterface) {
	if handlers, ok := ed.handlers[event.Name()]; ok {
		wg := &sync.WaitGroup{}
		for _, h := range handlers {
			wg.Add(1)
			go h.Handle(event, wg)
		}
		wg.Wait()
	}
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) {
	if handlers, ok := ed.handlers[eventName]; ok {
		for i, h := range handlers {
			if h == handler {
				ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
				return
			}
		}
	}
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if handlers, ok := ed.handlers[eventName]; ok {
		for _, h := range handlers {
			if h == handler {
				return true
			}
		}
	}

	return false
}

func (ed *EventDispatcher) Clear() {
	ed.handlers = make(map[string][]EventHandlerInterface)
}
