package patigo

import (
	"math/rand"
)

type Payout struct {
	PayoutPoint int
	Probability float64
}

type Payouts = []Payout

// probability: 0 - 1
func CalcHitProbability(probability float64, source rand.Source) bool {
	randFloat := rand.New(source).Float64()
	return randFloat < probability
}

func CalcHitProbabilityRound(probability *Payouts, source rand.Source) int {
	randFloat := rand.New(source).Float64()
	for i := 0; i < len(*probability); i++ {
		if randFloat < (*probability)[i].Probability && (*probability)[i].Probability != 0 {
			return i
		} else {
			randFloat -= (*probability)[i].Probability
		}
	}
	return len(*probability) + 1
}
