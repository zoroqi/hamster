package go_stream

type ChanStream[T any] struct {
	out chan T
}

func NewChanStreamByArray[T any](a []T) *ChanStream[T] {
	out := make(chan T)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return &ChanStream[T]{
		out: out,
	}
}

func NewChanStreamByChan[T any](a chan T) *ChanStream[T] {
	out := make(chan T)
	go func() {
		for v := range a {
			out <- v
		}
		close(out)
	}()
	return &ChanStream[T]{
		out: out,
	}
}

type Iterator[T any] interface {
	Next() T
	HasNext() bool
}

func NewChanStreamByIterator[T any](a Iterator[T]) *ChanStream[T] {
	out := make(chan T)
	go func() {
		for a.HasNext() {
			out <- a.Next()
		}
		close(out)
	}()
	return &ChanStream[T]{
		out: out,
	}
}

func Filter[T any](f func(T) bool, stream *ChanStream[T]) *ChanStream[T] {
	out := make(chan T)
	go func() {
		for v := range stream.out {
			if f(v) {
				out <- v
			}
		}
		close(out)
	}()
	return &ChanStream[T]{
		out: out,
	}
}

func Map[T any, R any](m func(T) R, stream *ChanStream[T]) *ChanStream[R] {
	out := make(chan R)
	go func() {
		for v := range stream.out {
			out <- m(v)
		}
		close(out)
	}()
	return &ChanStream[R]{
		out: out,
	}
}

func Reduce[T any, R any](r func(T, R) R, acc R, stream *ChanStream[T]) R {
	for v := range stream.out {
		acc = r(v, acc)
	}
	return acc
}

func CollectToArray[T any](stream *ChanStream[T]) []T {
	var c []T
	for v := range stream.out {
		c = append(c, v)
	}
	return c
}
