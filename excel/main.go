package main

import (
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

var sheetName = []int{
	1,
	//2,
	//3,
	//4,
	//5,
	//6,
	//7,
	//8,
	//9,
	//10,
	//11,
	//12,
}
var (
	year  = flag.String("y", "", "year")
	month = flag.Int("m", 0, "month")
)

func main() {
	flag.Parse()
	if *year == "" || *month <= 0 || *month >= 13 {
		fmt.Println("err", *year, *month)
		return
	}
	fmt.Println(*year, *month)
	excelPath := "/Users/wuming/note/other/账本/历史账本"
	f, err := excelize.OpenFile(path.Join(excelPath, *year+"账本.xlsx"), excelize.Options{Password: "ziyouben"})
	if err != nil {
		fmt.Println(err)
		return
	}
	outfile := fmt.Sprintf("%02d-temp.bean", *month)
	str := parse(f)
	ioutil.WriteFile(path.Join("./", outfile), []byte(str), 0644)
}

const (
	Charge    = "Expenses:%s:Banking:Charge"
	History   = "Expenses:%s:History"
	Alipay    = "Expenses:%s:Insurance:Alipay"
	Clothing  = "Expenses:%s:Life:Clothing"
	Commuting = "Expenses:%s:Life:Commuting"
	Dine      = "Expenses:%s:Life:Dine"
	House     = "Expenses:%s:Life:House"
	Medicine  = "Expenses:%s:Life:Medicine"
	Snacks    = "Expenses:%s:Life:Snacks"
	Appliance = "Expenses:%s:Other:Appliance"
	Machine   = "Expenses:%s:Other:Machine"
	Sundries  = "Expenses:%s:Other:Sundries"
	Party     = "Expenses:%s:Party"
	GAME      = "Expenses:%s:Virtual:GAME"
	P0535     = "Expenses:%s:Virtual:Phone:0535"
	P3005     = "Expenses:%s:Virtual:Phone:3005"
	VIP       = "Expenses:%s:Virtual:VIP"
	Book      = "Expenses:%s:Learn:Book"
)

const dayLength = 31

var monthDay = make([][2]string, dayLength)

func init() {
	start, _ := excelize.ColumnNameToNumber("L")
	for i := 0; i < dayLength; i++ {
		desc, _ := excelize.ColumnNumberToName(start + 2*i)
		money, _ := excelize.ColumnNumberToName(start + 2*i + 1)
		monthDay[i] = [...]string{desc, money}
	}
}

func parse(f *excelize.File) string {
	sb := strings.Builder{}
	bank := extractBank(f)

	sb.WriteString(defaultParse(f, Dine, []int{3, 4, 5, 6, 7, 8}, "dine", []string{"dine"}, bank))
	sb.WriteString(defaultParse(f, Snacks, []int{10, 11}, "snacks", []string{"snacks"}, bank))
	sb.WriteString(defaultParse(f, Commuting, []int{12, 13}, "commuting", []string{"commuting"}, bank))
	sb.WriteString(defaultParse(f, Sundries, []int{14, 15, 16, 17}, "history", []string{"sundries"}, bank))
	sb.WriteString(defaultParse(f, Book, []int{18, 19}, "book", []string{"book"}, bank))
	sb.WriteString(defaultParse(f, Clothing, []int{20, 21, 22}, "clothing", []string{"clothing"}, bank))
	sb.WriteString(defaultParse(f, GAME, []int{23, 24, 25}, "history", []string{"game"}, bank))
	sb.WriteString(defaultParse(f, P0535, []int{26, 27}, "phone", []string{"phone"}, bank))
	sb.WriteString(defaultParse(f, Machine, []int{28, 29}, "machine", []string{"machine"}, bank))
	sb.WriteString(defaultParse(f, Medicine, []int{32, 33}, "medicine", []string{"medicine"}, bank))
	sb.WriteString(defaultParse(f, Sundries, []int{34, 35, 36, 37, 38, 39}, "history", []string{"other"}, bank))
	for k, v := range bank {
		if len(v) >0 {
			fmt.Println(k, v)
		}
	}
	return sb.String()
}

func defaultParse(f *excelize.File, expenses string,
	rowRange []int, defaultDesc string, tag []string, banksAccount map[string][]string) string {

	monthStr := fmt.Sprintf("%d月", *month)
	dine := extract(f, rowRange, monthStr)

	sb := strings.Builder{}
	for i, pairs := range dine {
		if len(pairs) <= 0 {
			continue
		}
		sb.WriteString(formatAccount(fmt.Sprintf(expenses, *year), i+1, defaultDesc, pairs, tag, banksAccount))
	}
	return sb.String()
}

func formatAccount(expenses string, day int, defaultDesc string, ps []pair, tag []string, banksAccount map[string][]string) string {
	sb := strings.Builder{}
	//desc := ""
	//for _, p := range ps {
	//	if p.desc == "" {
	//		continue
	//	}
	//	desc += p.desc + ","
	//}
	//if desc == "" {
	//	desc = defaultDesc
	//} else {
	//	desc = desc[:len(desc)-1]
	//}

	for _, p := range ps {
		desc := p.desc
		if desc == "" {
			desc = defaultDesc
		}
		sb.WriteString(fmt.Sprintf("%s-%02d-%02d * \"%s-history\" \"%s\"", *year, *month, day, *year, desc))
		for _, s := range tag {
			sb.WriteString(fmt.Sprintf(" #%s", s))
		}
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("  %s   +%.2f CNY\n", expenses, p.money))
		bank := banksAccount[fmt.Sprintf("%.2f", p.money)]
		b := A1
		if len(bank) != 0 {
			b = bank[0]
			banksAccount[fmt.Sprintf("%.2f", p.money)] = bank[1:]
		}
		sb.WriteString(fmt.Sprintf("  %s\n\n", b))
	}
	return sb.String()
}

