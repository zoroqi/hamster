package gostream

type Iterator[T any] func() (T, bool)

type Function[T any, R any] func(T) R
type Predicate[T any] func(T) bool
type Collector[T any, R any] func(T, R) R
