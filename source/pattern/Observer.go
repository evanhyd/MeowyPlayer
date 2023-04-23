package pattern

type ZeroArgObservabler interface {
	AddCallback(func())
}
type OneArgObservabler[T0 any] interface {
	AddCallback(func(T0))
}
type TwoArgObservabler[T0, T1 any] interface {
	AddCallback(func(T0, T1))
}
type ThreeArgObservabler[T0, T1, T2 any] interface {
	AddCallback(func(T0, T1, T2))
}
type FourArgObservabler[T0, T1, T2, T3 any] interface {
	AddCallback(func(T0, T1, T2, T3))
}

type ZeroArgObservable struct {
	callbacks []func()
}

func (subject *ZeroArgObservable) NotifyAll() {
	for _, callback := range subject.callbacks {
		go callback()
	}
}

func (subject *ZeroArgObservable) AddCallback(callback func()) {
	subject.callbacks = append(subject.callbacks, callback)
}

type OneArgObservable[T0 any] struct {
	callbacks []func(T0)
}

func (subject *OneArgObservable[T0]) NotifyAll(t0 T0) {
	for _, callback := range subject.callbacks {
		go callback(t0)
	}
}

func (subject *OneArgObservable[T]) AddCallback(callback func(T)) {
	subject.callbacks = append(subject.callbacks, callback)
}

type TwoArgObservable[T0, T1 any] struct {
	callbacks []func(T0, T1)
}

func (subject *TwoArgObservable[T0, T1]) NotifyAll(t0 T0, t1 T1) {
	for _, callback := range subject.callbacks {
		go callback(t0, t1)
	}
}

func (subject *TwoArgObservable[T1, T2]) AddCallback(callback func(T1, T2)) {
	subject.callbacks = append(subject.callbacks, callback)
}

type ThreeArgObservable[T0, T1, T2 any] struct {
	callbacks []func(T0, T1, T2)
}

func (subject *ThreeArgObservable[T0, T1, T2]) NotifyAll(t0 T0, t1 T1, t2 T2) {
	for _, callback := range subject.callbacks {
		go callback(t0, t1, t2)
	}
}

func (subject *ThreeArgObservable[T0, T1, T2]) AddCallback(callback func(T0, T1, T2)) {
	subject.callbacks = append(subject.callbacks, callback)
}

type FourArgObservable[T0, T1, T2, T3 any] struct {
	callbacks []func(T0, T1, T2, T3)
}

func (subject *FourArgObservable[T0, T1, T2, T3]) NotifyAll(t0 T0, t1 T1, t2 T2, t3 T3) {
	for _, callback := range subject.callbacks {
		go callback(t0, t1, t2, t3)
	}
}

func (subject *FourArgObservable[T0, T1, T2, T3]) AddCallback(callback func(T0, T1, T2, T3)) {
	subject.callbacks = append(subject.callbacks, callback)
}