type pair struct {
	desc  string
	money float64
}

func extract(f *excelize.File, rowRange []int, month string) [][]pair {
	pairs := make([][]pair, dayLength)
	for i, day := range monthDay {
		n := make([]pair, 0, 3)
		for _, r := range rowRange {
			desc, _ := f.GetCellValue(month, fmt.Sprintf("%s%d", day[0], r))
			money, _ := f.GetCellValue(month, fmt.Sprintf("%s%d", day[1], r))
			if strings.TrimSpace(money) == "" {
				continue
			}
			m := replaceDollar(money)
			if m == 0 {
				continue
			}
			n = append(n, pair{desc: strings.TrimSpace(desc), money: m})
		}
		if len(n) > 0 {
			pairs[i] = n
		}
	}
	return pairs

}

func extractBank(f *excelize.File) map[string][]string {
	monthStr := fmt.Sprintf("%d月", *month)
	r := make(map[string][]string)
	for i := 31; i < 87; i++ {
		bank, _ := f.GetCellValue(monthStr, fmt.Sprintf("B%d", i))

		if nameMapper[bank] == "" {
			continue
		}
		accontes, _ := f.GetCellFormula(monthStr, fmt.Sprintf("E%d", i))
		if accontes == "" {
			accontes, _ = f.GetCellValue(monthStr, fmt.Sprintf("E%d", i))
		}

		var money []rune
		moneys := make([]string, 0, 0)
		for _, r := range []rune(accontes) {
			if r == '+' || r == '-' {
				moneys = append(moneys, string(money))
				money = money[:0]
			}
			money = append(money, r)
		}
		moneys = append(moneys, string(money))

		for _, m := range moneys {
			f, e := strconv.ParseFloat(m, 64)
			if e != nil {
				fmt.Println(f, e)
				continue
			}
			f = -f
			r[fmt.Sprintf("%0.2f", f)] = append(r[fmt.Sprintf("%0.2f", f)], nameMapper[bank])
		}
	}

	return r
}

func main2() {
	excelPath := "/Users/wuming/note/other/账本/历史账本"
	f, err := excelize.OpenFile(path.Join(excelPath, "2020账本.xlsx"), excelize.Options{Password: "ziyouben"})
	if err != nil {
		return
	}
	output := "/Users/wuming/note/historybean/2020"

	for _, v := range sheetName {
		s := format(f, fmt.Sprintf("%d月", v), v)
		ioutil.WriteFile(path.Join(output, fmt.Sprintf("%02d.bean", v)), []byte(s), 0644)
	}
}

