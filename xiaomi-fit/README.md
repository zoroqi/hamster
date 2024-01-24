将从 [account](https://account.xiaomi.com/) 导出的数据转换为 Markdown 格式. 具体的导出流程请参考 [Xiaomi Global Home](https://www.mi.com/global/support/article/KA-11566/).

我只是想要一个简单的 Markdown 格式进行备份, 所以这个脚本只会导出一些基本的信息.
我只选择导出心跳, 血压, 睡眠, 血氧和压力.

xiaomi fit 导出的数据是一个 zip 文件, 解压后的目录结构如下:
- xxx_MiFitness_hlth_center_fitness_data.csv 是身体指标的原始数据, 比如心跳, 睡眠, 血氧, 步行, 压力等
- xxx_MiFitness_hlth_center_sport_record.csv 是运动记录, 比如跑步, 游泳等

文件格式是 csv, 具体的数据字段是一个 json 串, 具体结构可以参考 types.go 文件中的定义(没有运动的信息).

使用 `xiaomi -start '2023-06-01' -end '2023-06-30' -path "data.csv"`

如果只是想简单看数据可以使用 sqlite 直接 load csv 文件, 然后使用 sql 查询.