package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	// not accounting for commission yet
	bankroll         = 1000
	numrolls         = 1000
	numgames         = 20000
	lay5payment      = 0.67
	place6or8payment = 1.17
	lay5Bet          = 100
	placeBet         = 25
	fieldBet         = 0
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
	shouldPress    bool
}

func createGameState() *gameState {
	newGame := &gameState{
		pointOn:        false,
		point:          0,
		lastRoll:       0,
		rollsThrown:    0,
		playerBankroll: bankroll,
		lay5value:      0,
		shouldPress:    true,
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
	//fmt.Printf("Bankroll: %v, rollsThrown: %v, localMax: %v\n", g.playerBankroll, g.rollsThrown, localMax)
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
	profit := 0
	if g.pointOn {
		switch roll {
		case 2:
			profit = g.payTheField(false)
		case 3:
			profit = g.payTheField(false)
		case 4:
			profit = g.payTheField(false)
		case 5:
			profit = g.clear5()
		case 6:
			profit = g.payThe6()
		case 7:
			profit = g.sevenOut()
		case 8:
			profit = g.payThe8()
		case 9:
			profit = g.payTheField(false)
		case 10:
			profit = g.payTheField(false)
		case 11:
			profit = g.payTheField(false)
		case 12:
			profit = g.payTheField(true)
		}
	} else {
		if roll == 4 || roll == 5 || roll == 6 || roll == 8 || roll == 9 || roll == 10 {
			g.pointOn = true
			g.point = roll
		}
	}
	if profit > 0 {
		//fmt.Printf("roll profit: %v\n", profit)
	}
}

func (g *gameState) payTheField(triple bool) int {
	var payment int
	if triple {
		payment = 3 * g.fieldBetValue
	} else {
		payment = 2 * g.fieldBetValue
	}

	g.playerBankroll = g.playerBankroll + payment
	g.fieldBetValue = 0
	return payment
}

func (g *gameState) payThe6() int {
	payment := 0
	if g.shouldPress && g.playerBankroll > g.place6value {
		g.place6value = 2 * g.place6value
		g.playerBankroll = g.playerBankroll - g.place6value
		g.shouldPress = false
	} else {
		payment = int(place6or8payment*float64(g.place6value)) + g.place6value - g.fieldBetValue
		g.playerBankroll = g.playerBankroll + payment
		g.place6value = 0
	}
	g.fieldBetValue = 0

	return payment
}

func (g *gameState) clear5() int {
	loss := g.lay5value - g.fieldBetValue
	g.lay5value = 0
	g.fieldBetValue = 0

	return loss
}

func (g *gameState) payThe8() int {
	payment := 0
	if g.shouldPress && g.playerBankroll > g.place6value {
		g.place8value = 2 * g.place8value
		g.playerBankroll = g.playerBankroll - g.place8value
		g.shouldPress = false
	} else {
		payment := int(place6or8payment*float64(g.place8value)) + g.place8value - g.fieldBetValue
		g.playerBankroll = g.playerBankroll + payment
		g.place8value = 0
	}
	g.fieldBetValue = 0

	return payment
}

func (g *gameState) sevenOut() int {
	profit := int(lay5payment*float64(g.lay5value)) + g.lay5value
	g.playerBankroll = g.playerBankroll + profit
	profit = profit - g.place8value
	g.place8value = 0
	profit = profit - g.place6value
	g.place6value = 0
	profit = profit - g.fieldBetValue
	g.fieldBetValue = 0
	g.lay5value = 0
	g.pointOn = false
	g.point = 0
	g.shouldPress = true

	return profit
}

func (g *gameState) bet() {
	if g.pointOn {
		if g.fieldBetValue == 0 && g.playerBankroll >= fieldBet {
			g.fieldBetValue = fieldBet
			g.playerBankroll = g.playerBankroll - fieldBet
		}

		if g.lay5value == 0 && g.playerBankroll >= lay5Bet {
			g.lay5value = lay5Bet
			g.playerBankroll = g.playerBankroll - lay5Bet
		}

		if g.place6value == 0 && g.playerBankroll >= placeBet {
			g.place6value = placeBet
			g.playerBankroll = g.playerBankroll - placeBet
		}

		if g.place8value == 0 && g.playerBankroll >= placeBet {
			g.place8value = placeBet
			g.playerBankroll = g.playerBankroll - placeBet
		}
	}
}
