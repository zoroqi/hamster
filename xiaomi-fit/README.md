将从 [account](https://account.xiaomi.com/) 导出的数据转换为 Markdown 格式. 具体的导出流程请参考 [Xiaomi Global Home](https://www.mi.com/global/support/article/KA-11566/).

我只是想要一个简单的 Markdown 格式进行备份, 所以这个脚本只会导出一些基本的信息.
我只选择导出心跳, 血压, 睡眠, 血氧和压力.

xiaomi fit 导出的数据是一个 zip 文件, 解压后的目录结构如下:
- xxx_MiFitness_hlth_center_fitness_data.csv 是身体指标的原始数据, 比如心跳, 睡眠, 血氧, 步行, 压力等
- xxx_MiFitness_hlth_center_sport_record.csv 是运动记录, 比如跑步, 游泳等

文件格式是 csv, 具体的数据字段是一个 json 串, 具体结构可以参考 types.go 文件中的定义(没有运动的信息).

使用 `xiaomi -start '2023-06-01' -end '2023-06-30' -path "data.csv"`

如果只是想简单看数据可以使用 sqlite 直接 load csv 文件, 然后使用 sql 查询.


开发中的一些问题和思考

功能并不复杂, 只是反复如何用泛型实现对代码的压缩, 然后发现 golang 泛型和 java 还是有很大的不同的.
比如以下的代码 java 是可以编译执行的, 而类似的 go 代码就不行了.
```java
Cell a = new Cell<String, Integer>();
a.t = "ab";
a.f = 10;
Cell b = new Cell<String, String>();
b.t = "cd";
b.f = "ef";
fu(a, b);
// 这句是可以编译通过, 也可以执行
fu(a,b);

static class Cell<T,F> {
    T t;
    F f;
}

static<T> void fu (Cell<T,Object> ...b) {
    for (Cell<T,Object> s : b) {
        System.out.println(s.t);
    }
}
```

```go
a := c2[string,int]{}
b := c2[string,string]{}
fu(a,b)

func fu[T any](c ...c2[T,any]) {
	for _, b := range c {
		fmt.Println(b.t)
	}
}

type c2[T, F any] struct {
	t T
	f F
}
```

我以为可以用 any 来实现兼容, 尝试下来发现没有成功, 而且不想去再去定义一个接口, 就作罢了.
现在想想可以补充一个 type set 的 interface 来实现, 但是现在不想改了.
当你泛型限制过于严苛的时候, 灵活度就不够, 但是类型是没有问题的. 而灵活度比较好的时候, 类型就不安全.
java 的擦除泛型实现了灵活性, 但是不安全, 而且为了方便我有时候就会选择灵活, 但是 go 让我无法钻空子了. :-P


另一个问题是 json 解析, 声明类型是 any, 而具体类型是 struct 执行json 解析的结果是 map 而不是我创建的类型
```
var t any
t := c2{}
json.Unmarshal(..., *t)
// t 的类型就是 map, 而不是 c2
// 最开始的想法就是 map[string]any 方式实现对任意指定类型的解析, 但是最后的结果就是, 最后发现不行.

var t any
t := &c2{}
json.Unmarshal(..., t)
// 这样 t 的类型就是 c2
```

