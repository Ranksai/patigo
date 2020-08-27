package patigo

import (
	"fmt"
	"math/rand"
)

type ST struct {
	Id                  int
	Name                string
	StName              string
	Border              int
	RushIn              float64
	RowProbability      float64
	HighProbability     float64
	RushRotations       int
	NoRushRotations     int
	RealRowProbability  float64
	RealHighProbability float64
	BonusPayouts        *Payouts
	ResultGame          *resultGame
}

type resultGame struct {
	BonusCount int
	PayOut     int
}

func NewSt(id int, name, stName string, border int, rushIn, rowProbability, highProbability float64, rushRotations, noRushRotations int, bonusPayouts *Payouts) *ST {
	game := &resultGame{
		BonusCount: 0,
		PayOut:     0,
	}
	return &ST{
		Id:                  id,
		Name:                name,
		StName:              stName,
		Border:              border,
		RushIn:              rushIn,
		RowProbability:      rowProbability,
		HighProbability:     highProbability,
		RushRotations:       rushRotations,
		NoRushRotations:     noRushRotations,
		RealRowProbability:  1.0 / rowProbability,
		RealHighProbability: 1.0 / highProbability,
		BonusPayouts:        bonusPayouts,
		ResultGame:          game,
	}
}

// normal bonus
func (st *ST) CalcBonus(source rand.Source) bool {
	return CalcHitProbability(st.RealRowProbability, source)
}

func (st *ST) CalcRushIn(source rand.Source) bool {
	return CalcHitProbability(st.RushIn, source)
}

// st bonus
func (st *ST) CalcRushInBonus(source rand.Source) bool {
	return CalcHitProbability(st.RealHighProbability, source)
}

func (st *ST) CalcRushInBonusPayout(source rand.Source) int {
	round := CalcHitProbabilityRound(st.BonusPayouts, source)
	fmt.Printf("ラウンド: %d\n", round)
	fmt.Printf("出玉: %d\n", (*st.BonusPayouts)[round].PayoutPoint)
	return round
}

func (st *ST) CalcReturnRush(source rand.Source) bool {
	fmt.Println("--- return rush start ---")
	for i := 0; i < st.NoRushRotations; i++ {
		if st.CalcBonus(source) {
			fmt.Printf("引き戻し、回転数: %d\n", i)
			return true
		}
	}
	fmt.Printf("時短終了、回転数: %d\n", st.NoRushRotations)
	fmt.Println("--- return rush end ---")
	return false
}

// play game
func (st *ST) PlayGame(source rand.Source) {
	normalCount := 0
	totalCount := 0
	for {
		normalCount++
		if st.CalcBonus(source) {
			fmt.Printf("大当たり! 回転数: %d\n", normalCount)
			totalCount += normalCount
			st.ResultGame.PayOut -= ((normalCount / st.Border) + 1) * 250
			normalCount = 0
			if st.CalcRushIn(source) {
				fmt.Printf("%s 突入\n", st.StName)
				result := st.PlayRush(source)
				st.ResultGame.PayOut += result.PayOut
				st.ResultGame.BonusCount += result.BonusCount
				break
			} else {
				fmt.Printf("%s\n", "激闘Bonus")
				if st.CalcReturnRush(source) {
					result := st.PlayRush(source)
					st.ResultGame.PayOut += result.PayOut
					st.ResultGame.BonusCount += result.BonusCount
					break
				}
				break
			}
		}
	}
	fmt.Printf("最終合計出玉: %d\n", st.ResultGame.PayOut)
}

// play st
func (st *ST) PlayRush(source rand.Source) *resultGame {
	fmt.Println("--- rush start ---")
	result := &resultGame{
		BonusCount: 0,
		PayOut:     0,
	}
	for i := 0; i < st.RushRotations; i++ {
		if st.CalcRushInBonus(source) {
			fmt.Printf("ST中当たり 回転数: %d\n", i)
			payoutId := st.CalcRushInBonusPayout(source)
			i = 0
			result.BonusCount++
			result.PayOut += (*st.BonusPayouts)[payoutId].PayoutPoint
			fmt.Printf("合計出玉: %d\n", result.PayOut)
			continue
		}
	}
	fmt.Println("ST終了")
	fmt.Printf("ST当たり回数: %d\n", result.BonusCount)
	fmt.Printf("合計出玉: %d\n", result.PayOut)
	fmt.Println("-- rush end---")
	return result
}
