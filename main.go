package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Ranksai/patigo/patigo"
)

const Debug = false

func main() {

	pay := &patigo.Payouts{
		{600, 0.42},
		{900, 0.0},
		{1200, 0.07},
		{2400, 0.51},
	}
	hokutoMusou := patigo.NewSt(
		0, "北斗無双", "幻闘ラッシュ",
		17, 0.5, 319.7,
		81.2, 130, 100, pay,
	)

	fmt.Printf("%+v\n", hokutoMusou)

	source := rand.NewSource(time.Now().UnixNano())
	if Debug {
		count := 0
		countRush := 0
		num := 1000000
		for i := 0; i < num; i++ {
			if hokutoMusou.CalcBonus(source) {
				count++
			}
			if hokutoMusou.CalcRushIn(source) {
				countRush++
			}
		}
		fmt.Printf("当たり回数: %d\n", count)
		fmt.Printf("確率: %g\n", float64(count)/float64(num))
		fmt.Printf("理論確率: %g\n", hokutoMusou.RealRowProbability)

		fmt.Printf("当たり回数: %d\n", countRush)
		fmt.Printf("確率: %g\n", float64(countRush)/float64(num))
		fmt.Printf("理論確率: %g\n", hokutoMusou.RushIn)
	}

	for i := 0; i < 10; i++ {
		hokutoMusou.PlayGame(source)
	}
}
