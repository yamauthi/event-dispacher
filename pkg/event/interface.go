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
	Clear()
	Dispatch(e EventInterface)
	Has(name string, h EventHandlerInterface) bool
	Register(name string, h EventHandlerInterface) error
	Remove(name string, h EventHandlerInterface)
}
