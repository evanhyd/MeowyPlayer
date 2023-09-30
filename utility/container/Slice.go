package container

type Vector[T any] []T

func (v *Vector[T]) PushBack(data T) {
	*v = append(*v, data)
}

func (v *Vector[T]) PopBack() {
	*v = (*v)[:len(*v)-1]
}

func (v *Vector[T]) Back() *T {
	return &(*v)[len(*v)-1]
}

func (v *Vector[T]) Size() int {
	return len(*v)
}

func (v *Vector[T]) Empty() bool {
	return v.Size() > 0
}

func (v *Vector[T]) Remove(index int) {
	if 0 <= index && index < v.Size() {
		(*v)[index] = *v.Back()
		v.PopBack()
	}
}

func (v *Vector[T]) Filter(filter func(T) bool) Vector[T] {
	var remain []T
	for _, v := range *v {
		if filter(v) {
			remain = append(remain, v)
		}
	}
	return remain
}
