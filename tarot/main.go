package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
)

type direction bool

func (d direction) String() string {
	if d {
		return "正位"
	}
	return "逆位"
}

type card struct {
	name      string // english
	zhName    string // chinese
	number    int
	romanNum  string
	direction direction
}

type spread struct {
	name      string
	zhname    string
	summary   string
	suited    string
	positions []string
}

func main() {

	spreadNameHelp := ""
	for i, s := range spreads {
		spreadNameHelp += fmt.Sprintf("%d: %s\n", i+1, s.zhname)
	}

	spreadNum := flag.Int("s", 1, "占卜使用的牌阵 (1-34)\n"+spreadNameHelp)

	spreadHelp := flag.Int("spread", 0, "显示牌阵帮助 (1-34)")
	flag.Parse()
	if *spreadHelp > 0 {
		fmt.Println("牌阵帮助:")
		s := spreads[*spreadHelp-1]
		fmt.Printf("%d:%s\n", *spreadHelp-1, s.zhname)
		fmt.Printf("适合场景:%s\n", s.suited)
		fmt.Printf("简介:%s\n", s.summary)
		fmt.Println("牌阵位置:")
		for i, p := range s.positions {
			fmt.Printf("%d: %s\n", i+1, p)
		}
		return
	}

	question := ""
	if flag.NArg() != 0 {
		question = flag.Arg(0)
	}

	if *spreadNum < 1 || *spreadNum > len(spreads) {
		fmt.Println("无效的牌阵选择，请输入 1 到 34 的数字。")
		return
	}
	*spreadNum = (*spreadNum) - 1
	tarotCards := []card{}
	tarotCards = append(tarotCards, majorArcana...)
	tarotCards = append(tarotCards, minorArcanaWands...)
	tarotCards = append(tarotCards, minorArcanaCups...)
	tarotCards = append(tarotCards, minorArcanaSwords...)
	tarotCards = append(tarotCards, minorArcanaPentacles...)
	shuffle := func(i, j int) {
		tarotCards[i], tarotCards[j] = tarotCards[j], tarotCards[i]
		tarotCards[i].direction = !tarotCards[i].direction
		tarotCards[j].direction = !tarotCards[j].direction
	}
	for range 5 {
		rand.Shuffle(len(tarotCards), shuffle)
	}

	fmt.Println("欢迎来到塔罗牌占卜程序!")

	position := spreads[*spreadNum].positions
	drawnCards := make([]card, len(position))
	for i := 0; i < len(position); {
		fmt.Printf("请选择一张牌 (输入 1 到 %d 的数字):\n", len(tarotCards))
		var choice int
		fmt.Scan(&choice)

		if choice < 1 || choice > len(tarotCards) {
			fmt.Printf("无效的选择，请输入 1 到 %d 的数字。", len(tarotCards))
			continue
		}

		card := tarotCards[choice-1]
		tarotCards = append(tarotCards[:choice-1], tarotCards[choice:]...)
		drawnCards[i] = card
		i++
	}
	fmt.Println("占卜问题:", question)
	fmt.Println("牌阵:", spreads[*spreadNum].zhname)
	fmt.Println("占卜结果:")
	for i, card := range drawnCards {
		fmt.Printf("%d. %s: %s (%s)\n", i+1, position[i], card.zhName, card.direction)
	}
}

var majorArcana = []card{
	{"The Fool", "愚者", 0, "Ο", true},
	{"The Magician", "魔术师", 1, "I", true},
	{"The High Priestess", "女祭司", 2, "II", true},
	{"The Empress", "女皇", 3, "III", true},
	{"The Emperor", "皇帝", 4, "IV", true},
	{"The Hierophant", "教皇", 5, "V", true},
	{"The Lovers", "恋人", 6, "VI", true},
	{"The Chariot", "战车", 7, "VII", true},
	{"Strength", "力量", 8, "VIII", true},
	{"The Hermit", "隐者", 9, "IX", true},
	{"Wheel of Fortune", "命运之轮", 10, "X", true},
	{"Justice", "正义", 11, "XI", true},
	{"The Hanged Man", "倒吊人", 12, "XII", true},
	{"Death", "死神", 13, "XIII", true},
	{"Temperance", "节制", 14, "XIV", true},
	{"The Devil", "恶魔", 15, "XV", true},
	{"The Tower", "塔", 16, "XVI", true},
	{"The Star", "星星", 17, "XVII", true},
	{"The Moon", "月亮", 18, "XVIII", true},
	{"The Sun", "太阳", 19, "XIX", true},
	{"Judgement", "审判", 20, "XX", true},
	{"The World", "世界", 21, "XXI", true},
}

