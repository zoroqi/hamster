package main

import (
	"flag"
	"fmt"
	"math/bits"
	"math/rand/v2"
)

type Gua struct {
	Name string
	Img  string
	Num  uint8
}

func (g Gua) String() string {
	return fmt.Sprintf("%s(%s)", g.Name, g.Img)
}

func yao(y uint8) (r uint8, d bool) {
	switch y {
	case 0:
		r = 0
		d = true
	case 1, 2, 4:
		r = 1
		d = false
	case 3, 5, 6:
		r = 0
		d = false
	case 7:
		r = 1
		d = true
	}
	return
}

func yaoString(y uint8, d bool) (s string) {
	if y == 1 {
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

const yao_shu = 6

func yao_dong_gui_ke(gui_ke func(int32) int32) ([]uint8, []bool, uint8) {
	gua_xiang := uint8(0)
	yao_xiang := make([]uint8, yao_shu)
	yao_dong := make([]bool, yao_shu)
	for i := 0; i < yao_shu; i++ {
		yao_xiang[i], yao_dong[i] = yao(uint8(gui_ke(8)))
		gua_xiang = yao_xiang[i]<<i | gua_xiang
	}
	return yao_xiang, yao_dong, gua_xiang
}

func main() {
	flag.Parse()
	question := ""
	if flag.NArg() != 0 {
		question = flag.Arg(0)
	}

	yao_xiang, yao_dong, gua_xiang := yao_dong_gui_ke(rand.Int32N)

	gua_xiang_shang := gua_xiang >> 3
	gua_xiang_xia := gua_xiang & 0x07

	fmt.Println("你是周文王, 根据我摇卦的结果进行占卜.")
	fmt.Println("- 占卜的问题:", question)
	fmt.Println("- 卦象:\n```")
	// [本卦/错卦/综卦/复卦/象卦/交卦/变卦/杂卦解读](https://www.shuozhouyi.com/25248.html)
	fmt.Printf("本卦: %s\n", gua64[gua_xiang])
	fmt.Printf("上卦: %s 下卦: %s\n", gua8[gua_xiang_shang], gua8[gua_xiang_xia])
	fmt.Printf("错卦: %s\n", gua64[gua_xiang^uint8(0x3F)])
	fmt.Printf("宗卦: %s\n", gua64[bits.Reverse8(gua_xiang)>>2])
	fmt.Printf("互(复)卦: %s\n", gua64[((((gua_xiang_shang&0x03)<<1)|(gua_xiang_xia>>2))<<3)|((gua_xiang_shang&0x01<<2)|(gua_xiang_xia>>1))])
	for i := yao_shu; i > 0; i-- {
		i := i - 1
		fmt.Printf("%s爻: %s\n", yao_ming[i], yaoString(yao_xiang[i], yao_dong[i]))
	}
	fmt.Println("```")
}

var gua8 = map[uint8]Gua{}
var gua64 = map[uint8]Gua{}
var yao_ming = []string{"初", "二", "三", "四", "五", "上"}

func buildGuaStruct(img, txt []string, steps []int) map[uint8]Gua {
	kv := map[string]Gua{}
	for i, v := range img {
		kv[v] = Gua{
			Name: txt[i],
			Img:  v,
		}
	}
	for _, step := range steps {
		for i := 0; i < len(img); i++ {
			h := step / 2
			ba := kv[img[i]]
			if i%step < h {
				ba.Num = ba.Num<<1 | 0
			} else {
				ba.Num = ba.Num<<1 | 1
			}
			kv[img[i]] = ba
		}
	}
	r := map[uint8]Gua{}
	for _, v := range kv {
		r[v.Num] = v
	}
	return r
}

func init() {
	gua8 = buildGuaStruct(
		[]string{"☷", "☶", "☵", "☴", "☳", "☲", "☱", "☰"},
		[]string{"坤", "艮", "坎", "巽", "震", "离", "兑", "乾"},
		[]int{2, 4, 8},
	)

	gua64 = buildGuaStruct(
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
