package util

var _ Subject[int] = &subjectBase[int]{}

type Observer[T any] interface {
	Notify(T)
}

type Subject[T any] interface {
	Attach(Observer[T])
	AttachFunc(func(T))
	NotifyAll(T)
}

type subjectBase[T any] struct {
	observers []func(T)
}

func MakeSubject[T any]() Subject[T] {
	return &subjectBase[T]{}
}

func (s *subjectBase[T]) Attach(observer Observer[T]) {
	s.observers = append(s.observers, observer.Notify)
}

func (s *subjectBase[T]) AttachFunc(callback func(T)) {
	s.observers = append(s.observers, callback)
}

func (s *subjectBase[T]) NotifyAll(t T) {
	for _, observer := range s.observers {
		go observer(t) //does threading work well with GUI?
	}
}