var minorArcanaCups = []card{
	{"Ace of Cups", "圣杯王牌", 22, "I", true},
	{"Two of Cups", "圣杯二", 23, "II", true},
	{"Three of Cups", "圣杯三", 24, "III", true},
	{"Four of Cups", "圣杯四", 25, "IV", true},
	{"Five of Cups", "圣杯五", 26, "V", true},
	{"Six of Cups", "圣杯六", 27, "VI", true},
	{"Seven of Cups", "圣杯七", 28, "VII", true},
	{"Eight of Cups", "圣杯八", 29, "VIII", true},
	{"Nine of Cups", "圣杯九", 30, "IX", true},
	{"Ten of Cups", "圣杯十", 31, "X", true},
	{"Page of Cups", "圣杯侍从", 32, "XI", true},
	{"Knight of Cups", "圣杯骑士", 33, "XII", true},
	{"Queen of Cups", "圣杯皇后", 34, "XIII", true},
	{"King of Cups", "圣杯国王", 35, "XIV", true},
}

var minorArcanaPentacles = []card{
	{"Ace of Pentacles", "星币王牌", 36, "I", true},
	{"Two of Pentacles", "星币二", 37, "II", true},
	{"Three of Pentacles", "星币三", 38, "III", true},
	{"Four of Pentacles", "星币四", 39, "IV", true},
	{"Five of Pentacles", "星币五", 40, "V", true},
	{"Six of Pentacles", "星币六", 41, "VI", true},
	{"Seven of Pentacles", "星币七", 42, "VII", true},
	{"Eight of Pentacles", "星币八", 43, "VIII", true},
	{"Nine of Pentacles", "星币九", 44, "IX", true},
	{"Ten of Pentacles", "星币十", 45, "X", true},
	{"Page of Pentacles", "星币侍从", 46, "XI", true},
	{"Knight of Pentacles", "星币骑士", 47, "XII", true},
	{"Queen of Pentacles", "星币皇后", 48, "XIII", true},
	{"King of Pentacles", "星币国王", 49, "XIV", true},
}

var minorArcanaSwords = []card{
	{"Ace of Swords", "宝剑王牌", 50, "I", true},
	{"Two of Swords", "宝剑二", 51, "II", true},
	{"Three of Swords", "宝剑三", 52, "III", true},
	{"Four of Swords", "宝剑四", 53, "IV", true},
	{"Five of Swords", "宝剑五", 54, "V", true},
	{"Six of Swords", "宝剑六", 55, "VI", true},
	{"Seven of Swords", "宝剑七", 56, "VII", true},
	{"Eight of Swords", "宝剑八", 57, "VIII", true},
	{"Nine of Swords", "宝剑九", 58, "IX", true},
	{"Ten of Swords", "宝剑十", 59, "X", true},
	{"Page of Swords", "宝剑侍从", 60, "XI", true},
	{"Knight of Swords", "宝剑骑士", 61, "XII", true},
	{"Queen of Swords", "宝剑皇后", 62, "XIII", true},
	{"King of Swords", "宝剑国王", 63, "XIV", true},
}

var minorArcanaWands = []card{
	{"Ace of Wands", "权杖王牌", 64, "I", true},
	{"Two of Wands", "权杖二", 65, "II", true},
	{"Three of Wands", "权杖三", 66, "III", true},
	{"Four of Wands", "权杖四", 67, "IV", true},
	{"Five of Wands", "权杖五", 68, "V", true},
	{"Six of Wands", "权杖六", 69, "VI", true},
	{"Seven of Wands", "权杖七", 70, "VII", true},
	{"Eight of Wands", "权杖八", 71, "VIII", true},
	{"Nine of Wands", "权杖九", 72, "IX", true},
	{"Ten of Wands", "权杖十", 73, "X", true},
	{"Page of Wands", "权杖侍从", 74, "XI", true},
	{"Knight of Wands", "权杖骑士", 75, "XII", true},
	{"Queen of Wands", "权杖皇后", 76, "XIII", true},
	{"King of Wands", "权杖国王", 77, "XIV", true},
}

