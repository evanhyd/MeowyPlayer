package pattern

import (
	"golang.org/x/exp/slices"
	"meowyplayer.com/utility/container"
)

/*
A generic observer interface.
*/
type Observer[T any] interface {
	Notify(T)
}

/*
A generic subject class.
*/
type Subject[T any] struct {
	observers container.Slice[Observer[T]]
}

func (s *Subject[T]) Attach(observer Observer[T]) {
	s.observers.PushBack(observer)
}

func (s *Subject[T]) Detach(observer Observer[T]) {
	s.observers.Remove(slices.Index(s.observers, observer))
}

func (s *Subject[T]) NotifyAll(t T) {
	for _, observer := range s.observers {
		go observer.Notify(t)
	}
}

/*
A generic value wrapper that notify the observers whenever the value gets set.
*/
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

/*
A generic function callback maker that can be used as a observer.
*/
type callbackWrapper[T any] struct {
	function func(T)
}

func (c *callbackWrapper[T]) Notify(data T) {
	c.function(data)
}

func MakeCallback[T any](function func(T)) Observer[T] {
	return &callbackWrapper[T]{function}
}
