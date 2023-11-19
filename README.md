# Event Dispatcher
Package for registering and handling events inside a go application

Structure:

  * EventInterface
    ------------------------------------
    ```go 
    Name() string
	Payload() interface{}
	SetPayload(payload interface{})
	OccurredAt() time.Time 
    ```
    * Contains event information
  * EventHandlerInterface
    ------------------------------------
    ```go 
    Handle(event EventInterface, waitGroup *sync.WaitGroup) 
    ```
    * Executed when an event is dispatched
  * EventDispatcher
    ------------------------------------
    ```go
    Clear()
    ```
    * Clear all the handlers
    ```go
    Dispatch(event EventInterface)
    ``` 
    * Execute all registered handlers for the Event.Name(), each one on a go routine
    ```go
    Has(eventName string, eventHandler EventHandlerInterface) bool
    ```
    * Check if a handler is registered for an Event Name
    ```go
    Register(eventName string, eventHandler EventHandlerInterface) error
    ```
    * Register a handler for an Event Name
    ```go
    Remove(eventName string, eventHandler EventHandlerInterface)
    ```
    * Remove a handler if it is registered for an Event Name