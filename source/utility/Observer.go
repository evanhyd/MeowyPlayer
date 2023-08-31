package utility

import "golang.org/x/exp/slices"

type Observer[T any] interface {
	Notify(T)
}

type Subject[T any] struct {
	observers []Observer[T]
}

func (s *Subject[T]) Attach(observer Observer[T]) {
	s.observers = append(s.observers, observer)
}

func (s *Subject[T]) Detach(observer Observer[T]) {
	index := slices.Index(s.observers, observer)
	last := len(s.observers) - 1
	s.observers[index] = s.observers[last]
	s.observers = s.observers[:last]
}

func (s *Subject[T]) NotifyAll(t T) {
	for _, observer := range s.observers {
		go observer.Notify(t)
	}
}
