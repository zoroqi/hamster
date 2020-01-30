package main

import (
	"math/rand"
)

type people struct {
	id         int
	money      int64                       // 手里的钱
	account    []account                   // 流水
	worker     func() int64                // 工作
	gambleFund func(selfMoney int64) int64 // 计算赌本
}

type account struct {
	money int64
	time  int
}

// 一个月的生活, 工作, 计算赌资, 赚钱
func (p *people) month(g *gamblingHouse) {
	p.keepAccounts()
	salary := p.worker()
	p.money += salary

	gambleFund := p.gambleFund(p.money)
	p.money -= gambleFund

	dice := rand.Intn(100)

	g.gamble(p, gambleFund, dice)
}

func (p *people) gambleResult(money int64) {
	p.money += money
}

func (p *people) keepAccounts() {
	p.account = append(p.account, account{p.money, len(p.account)})
}

type gamblingHouse struct {
	totalMoney    int64
	totalPoint    int64
	pointMoney    int64
	peoples       []gambler
	gamblerOffset int
}

type gambler struct {
	p    *people
	dice int
}

func (g *gamblingHouse) gamble(p *people, money int64, dice int) {
	g.peoples[g.gamblerOffset] = gambler{p: p, dice: dice}
	g.gamblerOffset++
	g.totalMoney += money
	g.totalPoint += int64(dice)
}

func (g *gamblingHouse) rock() {
	pointMoney := g.totalMoney / g.totalPoint
	g.pointMoney = pointMoney
}

func (g *gamblingHouse) apportion() {
	pointMoney := g.pointMoney
	loseMoney := int(g.totalMoney - g.pointMoney*g.totalPoint)
	for _, p := range g.peoples {
		p.p.gambleResult(pointMoney * int64(p.dice))
	}
	if loseMoney > 0 {
		l := len(g.peoples)
		rand.Shuffle(l, func(i, j int) {
			g.peoples[i], g.peoples[j] = g.peoples[j], g.peoples[i]
		})
		for i := 0; i < loseMoney; i++ {
			g.peoples[i%l].p.gambleResult(1)
		}
	}
}

func (g *gamblingHouse) reset() {
	g.totalPoint = 0
	g.totalMoney = 0
	g.pointMoney = 0
	g.gamblerOffset = 0
}
