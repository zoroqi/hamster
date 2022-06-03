package sstream

import "github.com/zoroqi/gostream"

type SliceStream[T any] struct {
	out []T
}

func NewSliceStreamBySlice[T any](a []T) *SliceStream[T] {
	return &SliceStream[T]{out: a}
}

func Filter[T any](f gostream.Predicate[T], stream *SliceStream[T]) *SliceStream[T] {
	out := make([]T, 0)
	for _, v := range stream.out {
		if f(v) {
			out = append(out, v)
		}
	}
	return &SliceStream[T]{
		out: out,
	}
}

func Map[T any, R any](m gostream.Function[T, R], stream *SliceStream[T]) *SliceStream[R] {
	out := make([]R, 0, len(stream.out))
	for _, v := range stream.out {
		out = append(out, m(v))
	}
	return &SliceStream[R]{
		out: out,
	}
}

func Reduce[T any, R any](r gostream.Collector[T, R], acc R, stream *SliceStream[T]) R {
	for _, v := range stream.out {
		acc = r(v, acc)
	}
	return acc
}