// 这部分信息来自 [AI塔罗牌：免费在线AI塔罗牌占卜](https://tarotap.com/).

var spreads = []spread{
	{name: "三张牌占卜法", zhname: "三张牌占卜法", summary: "三张牌占卜法是一种通用型牌阵，适用于多种场合和情况。它可以自由定义用于分析独立事件的不同方面，或用来占卜各种相关事物。这个万用型牌阵简单易用，适合新手，但同时也能直指问题核心。通过解读三张牌的关系，可以获得全面而深入的洞察。", suited: " 适合综合分析 & 单事解读", positions: []string{"第一张牌", "第二张牌", "第三张牌"}},
	{name: "时间流牌阵", zhname: "时间流牌阵", summary: "时间流牌阵是一种平行流向的时间解析法，将时间流应用于空间维度的占卜。它仿佛将流动的时间从过去延伸到未来，让事件铺陈于时间毯之上。这个牌阵特别适合有时间指向的占卜，是解读未来的有力工具。虽然不适合深层次占卜，但可以在时间轨道上占卜各种事物，因为任何事都具有时间属性。", suited: "适合预测未来 & 洞察未知", positions: []string{"过去", "现在", "未来"}},
	{name: "五张牌占卜法", zhname: "五张牌占卜法", summary: "五张牌占卜法是一种灵活的通用型占卜牌阵，不设置具体的位置意涵，可以根据需要自由定义。这个牌阵可以用在很多场合，比三张牌占卜法更深入具体，特别适合用来分析独立事情的某个方面。它适合较复杂的占卜，通过对每张牌的定义，我们可以探索复杂事物的不同方面，深入挖掘问题本质。五张牌占卜法的灵活性使它成为一个强大的工具，尤其在其他牌阵无法涵盖你全部问题的时候，这个自定义牌阵可以发挥重要作用，帮助你获得全面而深入的洞察。", suited: " 适合综合分析 & 单事解读", positions: []string{"万用牌阵第一张", "万用牌阵第二张", "万用牌阵第三张", "万用牌阵第四张", "万用牌阵第五张"}},
	{name: "事业金字塔阵", zhname: "事业金字塔阵", summary: "事业金字塔阵是一个专门用于解析财富运势和提升财商的牌阵。它能帮助我们明晰自我的优缺点，对事业的稳固发展有很大帮助。这个牌阵特别适合在事业遇到瓶颈期时使用，能为我们找到突破口。通过分析内核竞争力、个人优缺点以及最终可能达到的成就，事业金字塔阵为那些在事业中寻求成功的人提供了清晰的指引。它帮助我们除弊趋利，在事业中获得优势。如果你想要深入分析自己的事业前景，事业金字塔阵可以为你提供清晰明了的洞察。", suited: " 适合解析财富运势 & 提升财商", positions: []string{"你的内核竞争力", "代表你的缺点", "代表你的优点", "你最终的成就"}},
	{name: "凯尔特牌阵", zhname: "凯尔特牌阵", summary: "凯尔特牌阵是一个经久不衰的古老塔罗牌阵，被誉为永恒经典、备受推崇的塔罗牌阵。它通过对事物的层层剖析，抽丝剥茧，帮助我们做出明智的抉择。这个牌阵拥有很强的总结和规划能力，结构细致严谨，能够从宏观角度审视事件全貌，协助我们做出有利的决策。凯尔特牌阵不仅能分析问题的现状和过去，还能预测未来的发展，同时考虑到个人状况、环境影响以及希望和恐惧等多个方面。如果你想体验塔罗占卜的经典魅力，凯尔特牌阵绝对是不容错过的选择。", suited: "适合事件整体分析 & 宏观审视", positions: []string{"问题的现状", "目前的阻碍或麻烦", "你对问题的理想或目标", "造成现在状况的原因", "问题最近的过去状况", "问题最近的未来发展", "你本身的状况", "当前周边的环境", "你的希望和恐惧", "问题的最终结果"}},
	{name: "直指内核牌阵", zhname: "直指内核牌阵", summary: "直指内核牌阵专门用于占卜问题的核心因素，特别适用于具体问题遇到瓶颈时寻求突破。这个牌阵能帮助我们清楚地看到问题的根本所在，进而做出正确的决策。它有着极强的问题解决能力，可以快速找到问题的症结。如果你对某个问题感到纠结不定，这个牌阵可以帮你快速理清思路，找到解决方案。", suited: " 适合问题探索 & 切中要害", positions: []string{"问题的内核", "障碍或短处", "问题的对策", "资源或长处"}},
	{name: "上半年运势阵", zhname: "上半年运势阵", summary: "上半年运势阵是专门用来预测上半年运势的牌局。它可以针对特定方面进行独立解析，比如爱情半年运、事业半年运等。这个专用牌阵不仅可以占卜上半年的整体运势，还能分析爱情和财运的走势。由于时间的严格限定，本牌阵只适合占卜上半年运。通过逐月分析，上半年运势阵能帮助我们对未来半年有清晰的预见，为我们的生活和工作提供valuable指导。", suited: " 用于占卜半年运程 & 上半年运", positions: []string{"上半年一月运势", "上半年二月运势", "上半年三月运势", "上半年四月运势", "上半年五月运势", "上半年六月运势"}},
	{name: "四元素牌阵", zhname: "四元素牌阵", summary: "四元素牌阵通过感性、理性、物质、行动四个方面全面审视问题，帮助我们深入了解问题的实质。这个牌阵适合探索我们对世界的认知，从多角度审视固定的问题或单一事物。它提醒我们，只有全面地了解一个问题，才能真正解决它。四元素牌阵为我们提供了一个全面、系统的问题分析框架，有助于做出更明智的决策。", suited: "适合问题探索 & 多方解析", positions: []string{"行动与信心的暗示", "现实与金钱的暗示", "理性与决策的暗示", "感情或感性的暗示"}},
	{name: "金三角牌阵", zhname: "金三角牌阵", summary: "金三角牌阵是一种事业解析牌形，广受职场人士欢迎。它特别适用于事业发展遇到瓶颈时寻求突破。这个牌阵能帮助我们清楚地看到整个问题的来龙去脉，从而做出正确的决策。它对问题进行连贯的解剖，理清事物的前因后果，给出明确的思路。金三角牌阵以其一气呵成的特点，能有效帮助你寻找财富的方向，是解决职场和财务问题的有力工具。", suited: "适合财富探索 & 问题解析", positions: []string{"你现在的处境", "你遇到的困扰、问题", "问题对未来的影响", "你的解决方式或最终结果"}},
	{name: "吉普赛牌阵", zhname: "吉普赛牌阵", summary: "吉普赛牌阵是恋爱中人的专用牌阵，对双方各自进行剖析，比较相处方式和环境适应能力。它能总结恋爱中遇到的问题，并对未来做出预判。这个牌阵适用于婚姻、爱情、恋爱等各种感情方面的占卜，帮助探索彼此内心想法，找到合适的相处方式。浪漫而奔放的吉普赛牌阵是释放情感困扰的首选，能为你的感情生活带来新的洞察。", suited: "适合情侣分析 & 关系延展", positions: []string{"对方目前的想法", "自己目前的状况", "与对方相处应采取的方式", "目前的周遭状况", "关系最后的结果"}},
	{name: "圣三角牌阵", zhname: "圣三角牌阵", summary: "圣三角牌阵是时间流的一种变形，更注重事物的内在原因而非单纯的时间流向。它特别适合梳理问题的前因后果，理清事情的来龙去脉。这个牌阵常用于单项独立事物的占卜，能够简洁扼要地展示前因后果，清晰明了。在研究事情脉络，针对需要审证求因的问题时，这个牌阵尤其有效。", suited: "适合判断情势 & 寻找成因", positions: []string{"过去的原因", "问题的现状", "将来的结果"}},
	{name: "灵感对应牌阵", zhname: "灵感对应牌阵", summary: "灵感对应牌阵擅长占卜情绪链接或情感交互，具有简约直观的特点，常常有着神奇的应验效果。这个牌阵在情侣、同事、朋友间的应用较为广泛。它是一个双向的占卜牌阵，可以让我们看到同层次的交融，特别适合婚姻、爱情、恋爱等感情方面的交互关系占卜。当你需要揣摩对方心思，或者想要了解双方对关系的看法和期望时，灵感对应牌阵是上佳之选。它能帮助我们更好地理解彼此，促进关系的和谐发展。", suited: " 用于占卜情感因应 & 人际关系", positions: []string{"自己对对方的看法", "对方对自己的看法", "自己认为目前的关系", "对方认为目前的关系", "自己期望将来的发展", "对方期望将来的发展"}},
	{name: "面试求职牌阵", zhname: "面试求职牌阵", summary: "面试求职牌阵是专门用于求职新工作解析的牌阵。它能帮助我们洞悉对方的需求，了解自己需要注意的地方，有效提高面试成功率。这个牌阵最适合在求职之前使用，是预测面试求职结果的不二之选。它可以帮助求职者发现潜在的问题并找到解决方案。特别是在揣摩面试官心思时，这个牌阵往往能发挥神奇的效果。通过分析自己的心态、需要注意的情况、可能发生的状况以及对方的要求，面试求职牌阵为求职者提供了全面的指导。", suited: " 占卜应征面试 & 求职状况", positions: []string{"自己的心态及想法", "面试前需要注意的情况", "面试时将要发生的状况", "对方的要求或者问题", "最后的结果"}},
	{name: "爱情大十字", zhname: "爱情大十字", summary: "爱情大十字牌阵注重内心情感，主要应用于情侣之间。它善于洞悉彼此关系中的情感状况并分析结果。这个牌阵布局合理清晰，特别适合用于感情心理的占卜，能为感情生活带来实质性的帮助。通过分析双方的心理状态、当前关系状况以及外部环境因素，爱情大十字牌阵能全面地揭示感情问题的各个方面，帮助你找到解决感情困扰的方法。", suited: "解析两性关系 & 爱情状况", positions: []string{"自己目前的心情及想法", "对方现在的心理及态度", "彼此现在的状况", "目前周遭的环境情况", "彼此关系最后的结果"}},
	{name: "爱情树牌阵", zhname: "爱情树牌阵", summary: "爱情树牌阵是一个强大的工具，用于解析爱情关系的前因后果，回溯爱情过往。它不仅能探究感情问题的本源内核，还能指引未来趋势，为陷入爱情困境的情侣提供改善建议。这个牌阵特别适合在恋爱遇到困境时使用，能帮助寻找潜在的问题原因，从而改善感情关系。通过分析自己的想法、过去的原因、现在的建议、未来的指向以及潜在的影响，爱情树牌阵为解决复杂的感情问题提供了全面的视角。", suited: "适合溯本求源 & 寻找症结", positions: []string{"自己的想法", "过去的原因", "现在的建议", "未来的指向", "潜在的影响"}},
	{name: "恋人金字塔", zhname: "恋人金字塔", summary: "恋人金字塔牌阵简洁直接，涵盖了两人相恋的原始要素。它适合恋人情侣间的占卜，牌面一出就能明了易懂。这个牌阵的要素明确，条理清晰，不仅能够分析当前的关系状况，还能预测未来的发展趋势。当两个人之间遇到矛盾或迷惘时，这个牌阵可以帮助洞察问题，预见未来的走向，是解决感情问题的有力工具。", suited: "适合占卜恋人关系 & 交互解析", positions: []string{"代表你自己", "代表你的恋人", "你们彼此的关系", "你们未来的发展"}},
	{name: "恋人树牌阵", zhname: "恋人树牌阵", summary: "恋人树牌阵能够详细分析对方的心理和态度，帮助我们找到存在于两人之间的症结。通过仔细揣摩利弊细节，这个牌阵能够帮助我们做出正确的选择。它采用树形分支的结构，可以清晰地对比情侣间的行为差异，进而根据事实做出判断。如果你想了解对方在想什么，以及自己在爱情中的真实感受，恋人树牌阵是一个理想的选择。它不仅能帮助我们理解双方的想法和态度，还能揭示彼此对关系的期望，以及最终的发展趋势。", suited: "适合探究对方心理 & 行为模式", positions: []string{"自己对对方的想法与态度", "对方对自己的想法与态度", "自己对彼此关系的期望", "对方对彼此关系的期望", "自己的外在行动", "对方的外在行动", "两人最后的结果"}},
	{name: "婚姻牌阵", zhname: "婚姻牌阵", summary: "婚姻牌阵专门针对婚姻状况和期望进行解析。通过牌面分析姻缘，它能帮助我们掌握婚姻走势，并根据现实情况推演出最后的结果。这个牌阵特别适合婚姻专项的占卜，能够探明姻缘脉动，深入解析婚姻关系。无论是在婚前还是婚后，当你需要了解婚姻状况时，婚姻牌阵都是一个很好的选择。它不仅能帮助我们理解当前的婚姻状况，还能洞察未来的发展趋势，为我们的婚姻生活提供指引。", suited: " 适合结婚运势 & 婚姻剖析", positions: []string{"代表你自己", "过去的状况", "现在的状况", "遇到的问题", "对婚姻的期望", "对婚姻的恐惧", "未来的发展"}},
	{name: "身心灵牌阵", zhname: "身心灵牌阵", summary: "身心灵牌阵通过从灵性、心理、身体三个方面审视自我，帮助我们更全面地了解自身状况。这个牌阵特别适合在自己感到迷惘时使用，通过向内寻求答案来获得洞察。它能帮助我们探索自己的来处和去向，让塔罗为我们指明方向。身心灵牌阵是一个强大的自我探索工具，能带来深刻的自我认知。", suited: "适合自我探索 & 了解自己", positions: []string{"身体的状况", "心理的状况", "灵魂的状况", "可学习提升的元素"}},
	{name: "寻找对象牌阵", zhname: "寻找对象牌阵", summary: "寻找对象牌阵专门用于未来爱人的占卜，适合寻找理想的伴侣。这个牌阵不仅能帮助单身人士创建意中人的愿景，还能帮助确定自己的目标。它鼓励我们先通过塔罗的指引做好准备，然后再出发寻找真爱。如果你不确定自己想要什么样的伴侣，这个牌阵可以帮你梳理思路，明确期望。通过分析自己的现状、理想对象的特质以及可能的行动方向，寻找对象牌阵为寻找真爱提供了全面的指导。", suited: " 适合寻找意中人 & 有缘人", positions: []string{"代表你现在的心情、处境", "代表你希望追求的对象", "代表你不喜欢的对象", "代表该采取的行动", "代表未来发展、最后结果"}},
	{name: "问题解决牌阵", zhname: "问题解决牌阵", summary: "问题解决牌阵是一个专门用于解决具体问题的强大工具。通过对问题的深入剖析，对比前因后果，根据逻辑关系，这个牌阵能够帮助提高处理问题的能力。它不仅能总结问题、评判损益得失，还能对未来做出预判。这个牌阵对解决问题有积极的帮助，可以深入挖掘问题的本质。如果你面临一个棘手的问题，或者需要做出重要决策，问题解决牌阵可以为你提供全面的分析和洞察。", suited: " 适合问题剖析 & 答疑解惑", positions: []string{"问题发生的原因", "问题现在的状况", "周遭的环境情况", "将会遇到的阻碍", "问题的解决方式"}},
	{name: "恋人复合牌阵", zhname: "恋人复合牌阵", summary: "恋人复合牌阵专门针对分手后是否复合进行占卜。通过详细分析牌面，这个牌阵能帮我们了解自己的状况和前任的想法，理清彼此的关系。它提醒我们，得不到的永远在骚动，被偏爱的都有恃无恐。通过对照彼此的内心感受，恋人复合牌阵能帮我们揭开对方扑朔迷离的面纱。如果你仍然念念不忘，觉得彼此前缘未了，这个牌阵可以帮你全面分析情况，做出明智的选择。它不仅能揭示过去和现在的状况，还能预测未来的发展，是处理复杂感情问题的有力工具。", suited: " 适合分手复合 & 前任想法", positions: []string{"你和他的过去", "你现在的状况", "他现在的状况", "你对复合的感受", "他对复合的感受", "阻碍你的", "帮助你的", "你不知道的重要事", "整体结果"}},
	{name: "下半年运势阵", zhname: "下半年运势阵", summary: "下半年运势阵是专门用来预测下半年运势的牌局。与上半年运势阵类似，它可以针对特定方面进行独立解析，如爱情半年运、事业半年运等。这个专用牌阵不仅可以占卜下半年的整体运势，还能分析爱情和财运的走势。由于时间的严格限定，本牌阵只适合占卜下半年运。通过逐月分析，下半年运势阵能帮助我们对未来半年有清晰的预见，为我们的生活和工作提供valuable指导。", suited: " 用于占卜半年运程 & 下半年运", positions: []string{"下半年一月运势", "下半年二月运势", "下半年三月运势", "下半年四月运势", "下半年五月运势", "下半年六月运势"}},
	{name: "自我探索牌阵", zhname: "自我探索牌阵", summary: "自我探索牌阵是一种自我修为成长牌型，专门用于认识自我、开发自身潜力，提高对自己的认识。当我们遇到成长瓶颈需要寻求突破时，这个牌阵能帮助我们从潜意识底层发现真实的自己。它用于提升自我的觉知，可以突破自身瓶颈，具有极强的领悟能力。自我探索牌阵通过分析我们的内心深处、精神生活、知识领域和感情生活，为我们提供了全面的自我认知。如果你想体验醍醐灌顶的感觉，用心去感悟这个牌阵，它将为你打开自我认知的新世界。", suited: " 适合认识自我 & 提升潜能", positions: []string{"代表内心深处、内在的自我", "代表精神生活、灵性面", "代表知识领域、理性面", "代表感情生活、感情面"}},
	{name: "单张牌占卜法", zhname: "单张牌占卜法", summary: "单张牌占卜法简单易懂，擅长回答是非问题和预测单日运势。抽取一张牌即可做出判断，是塔罗入门者的理想选择。它不仅可以用于日常快速占卜，还可以在网络占卜中作为切牌或补牌使用。每天清晨抽一张塔罗牌，能让你快速了解当天的整体运势走向，为你的一天做好准备。", suited: " 适合是非判断 & 单日运势", positions: []string{"无牌阵单张"}},
	{name: "六芒星牌阵", zhname: "六芒星牌阵", summary: "六芒星牌阵是一个意义深远的塔罗牌阵，专门用于指向未来和预判事物的发展方向。这个牌阵能帮助我们理清事情的本源，是真正拥有预测属性的塔罗牌阵。它不仅能判断事情的走向，还能分析潜意识与显意识的表达，具有极强的窥视未来的能力。六芒星牌阵对事物发展的预测有着积极的指导意义。如果你想真正地预测未来，深入分析和了解六芒星牌阵将会给你带来全面而深刻的洞察。", suited: "用于占卜事物发展 & 预测未来", positions: []string{"问题过去状况", "问题现在状况", "问题未来状况", "解决问题的对应策略", "周遭的环境状况", "本人的心理态度", "事物的最终结果"}},
	{name: "真命天子牌阵", zhname: "真命天子牌阵", summary: "真命天子牌阵帮助我们探索情投意合的人。它通过引导我们探索内心想法，改善外在行为，来增加遇到适合之人的机会。这个牌阵特别适合单身人士用来占卜邂逅有缘人的可能性。通过内在和外在的双重改观，它能帮助我们找到真正合适的伴侣。真命天子牌阵不仅帮助我们了解理想伴侣的类型，还引导我们思考需要做出的改变，以及应该相信的事情。如果你正在寻找真爱，这个牌阵可能会为你带来意想不到的洞察和惊喜。", suited: " 单身探索有缘人 & 合意对象", positions: []string{"我真正的伴侣的类型", "我真正的伴侣已经走入我的生命中了吗？", "将会有困难产生吗？", "什么样的改变是需要的？", "我将相信什么？"}},
	{name: "爱之星牌阵", zhname: "爱之星牌阵", summary: "爱之星牌阵是专门用于占卜爱情的塔罗牌阵。它可以深入分析两个人的心理状况、对未来的期望，以及两性关系的最终结果。这个牌阵特别适合在遇到感情问题时使用，能帮助理清感情问题的前因后果，促进关系的融洽。当两个人的关系走到十字路口时，爱之星牌阵可以帮助我们回顾过去，理解现在，预见未来。通过全面的分析，它能为爱情关系提供新的洞察和指引，帮助我们做出明智的决定。", suited: "适合占卜爱情 & 两性关系", positions: []string{"目前的状况", "女方的心情", "男方的心情", "你们过去的状况", "你们对未来的期望", "你们的最终结果"}},
	{name: "三选一牌阵", zhname: "三选一牌阵", summary: "三选一牌阵专门针对三个不同的选项进行占卜。通过详细分析牌面，这个牌阵能帮助我们对各个选项的利害关系进行深入的分析判断，从而权衡利弊，做出最终的选择。它特别适合在较为复杂的环境中使用，能帮助我们做出正确的判断，清晰地分辨各方的优劣。当你面临多个选择，尤其是在选项较多时，三选一牌阵可以提供全面而深入的指引，帮助你做出最佳决策。", suited: "适合事情抉择 & 选择占卜", positions: []string{"求问者本身状况", "代表选择A的发展", "代表选择B的发展", "代表选择C的发展", "代表选择A的结果", "代表选择B的结果", "代表选择C的结果"}},
	{name: "财富之树", zhname: "财富之树", summary: "财富之树是一个专门用来解析财富运势生长的牌阵。它不仅用于占卜事业和财运，还通过对财富的梳理，帮助创建适合自己的财富模式。这个牌阵象征财富的生成过程，可以揭示财富脉搏，对求财有积极的指导意义。如果你想了解自己的财富指数，或者寻求事业发展的方向，财富之树牌阵是一个很好的选择。它能帮助你分析财富增长的根基、所需能量、可能遇到的阻碍和潜在危险，最终预测你能达到的财富高度。", suited: " 适合事业发展 & 财运状况", positions: []string{"生长的根基", "依赖的能量", "遇到的阻碍", "潜在的危险", "最终的高度"}},
	{name: "二选一牌阵", zhname: "二选一牌阵", summary: "二选一牌阵专门用于在两种情况中做出选择。它在判断情势、决定方向等解析中应用广泛，可用于感情、事业、学业等多个领域的占卜。这个牌阵的用途相对广泛，特别适合在犹豫不决时使用。通过分析每个选择的发展和最终结果，二选一牌阵能帮助你更清晰地看清每个选项的利弊，从而做出更明智的决定。", suited: "主要用于抉择 & 判断", positions: []string{"问题的现况", "选择A的发展", "选择B的发展", "选择A的最终结果", "选择B的最终结果"}},
	{name: "维纳斯牌阵", zhname: "维纳斯牌阵", summary: "维纳斯牌阵以其直指未来的超然力量而闻名，能够快速而真实地预测爱情的内核和发展。这是一个罕见的经典塔罗牌阵，特别适合用于婚姻和恋爱方面的未来指向占卜。它能够深入分析爱情的未来，让我们清晰地看到双方在未来的状况。维纳斯牌阵不仅能揭示双方的真实想法，还能分析关系对双方的影响，预见可能遇到的障碍，最终给出关系的发展结果。如果你想占卜自己爱情的未来走向，维纳斯牌阵无疑是首选。", suited: "用于占卜爱情发展 & 预测走向", positions: []string{"自己的真实想法", "对方的真实想法", "彼此关系对自己的影响", "彼此关系对对方的影响", "你们将会遇到的障碍", "最后的结果", "将来自己的心情", "将来对方的心情"}},
	{name: "周运势牌阵", zhname: "周运势牌阵", summary: "周运势牌阵是一个专门用于分析一周七天运势的牌阵。它能够帮助我们逐天进行分析，提前做好预判，同时也可以通过整体牌面来预测本周运势的高低。这个牌阵是周运占卜的专用工具，不仅可以用来占卜下一周的运势，还可以应用在任何有七天期限的占卜中，而不影响效果。通过周运势牌阵，你可以对未来一周的整体趋势有清晰的把握，为每一天做好充分的准备。", suited: "适合周运分析 & 单周占卜", positions: []string{"周运势第一天", "周运势第二天", "周运势第三天", "周运势第四天", "周运势第五天", "周运势第六天", "周运势第七天"}},
	{name: "X机会牌阵", zhname: "X机会牌阵", summary: "X机会牌阵专门用于解决参与时机的问题。它以事物处理时机为主轴线，帮助我们拿捏问题解决的成功几率。这个牌阵能增强自身的洞察能力，让我们做出对未来的准确预判。它特别适合用于审查时机，捕捉稍纵即逝的机会，给出一个或弃或留的选择。如果你不想错过大好时机，X机会牌阵就是你审视机会的最佳工具。通过分析自己的心态、眼前的时机、成功的几率以及影响因素，这个牌阵能帮助你做出明智的决策。", suited: "适合时机捕捉 & 临场决策", positions: []string{"代表你自己的心态", "代表眼前的时机", "代表成功的几率", "代表影响的因素", "代表未来发展、最后结果"}},
}
