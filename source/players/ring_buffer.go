package players

type ringBuffer[T any] struct {
	data  []T
	index int
}

func (b *ringBuffer[T]) set(data []T) {
	b.data = data
	b.index = 0
}

func (b *ringBuffer[T]) empty() bool {
	return len(b.data) == 0
}

func (b ringBuffer[T]) current() T {
	return b.data[b.index]
}

func (b *ringBuffer[T]) next() T {
	b.index = (b.index + 1) % len(b.data)
	return b.data[b.index]
}

func (b *ringBuffer[T]) prev() T {
	b.index = (b.index - 1 + len(b.data)) % len(b.data)
	return b.data[b.index]
}
