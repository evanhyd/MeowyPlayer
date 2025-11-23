package context

type Event int64
type EventListener = func(Event, any)

const (
	OnLoadPlaylistToMusicPlayerEvent Event = iota
)

type EventDispatcher struct {
	listeners map[Event][]EventListener
}

func makeEventDispatcher() EventDispatcher {
	return EventDispatcher{
		listeners: make(map[Event][]EventListener),
	}
}

func (e *EventDispatcher) addListener(event Event, listener EventListener) {
	e.listeners[event] = append(e.listeners[event], listener)
}

func (e *EventDispatcher) dispatch(event Event, value any) {
	for _, listener := range e.listeners[event] {
		listener(event, value)
	}
}
