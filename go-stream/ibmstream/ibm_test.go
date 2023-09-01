package ibmstream

import (
	"github.com/IBM/fp-go/iterator/stateless"
	"testing"
)

func BenchmarkMulti(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		TMulti()
	}
}

func Test_Ibm(t *testing.T) {
	if TMulti() != 16242924 {
		t.Fatal("err")
	}
}

func TestAlloc(t *testing.T) {
	f := testing.AllocsPerRun(5, func() {
		TMulti()
	})
	t.Log(f)
}

func TMulti() int64 {
	i64to32 := stateless.Map(func(t int64) int {
		return int(t)
	})
	i32to64 := stateless.Map(func(t int) int64 {
		return int64(t)
	})
	even := stateless.Filter(func(n int) bool {
		return n%2 == 0
	})
	no4 := stateless.Filter(func(t int64) bool {
		n := t
		for n/10 != 0 {
			if n%10 == 4 {
				return false
			}
			n = n / 10
		}
		return true
	})

	sum := stateless.Reduce(func(n int64, r int) int64 {
		return n + int64(r)
	}, int64(0))
	length := 10000
	nums := make([]int, 0, length)
	for i := 0; i < length; i++ {
		nums = append(nums, i+1)
	}
	arr := stateless.FromArray(nums)
	return sum(i64to32(no4(i32to64(even(arr)))))
}
