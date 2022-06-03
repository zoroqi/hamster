# gostream

一个实验性的工程, go1.18 引入了泛型, 想看看实现 stream 效果如何. 总共有三种方案
1. 基于 chan 的
2. 基于迭代器的
3. 基于 slice 的

最后测试结果使用迭代器是最优的. 因为 chan 方案耗时过于严重; slice 方案会使用过多的内存, 并且无法对无限流进行求值. 迭代器可以两者兼顾, 所以是一个不错的选择.

测试结果
```
pkg: github.com/zoroqi/gostream/cstream
cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
BenchmarkMulti-8         492       2587500 ns/op       82985 B/op         16 allocs/op
BenchmarkMulti-8         463       2633396 ns/op       82979 B/op         16 allocs/op
BenchmarkMulti-8         472       3053850 ns/op       82939 B/op         16 allocs/op
BenchmarkMulti-8         439       3624352 ns/op       82930 B/op         16 allocs/op

pkg: github.com/zoroqi/gostream/istream
cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
BenchmarkMulti-8        5974        197930 ns/op       82161 B/op         12 allocs/op
BenchmarkMulti-8        7527        151803 ns/op       82160 B/op         12 allocs/op
BenchmarkMulti-8        6940        185106 ns/op       82160 B/op         12 allocs/op
BenchmarkMulti-8        5409        203937 ns/op       82160 B/op         12 allocs/op

pkg: github.com/zoroqi/gostream/sstream
cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
BenchmarkMulti-8        7393        193657 ns/op      365778 B/op         38 allocs/op
BenchmarkMulti-8        6175        236944 ns/op      365777 B/op         38 allocs/op
BenchmarkMulti-8        6452        172061 ns/op      365777 B/op         38 allocs/op
BenchmarkMulti-8        7548        144701 ns/op      365777 B/op         38 allocs/op
```

## 流式现在适合 golang 吗?

从写代码可以看出, 现在流式代码并不适合 golang.

### 原因一, method 不支持单独指定类型参数

这个原因导致无法实现类似效果的代码, 只能选择使用参数进行传递 Stream.
```golang
type Stream[T any] struct {
    ...
}
func (s *Stream[T]) Map[R any](f func(T)R) *Stream[R] {
    ...
}
```

我使用了是两种嵌套方式. 这两种方式不管怎么看都感觉怪怪的(我更喜欢第一种). 当流程稍微长一点代码就不易理解.
```golang
Reduce(
        func(n int, r int64) int64 {
            return int64(n) + r
        }, 0,
        Map(func(t int64) int {
            return int(t)
        },
            Filter(func(n int64) bool {
                return n%2 == 0
            }, stream),
        ),
    )
```

另一种方式
```golang
Reduce(
    Map(
        Filter(stream,
            func(n int64) bool {
                return n%2 == 0
            }),
                func(t int64) int {
                    return int(t)
                }),
                func(n int, r int64) int64 {
                    return int64(n) + r
                }, 0 )
```

### 原因二, 没有 lambda 语法糖

原因二导致写起来代码很啰嗦, 不能省略方法声明.
