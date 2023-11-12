package container

type Slice[T any] []T

func (v *Slice[T]) PushBack(data T) {
	*v = append(*v, data)
}

func (v *Slice[T]) PopBack() {
	*v = (*v)[:len(*v)-1]
}

func (v *Slice[T]) Back() *T {
	return &(*v)[len(*v)-1]
}

func (v *Slice[T]) Empty() bool {
	return len(*v) == 0
}

func (v *Slice[T]) Clear() {
	*v = []T{}
}
