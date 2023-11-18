package event

import (
	"sync"
	"time"
)

type EventInterface interface {
	Name() string
	Payload() interface{}
	SetPayload(p interface{})
	OccurredAt() time.Time
}

type EventHandlerInterface interface {
	Handle(e EventInterface, wg *sync.WaitGroup)
}

type EventDispatcherInterface interface {
	Register(name string, h EventHandlerInterface) error
	Dispatch(e EventInterface) error
	Remove(name string, h EventHandlerInterface) error
	Has(name string, h EventHandlerInterface) bool
	Clear()
}
