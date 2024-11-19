package main

import (
	"flag"
	"fmt"
	"math/bits"
	"math/rand/v2"
)

type 序卦 = uint8

type 卦 struct {
	卦名 string
	卦象 string
	序号 序卦
}

func (g 卦) String() string {
	return fmt.Sprintf("%s(%s)", g.卦名, g.卦象)
}

type 爻 = uint8

const 阴爻 = 爻(0)
const 阳爻 = 爻(1)

func 解(y 爻) (r 爻, d bool) {
	switch y {
	case 0:
		r = 阴爻
		d = true
	case 1, 2, 4:
		r = 阳爻
		d = false
	case 3, 5, 6:
		r = 阴爻
		d = false
	case 7:
		r = 阳爻
		d = true
	}
	return
}

func yaoString(y 爻, d bool) (s string) {
	if y == 阳爻 {
		s = "阳"
	} else {
		s = "阴"
	}
	s += "爻"
	if d {
		s += ", 动爻"
	}
	return
}

const 次数 = 6

func 摇动(龟壳 func() 爻) ([]uint8, []bool, uint8) {
	卦象 := uint8(0)
	爻像 := make([]uint8, 次数)
	爻动 := make([]bool, 次数)
	for i := 0; i < 次数; i++ {
		爻像[i], 爻动[i] = 解(龟壳())
		卦象 = 爻像[i]<<i | 卦象
	}
	return 爻像, 爻动, 卦象
}

var 龟壳 = func() 爻 {
	return 爻(rand.Int32N(8))
}

func main() {
	flag.Parse()
	question := ""
	if flag.NArg() != 0 {
		question = flag.Arg(0)
	}

	爻像, 爻动, 卦象 := 摇动(龟壳)

	上挂 := 卦象 >> 3
	下卦 := 卦象 & 0x07

	fmt.Println("你是鬼谷子, 根据我摇卦的结果进行占卜.")
	fmt.Println("- 占卜的问题:", question)
	fmt.Println("- 卦象:\n```")
	// [本卦/错卦/综卦/复卦/象卦/交卦/变卦/杂卦解读](https://www.shuozhouyi.com/25248.html)
	fmt.Printf("本卦: %s\n", 六十四卦[卦象])
	fmt.Printf("上卦: %s 下卦: %s\n", 八卦[上挂], 八卦[下卦])
	fmt.Printf("错卦: %s\n", 六十四卦[卦象^uint8(0x3F)])
	fmt.Printf("宗卦: %s\n", 六十四卦[bits.Reverse8(卦象)>>2])
	fmt.Printf("互(复)卦: %s\n", 六十四卦[((((上挂&0x03)<<1)|(下卦>>2))<<3)|((上挂&0x01<<2)|(下卦>>1))])

	for i := 次数; i > 0; i-- {
		i := i - 1
		fmt.Printf("%s爻: %s\n", 爻名[i], yaoString(爻像[i], 爻动[i]))
	}
	fmt.Println("```")
    fmt.Println("- 解卦的模板:\n```")
    fmt.Println(`卦象解析:
上挂解析:
下卦解析:
错卦解析:
宗卦解析:
互卦解析:
建议:`)
    fmt.Println("```")
}

var 八卦 = map[序卦]卦{}
var 六十四卦 = map[序卦]卦{}
var 爻名 = []string{"初", "二", "三", "四", "五", "上"}

func buildGuaStruct(img, txt []string, steps []int) map[序卦]卦 {
	kv := map[string]卦{}
	for i, v := range img {
		kv[v] = 卦{
			卦名: txt[i],
			卦象: v,
		}
	}
	for _, step := range steps {
		for i := 0; i < len(img); i++ {
			h := step / 2
			ba := kv[img[i]]
			if i%step < h {
				ba.序号 = ba.序号<<1 | 0
			} else {
				ba.序号 = ba.序号<<1 | 1
			}
			kv[img[i]] = ba
		}
	}
	r := map[序卦]卦{}
	for _, v := range kv {
		r[v.序号] = v
	}
	return r
}

func init() {
	八卦 = buildGuaStruct(
		[]string{"☷", "☶", "☵", "☴", "☳", "☲", "☱", "☰"},
		[]string{"坤", "艮", "坎", "巽", "震", "离", "兑", "乾"},
		[]int{2, 4, 8},
	)

	六十四卦 = buildGuaStruct(
		[]string{
			"䷁", "䷖", "䷇", "䷓", "䷏", "䷢", "䷬", "䷋",
			"䷎", "䷳", "䷦", "䷴", "䷽", "䷷", "䷞", "䷠",
			"䷆", "䷃", "䷜", "䷺", "䷧", "䷿", "䷮", "䷅",
			"䷭", "䷑", "䷯", "䷸", "䷟", "䷱", "䷛", "䷫",
			"䷗", "䷚", "䷂", "䷩", "䷲", "䷔", "䷐", "䷘",
			"䷣", "䷕", "䷾", "䷤", "䷶", "䷝", "䷰", "䷌",
			"䷒", "䷨", "䷻", "䷼", "䷵", "䷥", "䷹", "䷉",
			"䷊", "䷙", "䷄", "䷈", "䷡", "䷍", "䷪", "䷀",
		},
		[]string{
			"坤", "剥", "比", "观", "豫", "晋", "萃", "否",
			"谦", "艮", "蹇", "渐", "小过", "旅", "咸", "遁",
			"师", "蒙", "坎", "涣", "解", "未济", "困", "讼",
			"升", "蛊", "井", "巽", "恒", "鼎", "大过", "姤",
			"复", "颐", "屯", "益", "震", "噬嗑", "随", "无妄",
			"明夷", "贲", "既济", "家人", "丰", "离", "革", "同人",
			"临", "损", "节", "中孚", "归妹", "睽", "兑", "履",
			"泰", "大畜", "需", "小畜", "大壮", "大有", "夬", "乾",
		},
		[]int{2, 4, 8, 16, 32, 64},
	)
}
