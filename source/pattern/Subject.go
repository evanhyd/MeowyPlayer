package pattern

type ZeroArgSubject struct {
	observers []ZeroArgObserver
	callbacks []func()
}

func (subject *ZeroArgSubject) NotifyAll() {
	for _, observer := range subject.observers {
		observer.Notify()
	}
	for _, callback := range subject.callbacks {
		callback()
	}
}

func (subject *ZeroArgSubject) AddObserver(observer ZeroArgObserver) {
	subject.observers = append(subject.observers, observer)
}

func (subject *ZeroArgSubject) AddCallback(callback func()) {
	subject.callbacks = append(subject.callbacks, callback)
}

type OneArgSubject[T any] struct {
	observers []OneArgObserver[T]
	callbacks []func(T)
}

func (subject *OneArgSubject[T]) NotifyAll(t T) {
	for _, observer := range subject.observers {
		observer.Notify(t)
	}
	for _, observer := range subject.callbacks {
		observer(t)
	}
}

func (subject *OneArgSubject[T]) AddObserver(observer OneArgObserver[T]) {
	subject.observers = append(subject.observers, observer)
}

func (subject *OneArgSubject[T]) AddCallback(callback func(T)) {
	subject.callbacks = append(subject.callbacks, callback)
}

type TwoArgSubject[T0, T1 any] struct {
	observers []TwoArgObserver[T0, T1]
	callbacks []func(T0, T1)
}

func (subject *TwoArgSubject[T0, T1]) NotifyAll(t0 T0, t1 T1) {
	for _, observer := range subject.observers {
		observer.Notify(t0, t1)
	}
	for _, observer := range subject.callbacks {
		observer(t0, t1)
	}
}

func (subject *TwoArgSubject[T0, T1]) AddObserver(observer TwoArgObserver[T0, T1]) {
	subject.observers = append(subject.observers, observer)
}

func (subject *TwoArgSubject[T0, T1]) AddCallback(callback func(T0, T1)) {
	subject.callbacks = append(subject.callbacks, callback)
}

type ThreeArgSubject[T0, T1, T2 any] struct {
	observers []ThreeArgObserver[T0, T1, T2]
	callbacks []func(T0, T1, T2)
}

func (subject *ThreeArgSubject[T0, T1, T2]) NotifyAll(t0 T0, t1 T1, t2 T2) {
	for _, observer := range subject.observers {
		observer.Notify(t0, t1, t2)
	}
	for _, observer := range subject.callbacks {
		observer(t0, t1, t2)
	}
}

func (subject *ThreeArgSubject[T0, T1, T2]) AddObserver(observer ThreeArgObserver[T0, T1, T2]) {
	subject.observers = append(subject.observers, observer)
}

func (subject *ThreeArgSubject[T0, T1, T2]) AddCallback(callback func(T0, T1, T2)) {
	subject.callbacks = append(subject.callbacks, callback)
}

type FourArgSubject[T0, T1, T2, T3 any] struct {
	observers []FourArgObserver[T0, T1, T2, T3]
	callbacks []func(T0, T1, T2, T3)
}

func (subject *FourArgSubject[T0, T1, T2, T3]) NotifyAll(t0 T0, t1 T1, t2 T2, t3 T3) {
	for _, observer := range subject.observers {
		observer.Notify(t0, t1, t2, t3)
	}
	for _, observer := range subject.callbacks {
		observer(t0, t1, t2, t3)
	}
}

func (subject *FourArgSubject[T0, T1, T2, T3]) AddObserver(observer FourArgObserver[T0, T1, T2, T3]) {
	subject.observers = append(subject.observers, observer)
}

func (subject *FourArgSubject[T0, T1, T2, T3]) AddCallback(callback func(T0, T1, T2, T3)) {
	subject.callbacks = append(subject.callbacks, callback)
}
