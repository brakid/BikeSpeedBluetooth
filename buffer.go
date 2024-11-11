package main

import (
	"fmt"
)

type Buffer[T any] struct {
	values   []T
	index    uint
	capacity uint
}

func CreateBuffer[T any](capacity uint) (*Buffer[T], error) {
	if capacity == 0 {
		return nil, fmt.Errorf("capacity must be > 0")
	}
	buffer := Buffer[T]{
		make([]T, capacity),
		0,
		capacity,
	}

	return &buffer, nil
}

func (b *Buffer[T]) Append(value T) {
	b.values[b.index] = value
	b.index = (b.index + 1) % b.capacity
}

func (b *Buffer[T]) Get() []T {
	return append(b.values[b.index:], b.values[:b.index]...)
}
