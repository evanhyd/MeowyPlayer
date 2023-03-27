package pattern

type ZeroArgSubject struct {
	observers []ZeroArgObserver
}

func (subject *ZeroArgSubject) NotifyAll() {
	for _, observer := range subject.observers {
		observer.Notify()
	}
}

func (subject *ZeroArgSubject) AddObserver(observer ZeroArgObserver) {
	subject.observers = append(subject.observers, observer)
}

func (subject *ZeroArgSubject) RemoveObserver(observer ZeroArgObserver) {
	lastIndex := len(subject.observers) - 1
	for i := range subject.observers {
		if subject.observers[i] == observer {
			subject.observers[i], subject.observers[lastIndex] = subject.observers[lastIndex], subject.observers[i]
			subject.observers = subject.observers[:lastIndex]
			break
		}
	}
}

type OneArgSubject[T any] struct {
	observers []OneArgObserver[T]
}

func (subject *OneArgSubject[T]) NotifyAll(t T) {
	for _, observer := range subject.observers {
		observer.Notify(t)
	}
}

func (subject *OneArgSubject[T]) AddObserver(observer OneArgObserver[T]) {
	subject.observers = append(subject.observers, observer)
}

func (subject *OneArgSubject[T]) RemoveObserver(observer OneArgObserver[T]) {
	lastIndex := len(subject.observers) - 1
	for i := range subject.observers {
		if subject.observers[i] == observer {
			subject.observers[i], subject.observers[lastIndex] = subject.observers[lastIndex], subject.observers[i]
			subject.observers = subject.observers[:lastIndex]
			break
		}
	}
}

type TwoArgSubject[T0, T1 any] struct {
	observers []TwoArgObserver[T0, T1]
}

func (subject *TwoArgSubject[T0, T1]) NotifyAll(t0 T0, t1 T1) {
	for _, observer := range subject.observers {
		observer.Notify(t0, t1)
	}
}

func (subject *TwoArgSubject[T0, T1]) AddObserver(observer TwoArgObserver[T0, T1]) {
	subject.observers = append(subject.observers, observer)
}

func (subject *TwoArgSubject[T0, T1]) RemoveObserver(observer TwoArgObserver[T0, T1]) {
	lastIndex := len(subject.observers) - 1
	for i := range subject.observers {
		if subject.observers[i] == observer {
			subject.observers[i], subject.observers[lastIndex] = subject.observers[lastIndex], subject.observers[i]
			subject.observers = subject.observers[:lastIndex]
			break
		}
	}
}

type ThreeArgSubject[T0, T1, T2 any] struct {
	observers []ThreeArgObserver[T0, T1, T2]
}

func (subject *ThreeArgSubject[T0, T1, T2]) NotifyAll(t0 T0, t1 T1, t2 T2) {
	for _, observer := range subject.observers {
		observer.Notify(t0, t1, t2)
	}
}

func (subject *ThreeArgSubject[T0, T1, T2]) AddObserver(observer ThreeArgObserver[T0, T1, T2]) {
	subject.observers = append(subject.observers, observer)
}

func (subject *ThreeArgSubject[T0, T1, T2]) RemoveObserver(observer ThreeArgObserver[T0, T1, T2]) {
	lastIndex := len(subject.observers) - 1
	for i := range subject.observers {
		if subject.observers[i] == observer {
			subject.observers[i], subject.observers[lastIndex] = subject.observers[lastIndex], subject.observers[i]
			subject.observers = subject.observers[:lastIndex]
			break
		}
	}
}

type FourArgSubject[T0, T1, T2, T3 any] struct {
	observers []FourArgObserver[T0, T1, T2, T3]
}

func (subject *FourArgSubject[T0, T1, T2, T3]) NotifyAll(t0 T0, t1 T1, t2 T2, t3 T3) {
	for _, observer := range subject.observers {
		observer.Notify(t0, t1, t2, t3)
	}
}

func (subject *FourArgSubject[T0, T1, T2, T3]) AddObserver(observer FourArgObserver[T0, T1, T2, T3]) {
	subject.observers = append(subject.observers, observer)
}

func (subject *FourArgSubject[T0, T1, T2, T3]) RemoveObserver(observer FourArgObserver[T0, T1, T2, T3]) {
	lastIndex := len(subject.observers) - 1
	for i := range subject.observers {
		if subject.observers[i] == observer {
			subject.observers[i], subject.observers[lastIndex] = subject.observers[lastIndex], subject.observers[i]
			subject.observers = subject.observers[:lastIndex]
			break
		}
	}
}
