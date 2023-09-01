package istream

import (
	"testing"
)

func BenchmarkMulti(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		TMulti()
	}
}

func TestAlloc(t *testing.T) {
	f := testing.AllocsPerRun(5, func() {
		TMulti()
	})
	t.Log(f)
}

func TMulti() int64 {
	length := 10000
	nums := make([]int, 0, length)
	for i := 0; i < length; i++ {
		nums = append(nums, i+1)
	}
	stream := NewIteratorStreamBySlice(nums)
	return Reduce(
		func(n int, r int64) int64 {
			return int64(n) + r
		}, 0,
		Map(func(t int64) int {
			return int(t)
		},
			Filter(func(t int64) bool {
				n := t
				for n/10 != 0 {
					if n%10 == 4 {
						return false
					}
					n = n / 10
				}
				return true
			},

				Map(func(t int) int64 {
					return int64(t)
				},
					Filter(func(n int) bool {
						return n%2 == 0
					}, stream),
				),
			),
		),
	)
}
