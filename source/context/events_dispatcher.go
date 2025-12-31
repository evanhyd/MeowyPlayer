package context

type EventsDispatcher struct {
	listeners map[EventType][]EventListener
}

func makeEventsDispatcher() EventsDispatcher {
	return EventsDispatcher{listeners: make(map[EventType][]EventListener)}
}

func (e *EventsDispatcher) AddListener(event EventType, listener EventListener) {
	e.listeners[event] = append(e.listeners[event], listener)
}

func (e *EventsDispatcher) Dispatch(event EventType, value any) {
	for _, listener := range e.listeners[event] {
		listener(event, value)
	}
}
