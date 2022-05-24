# go stream

测试的 go stream 项目, 提供了几个链式调用的基础函数. 使用效果整体评价不是很好.

最开始再想用 chan 做数据传输会泄漏吗? 或会出现没有消费干净吗? 想了下只要最后一个处理性质的是同步的就不会出现这两种情况

主要原因

一. 是 golang 语法本身是无法在方法上再追加指定新的泛型

二. 没有简写的语法糖导致链式长一点, 代码看起来就比较痛苦了.
```go
// 代码没啥目的为了嵌套而嵌套的
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
)
```
