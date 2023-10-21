package container

import "meowyplayer.com/utility/assert"

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

func (v *Slice[T]) Size() int {
	return len(*v)
}

func (v *Slice[T]) Empty() bool {
	return v.Size() == 0
}

func (v *Slice[T]) Clear() {
	*v = []T{}
}

func (v *Slice[T]) Remove(index int) {
	assert.Ensure(func() bool { return 0 <= index && index < v.Size() })
	(*v)[index] = *v.Back()
	v.PopBack()
}

func (v *Slice[T]) Filter(filter func(T) bool) Slice[T] {
	var remain []T
	for _, v := range *v {
		if filter(v) {
			remain = append(remain, v)
		}
	}
	return remain
}
