package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	bankroll         = 10000
	numrolls         = 1000
	numgames         = 10000
	lay5payment      = 0.67
	place6or8value   = 25
	place6or8payment = 1.17
	fieldbetvalue    = 15
	fieldpaysdouble  = true
	fieldhastriple   = true
)

var maxState = 0
var maxStateRolls = 0

func main() {
	highScore := 0
	totalPlayers := 0
	playersSurvived := 0
	winners := 0
	var averageOutcome float64 = 0
	var averageRolls float64 = 0
	var averagePeak float64 = 0

	for i := numgames; i > 0; i-- {
		game := createGameState()
		total, rolls, localMax := game.performIronCrossLay5()

		totalPlayers = totalPlayers + 1

		if total > 200 {
			playersSurvived = playersSurvived + 1
		}

		if total > bankroll {
			winners = winners + 1
		}

		averageOutcome = averageOutcome + (float64(total)-averageOutcome)/float64(totalPlayers)
		averageRolls = averageRolls + (float64(rolls)-averageRolls)/float64(totalPlayers)
		averagePeak = averagePeak + (float64(localMax)-averagePeak)/float64(totalPlayers)

		if total > highScore {
			highScore = total
		}
	}
	fmt.Printf("After %v players each taking %v rolls, we have: \n - highscore: %v\n - playersSurvived: %v\n - winners: %v\n - averageOutcome: %v\n - averageRolls %v\n - maxState %v in %v rolls\n - averagePeakBankroll %v", numgames, numrolls, highScore, playersSurvived, winners, averageOutcome, averageRolls, maxState, maxStateRolls, averagePeak)
}

type gameState struct {
	pointOn        bool
	point          int
	lastRoll       int
	rollsThrown    int
	playerBankroll int
	lay5value      int
	place6value    int
	place8value    int
	fieldBetValue  int
}

func createGameState() *gameState {
	newGame := &gameState{
		pointOn:        false,
		point:          0,
		lastRoll:       0,
		rollsThrown:    0,
		playerBankroll: bankroll,
		lay5value:      0,
	}

	return newGame
}

func (g *gameState) performIronCrossLay5() (int, int, int) {
	localMax := g.playerBankroll
	for g.rollsThrown < numrolls {
		g.bet()
		if g.noMoneyOnTable() {
			break
		}
		newRoll := g.rollTheDice()
		g.processRoll(int(newRoll))
		if g.playerBankroll > maxState {
			maxState = g.playerBankroll
			maxStateRolls = g.rollsThrown
		}

		if g.playerBankroll > localMax {
			localMax = g.playerBankroll
		}
	}
	//fmt.Printf("Bankroll: %v, rollsThrown: %v \n", g.playerBankroll, g.rollsThrown)
	return g.playerBankroll, g.rollsThrown, localMax
}

func (g *gameState) noMoneyOnTable() bool {
	return g.pointOn && g.lay5value == 0 && g.place6value == 0 && g.place8value == 0 && g.fieldBetValue == 0
}

func (g *gameState) rollTheDice() int64 {
	n1Big, err := rand.Int(rand.Reader, big.NewInt(5))
	if err != nil {
		panic(err)
	}
	n2Big, err := rand.Int(rand.Reader, big.NewInt(5))
	if err != nil {
		panic(err)
	}
	n := n1Big.Int64() + 1
	n2 := n2Big.Int64() + 1

	g.rollsThrown = g.rollsThrown + 1
	return n + n2
}

func (g *gameState) processRoll(roll int) {
	if g.pointOn {
		switch roll {
		case 2:
			g.payTheField(false)
		case 3:
			g.payTheField(false)
		case 4:
			g.payTheField(false)
		case 5:
			g.clear5()
		case 6:
			g.payThe6()
		case 7:
			g.sevenOut()
		case 8:
			g.payThe8()
		case 9:
			g.payTheField(false)
		case 10:
			g.payTheField(false)
		case 11:
			g.payTheField(false)
		case 12:
			g.payTheField(true)
		}
	} else {
		if roll == 4 || roll == 5 || roll == 6 || roll == 8 || roll == 9 || roll == 10 {
			g.pointOn = true
			g.point = roll
		}
	}
}

func (g *gameState) payTheField(triple bool) {
	if triple {
		g.playerBankroll = g.playerBankroll + 3*g.fieldBetValue
	} else {
		g.playerBankroll = g.playerBankroll + 2*g.fieldBetValue
	}
	g.fieldBetValue = 0
}

func (g *gameState) payThe6() {
	g.playerBankroll = g.playerBankroll + int(place6or8payment*float64(g.place6value)) + g.place6value
	g.place6value = 0
}

func (g *gameState) clear5() {
	g.lay5value = 0
}

func (g *gameState) payThe8() {
	g.playerBankroll = g.playerBankroll + int(place6or8payment*float64(g.place8value)) + g.place8value
	g.place8value = 0
}

func (g *gameState) sevenOut() {
	g.playerBankroll = g.playerBankroll + int(lay5payment*float64(g.lay5value)) + g.lay5value
	g.place8value = 0
	g.place6value = 0
	g.fieldBetValue = 0
	g.lay5value = 0
	g.pointOn = false
	g.point = 0
}

func (g *gameState) bet() {
	if g.pointOn {
		if g.fieldBetValue == 0 && g.playerBankroll >= 15 {
			g.fieldBetValue = 15
			g.playerBankroll = g.playerBankroll - 15
		}

		if g.lay5value == 0 && g.playerBankroll >= 150 {
			g.lay5value = 150
			g.playerBankroll = g.playerBankroll - 150
		}

		if g.place6value == 0 && g.playerBankroll >= 25 {
			g.place6value = 25
			g.playerBankroll = g.playerBankroll - 25
		}

		if g.place8value == 0 && g.playerBankroll >= 25 {
			g.place8value = 25
			g.playerBankroll = g.playerBankroll - 25
		}
	}
}
