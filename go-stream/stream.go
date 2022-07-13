package gostream

import "fmt"

type Iterator[T any] func() (T, bool)

type Function[T any, R any] func(T) R
type Predicate[T any] func(T) bool
type Collector[T any, R any] func(T, R) R

type F2 interface {
	Get(int) int
}

type F3[T any] interface {
	Get(T) T
}

type Fu func(int) int

func (f Fu) Get(i int) int {
	return f(i)
}

func dd(f F3[int], n int) int {
	return f.Get(n)
}

func ddd() {
	var f Fu
	f = func(i int) int {
		return i * 2
	}
	fmt.Println(dd(f, 10))
}
