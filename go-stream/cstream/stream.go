package cstream

import "github.com/zoroqi/gostream"

type ChanStream[T any] struct {
	out chan T
}

const chanl = 5

func NewChanStreamBySlice[T any](a []T) *ChanStream[T] {
	out := make(chan T, chanl)
	go func() {
		defer close(out)
		for _, v := range a {
			out <- v
		}
	}()
	return &ChanStream[T]{
		out: out,
	}
}

func NewChanStreamByChan[T any](a chan T) *ChanStream[T] {
	out := make(chan T, chanl)
	go func() {
		defer close(out)
		for v := range a {
			out <- v
		}
	}()
	return &ChanStream[T]{
		out: out,
	}
}

func NewChanStreamByIterator[T any](a gostream.Iterator[T]) *ChanStream[T] {
	out := make(chan T, chanl)
	go func() {
		defer close(out)
		for v, r := a(); r; v, r = a() {
			out <- v
		}

	}()
	return &ChanStream[T]{
		out: out,
	}
}

func Filter[T any](f gostream.Predicate[T], stream *ChanStream[T]) *ChanStream[T] {
	out := make(chan T, chanl)
	go func() {
		defer close(out)
		for v := range stream.out {
			if f(v) {
				out <- v
			}
		}
	}()
	return &ChanStream[T]{
		out: out,
	}
}

func Map[T any, R any](m gostream.Function[T, R], stream *ChanStream[T]) *ChanStream[R] {
	out := make(chan R, chanl)
	go func() {
		defer close(out)
		for v := range stream.out {
			out <- m(v)
		}
	}()
	return &ChanStream[R]{
		out: out,
	}
}

func Reduce[T any, R any](r gostream.Collector[T, R], acc R, stream *ChanStream[T]) R {
	for v := range stream.out {
		acc = r(v, acc)
	}
	return acc
}