func format(f *excelize.File, month string, num int) string {
	balanceDate := fmt.Sprintf("%s-%02d-01", year, num)
	countDate := fmt.Sprintf("%s-%02d-02", year, num)

	alipay, am := account(f, month, 19)
	beiyin, bm := account(f, month, 22)
	zhaoshangxin, zxm := account(f, month, 24)
	zhaoshang, zm := account(f, month, 20)
	huabei, hm := account(f, month, 25)
	zhongguo, zgm := account(f, month, 26)
	weixin, wm := account(f, month, 27)

	sb := strings.Builder{}
	cash, _ := f.GetCellValue(fmt.Sprintf("%d月", num-1), "G16")

	sb.WriteString(balance(balanceDate, "现金", replaceDollar(cash)))
	sb.WriteString(balance(balanceDate, alipay, am))
	sb.WriteString(balance(balanceDate, beiyin, bm))
	sb.WriteString(balance(balanceDate, zhaoshang, zm))
	sb.WriteString(balance(balanceDate, zhaoshangxin, zxm-60000))
	sb.WriteString(balance(balanceDate, huabei, hm-32000))
	sb.WriteString(balance(balanceDate, zhongguo, zgm))
	sb.WriteString(balance(balanceDate, weixin, wm))

	sb.WriteString(liushui(countDate))

	return sb.String()
}

func liushui(countDate string) string {
	sb := strings.Builder{}
	sb.WriteString("\n\n\n")

	sb.WriteString(fmt.Sprintf("%s * \"优酷\" \"工资\" #incom\n", countDate))
	sb.WriteString(fmt.Sprintf("  %s   \n", I1))
	sb.WriteString(fmt.Sprintf("  %s   \n", A2))
	sb.WriteString(fmt.Sprintf("\n\n"))

	sb.WriteString(fmt.Sprintf(`
%s * "个人" "取现" #cash
  Assets:Cash

%s * "history" "history"
  Expenses:History

%s * "history" "history"
  Expenses:History  0 CNY
  Assets:Cash

%s * "history" "还款" #bank #credit_card
  Liabilities:Bank:招商银行:8336 

%s * "alipay" "还款" #bank #credit_card
  Liabilities:Web:Alipay:花呗
   
`, countDate, countDate, countDate, countDate, countDate))

	return sb.String()
}

const (
	I1  = "Income:Work:凤凰网"
	I2  = "Income:Bank:招商银行:聚益生金30"
	I3  = "Income:Bank:招商银行:聚益生金63"
	I4  = "Income:Bank:招商银行:日日欣"
	I5  = "Income:Gov:个人医保"
	I6  = "Income:Web:Alipay:余额宝"
	I7  = "Income:Other:LuckyMoney"
	I8  = "Income:Bank:招商银行:定期"
	I9  = "Income:Work:外包公司"
	I10 = "Income:Work:顶尖时代"
	I11 = "Income:Work:百分点"
	I12 = "Income:Work:优酷"

	A1 = "Assets:Cash"
	A2 = "Assets:Bank:招商银行:0829"
	A3 = "Assets:Bank:中国银行:1933"
	A4 = "Assets:Web:Alipay:余额宝"
	A5 = "Assets:Web:Alipay:其他"
	A6 = "Assets:Web:Wechat:Zoroqi"
	A7 = "Assets:Bank:北京银行:1760"
	A8 = "Assets:Bank:农业银行:3515"
)

var nameMapper = map[string]string{
	"中国银行信用卡":       "Liabilities:Bank:中国银行:5141",
	"支付宝":           "Assets:Web:Alipay:余额宝",
	"招商银行 信用卡":      "Liabilities:Bank:招商银行:8336",
	"招商银行-公积金-0829": "Assets:Bank:招商银行:0829",
	"招商银行-公积金":      "Assets:Bank:招商银行:0829",
	"蚂蚁花呗":          "Liabilities:Web:Alipay:花呗",
	"中国银行1933":      "Assets:Bank:中国银行:1933",
	"北银1760":        "Assets:Bank:北京银行:1760",
	"微信":            "Assets:Web:Wechat:Zoroqi",
	"农行3515":        "Assets:Bank:农业银行:3515",
	"现金":            A1,
}

func balance(date string, name string, money float64) string {
	accountName, exist := nameMapper[name]
	if !exist {
		return ""
	}
	return fmt.Sprintf("%s balance %s %0.2f CNY\n", date, accountName, money)
}

func account(f *excelize.File, month string, row int) (name string, money float64) {
	name, _ = f.GetCellValue(month, fmt.Sprintf("B%d", row))
	moneyS, _ := f.GetCellValue(month, fmt.Sprintf("H%d", row))
	money = replaceDollar(moneyS)
	return
}

func replaceDollar(moneyS string) float64 {
	f, _ := strconv.ParseFloat(strings.ReplaceAll(moneyS, "￥", ""), 64)
	return f
}
