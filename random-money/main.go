package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const initMoney = int64(2000)
const maxBetMoney = int64(200)

func maxBet(money int64) int64 {
	if money < maxBetMoney {
		return money
	}
	return maxBetMoney
}

func fixedRatioBet(money int64) int64 {
	n := int64(float64(money) * 0.1)
	return n
}

func main() {
	rand.Seed(time.Now().Unix())
	humanNums := 100
	roll := 10000
	gh := &gamblingHouse{peoples: make([]gambler, humanNums)}

	humans := make([]*people, humanNums)

	worker := func() int64 {
		return 2
	}

	gambleFund := fixedRatioBet

	for i := 0; i < humanNums; i++ {
		p := &people{money: initMoney, id: i, worker: worker, gambleFund: gambleFund}
		humans[i] = p
	}
	for i := 0; i < roll; i++ {
		for _, p := range humans {
			p.month(gh)
		}
		gh.rock()
		gh.apportion()
		gh.reset()
	}

	for _, h := range humans {
		h.keepAccounts()
	}

	sort.Slice(humans, func(i, j int) bool {
		return humans[i].money <= humans[j].money
	})

	printNum := 5
	num := make([]int, printNum*2)
	for i := 0; i < printNum; i++ {
		num[i] = i
		num[printNum*2-i-1] = len(humans) - i - 1
	}
	roll = len(humans[0].account)
	for i := 0; i < roll; i++ {
		for _, v := range num {
			fmt.Printf("%d\t", humans[v].account[i].money)
		}
		fmt.Println()
	}
}
