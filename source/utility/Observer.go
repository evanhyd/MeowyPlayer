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

type Data[T any] struct {
	Subject[*T]
	value T
}

func (d *Data[T]) NotifyAll() {
	d.Subject.NotifyAll(&d.value)
}

func (d *Data[T]) Set(v *T) {
	d.value = *v
	d.NotifyAll()
}

func (d *Data[T]) Get() *T {
	return &d.value
}

type callbackWrapper[T any] struct {
	function func(T)
}

func (c *callbackWrapper[T]) Notify(data T) {
	c.function(data)
}

func MakeCallback[T any](function func(T)) Observer[T] {
	call := callbackWrapper[T]{function}
	return &call
}
