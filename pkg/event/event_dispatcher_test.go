package event

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	name       string
	payload    interface{}
	occurredAt time.Time
}

func (te *TestEvent) Name() string {
	return te.name
}

func (te *TestEvent) Payload() interface{} {
	return te.payload
}

func (te *TestEvent) SetPayload(p interface{}) {
	te.payload = p
}

func (te *TestEvent) OccurredAt() time.Time {
	return te.occurredAt
}

type TestEventHandler struct {
	ID int
}

func (teh *TestEventHandler) Handle(e EventInterface, wg *sync.WaitGroup) {
}

type MockEventHandler struct {
	mock.Mock
}

func (meh *MockEventHandler) Handle(e EventInterface, wg *sync.WaitGroup) {
	meh.Called(e)
	wg.Done()
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()

	suite.handler = TestEventHandler{
		ID: 1,
	}
	suite.handler2 = TestEventHandler{
		ID: 2,
	}
	suite.handler3 = TestEventHandler{
		ID: 3,
	}

	suite.event = TestEvent{
		name:       "Test Event 1",
		payload:    "Payload for TestEvent 1",
		occurredAt: time.Now(),
	}

	suite.event2 = TestEvent{
		name: "Test Event 2",
		payload: struct {
			A int
			B int
		}{A: 1, B: 2},
		occurredAt: time.Now(),
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_SameNameAndHandlerError() {
	err := suite.eventDispatcher.Register(suite.event.name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.name]))

	err = suite.eventDispatcher.Register(suite.event.name, &suite.handler)
	suite.Equal(ErrHandlerAlreadyRegistered, err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.name]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	//Test Event 1
	err := suite.eventDispatcher.Register(suite.event.name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.name]))

	err = suite.eventDispatcher.Register(suite.event.name, &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.name]))

	err = suite.eventDispatcher.Register(suite.event.name, &suite.handler3)
	suite.Nil(err)
	suite.Equal(3, len(suite.eventDispatcher.handlers[suite.event.name]))

	//Test Event 2
	err = suite.eventDispatcher.Register(suite.event2.name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.name]))

	err = suite.eventDispatcher.Register(suite.event2.name, &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event2.name]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	meh := MockEventHandler{}
	meh2 := MockEventHandler{}

	meh.On("Handle", &suite.event)
	meh2.On("Handle", &suite.event)

	suite.eventDispatcher.Register(suite.event.name, &meh)
	suite.eventDispatcher.Register(suite.event.name, &meh2)

	suite.eventDispatcher.Dispatch(&suite.event)

	meh.AssertExpectations(suite.T())
	meh2.AssertExpectations(suite.T())
	meh.AssertNumberOfCalls(suite.T(), "Handle", 1)
	meh2.AssertNumberOfCalls(suite.T(), "Handle", 1)

	suite.eventDispatcher.Dispatch(&suite.event)
	meh.AssertNumberOfCalls(suite.T(), "Handle", 2)
	meh2.AssertNumberOfCalls(suite.T(), "Handle", 2)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	//Test Event 1
	suite.eventDispatcher.Register(suite.event.name, &suite.handler)
	suite.eventDispatcher.Register(suite.event.name, &suite.handler2)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.name]))

	//Test Event 2
	suite.eventDispatcher.Register(suite.event2.name, &suite.handler)
	suite.eventDispatcher.Register(suite.event2.name, &suite.handler3)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event2.name]))

	suite.eventDispatcher.Remove(suite.event.name, &suite.handler)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.name]))
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event.name][0])
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event2.name]))

	suite.eventDispatcher.Remove(suite.event.name, &suite.handler2)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event.name]))

	suite.eventDispatcher.Remove(suite.event2.name, &suite.handler3)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.name]))
	assert.Equal(suite.T(), &suite.handler, suite.eventDispatcher.handlers[suite.event2.name][0])

	suite.eventDispatcher.Remove(suite.event2.name, &suite.handler)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event2.name]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	//Test Event 1
	suite.eventDispatcher.Register(suite.event.name, &suite.handler)
	suite.eventDispatcher.Register(suite.event.name, &suite.handler2)

	//Test Event 2
	suite.eventDispatcher.Register(suite.event2.name, &suite.handler)
	suite.eventDispatcher.Register(suite.event2.name, &suite.handler3)

	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.name, &suite.handler))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.name, &suite.handler2))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event.name, &suite.handler3))

	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event2.name, &suite.handler))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event2.name, &suite.handler3))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event2.name, &suite.handler2))

	suite.eventDispatcher.Remove(suite.event.name, &suite.handler)
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event.name, &suite.handler))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.name, &suite.handler2))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event2.name, &suite.handler))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	//Test Event 1
	suite.eventDispatcher.Register(suite.event.name, &suite.handler)
	suite.eventDispatcher.Register(suite.event.name, &suite.handler2)

	//Test Event 2
	suite.eventDispatcher.Register(suite.event2.name, &suite.handler)
	suite.eventDispatcher.Register(suite.event2.name, &suite.handler3)

	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.name]))
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event2.name]))

	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}
