package go_stream

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	stream := NewChanStreamByArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	r := CollectToArray(
		Filter(func(n int) bool {
			return n <= 5
		}, stream))
	if !reflect.DeepEqual(r, []int{1, 2, 3, 4, 5}) {
		t.Error("Expected to return 1,2,3,4,5; but returned ", r)
	}
}

func TestMap(t *testing.T) {
	stream := NewChanStreamByArray([]int{1, 2, 3, 4, 5})
	r := CollectToArray(
		Map(func(n int) int64 {
			return int64(n * 2)
		}, stream))
	if !reflect.DeepEqual(r, []int64{2, 4, 6, 8, 10}) {
		t.Error("Expected to return 2,4,6,8,10; but returned ", r)
	}
}

func TestReduce(t *testing.T) {
	stream := NewChanStreamByArray([]int{1, 2, 3, 4, 5})
	r := Reduce(func(n int, r int64) int64 {
		return int64(n) + r
	}, 0, stream)
	if !reflect.DeepEqual(r, int64(15)) {
		t.Error("Expected to return 15; but returned ", r)
	}
}

func NewNumIter(start, end int) Iterator[int] {
	return &iter{start: start, end: end, i: start}
}

type iter struct {
	start int
	end   int
	i     int
}

func (i *iter) Next() int {
	n := i.i
	i.i++
	return n
}

func (i *iter) HasNext() bool {
	return i.i <= i.end
}

// haskell: sum . map (\x -> read x :: Int) . filter (\x -> not . elem '4' $ x) . map show . filter (\x-> mod x 2 == 0) $ [1..100]
func TestMulti(t *testing.T) {
	nums := NewNumIter(1, 100)
	stream := NewChanStreamByIterator(nums)
	r := Reduce(
		func(n int, r int64) int64 {
			return int64(n) + r
		}, 0,
		Map(func(t string) int {
			n, _ := strconv.Atoi(t)
			return n
		},
			Filter(func(t string) bool {
				return !strings.Contains(t, "4")
			},
				Map(func(t int64) string {
					return fmt.Sprintf("%d", t)
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
		),
	)
	if !reflect.DeepEqual(r, int64(1884)) {
		t.Error("Expected to return 15; but returned ", r)
	}
}
