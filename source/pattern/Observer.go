package pattern

type ZeroArgObserver interface {
	Notify()
}

type OneArgObserver[T any] interface {
	Notify(T)
}

type TwoArgObserver[T0, T1 any] interface {
	Notify(T0, T1)
}

type ThreeArgObserver[T0, T1, T2 any] interface {
	Notify(T0, T1, T2)
}

type FourArgObserver[T0, T1, T2, T3 any] interface {
	Notify(T0, T1, T2, T3)
}
