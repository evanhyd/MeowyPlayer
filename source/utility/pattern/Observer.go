package pattern

import (
	"slices"
)

/*
A generic observer interface.
*/
type Observer[T any] interface {
	Notify(T)
}

/*
A generic subject interface.
*/
type Subject[T any] interface {
	Attach(Observer[T])
	Detach(Observer[T])
	NotifyAll(T)
}

/*
A generic subject base class.
*/
type SubjectBase[T any] struct {
	observers []Observer[T]
}

func (s *SubjectBase[T]) Attach(observer Observer[T]) {
	s.observers = append(s.observers, observer)
}

func (s *SubjectBase[T]) Detach(observer Observer[T]) {
	index := slices.Index(s.observers, observer)
	s.observers[index] = s.observers[len(s.observers)-1]
	s.observers = s.observers[:len(s.observers)-1]
}

func (s *SubjectBase[T]) NotifyAll(t T) {
	for _, observer := range s.observers {
		go observer.Notify(t)
	}
}

/*
A generic value wrapper that notify the observers whenever the value gets set.
*/
type Data[T any] struct {
	SubjectBase[T]
	value T
}

func (d *Data[T]) Get() T {
	return d.value
}

func (d *Data[T]) Set(v T) {
	d.value = v
	d.SubjectBase.NotifyAll(d.value)
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
