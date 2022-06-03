package istream

import "github.com/zoroqi/gostream"

type IteratorStream[T any] struct {
	i gostream.Iterator[T]
}

func NewIteratorStreamBySlice[T any](a []T) *IteratorStream[T] {
	i := 0
	l := len(a)
	iter := func() (T, bool) {
		if i < l {
			m := a[i]
			i++
			return m, true
		}
		var vv T
		return vv, false
	}
	return &IteratorStream[T]{
		i: iter,
	}
}

func NewIteratorStreamByChan[T any](a chan T) *IteratorStream[T] {
	iter := func() (T, bool) {
		select {
		case v, c := <-a:
			return v, c
		}
	}
	return &IteratorStream[T]{
		i: iter,
	}
}

func NewIteratorStreamByIterator[T any](a gostream.Iterator[T]) *IteratorStream[T] {
	return &IteratorStream[T]{
		i: a,
	}
}

func Filter[T any](f gostream.Predicate[T], stream *IteratorStream[T]) *IteratorStream[T] {
	iter := func() (T, bool) {
		for v, run := stream.i(); run; v, run = stream.i() {
			if f(v) {
				return v, run
			}
		}
		var vv T
		return vv, false
	}
	return &IteratorStream[T]{
		i: iter,
	}
}

func Map[T any, R any](m gostream.Function[T, R], stream *IteratorStream[T]) *IteratorStream[R] {
	iter := func() (R, bool) {
		v, e := stream.i()
		if e {
			return m(v), e
		} else {
			var vv R
			return vv, false
		}
	}
	return &IteratorStream[R]{
		i: iter,
	}
}

func Reduce[T any, R any](r gostream.Collector[T, R], acc R, stream *IteratorStream[T]) R {
	for v, run := stream.i(); run; v, run = stream.i() {
		acc = r(v, acc)
	}
	return acc
}
