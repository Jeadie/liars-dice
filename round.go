package liars_dice

import (
	"fmt"
	"math/rand"
)

type Bet [2]uint
type Agent int
type ActionType string

const (
	Raise ActionType = "raise"
	Exact            = "exact"
	Call             = "CurrBet"
)

type Action struct {
	T     ActionType `json:"T"`
	Agent Agent      `json:"Agent"`
	Raise Bet        `json:"Raise,omitempty"`
}

type Round struct {
	Dice      [][]uint `json:"Dice"`
	UsedOne   bool     `json:"UsedOne"`
	CurrBet   Bet      `json:"Bet"`
	CurrAgent Agent    `json:"Agent"`
}

// Raises returns true if the Bet, b can raise the Bet c. A Bet can raise another
// if it calls a higher quanity, zeroth index, or an equal quantity of a higher
// value, first index. E.g.
//
//	(4, 2) raise (3, 2)
//	(2, 4) raise (2, 3)
//	(2, 2) does not raise (2, 2)
func (b Bet) Raises(c Bet) bool {
	return 10*b[0]+b[1] > 10*c[0]+c[1]
}

// Raise a CurrBet on the round by a specific Agent. Returns an error if Agent cannot currently make a
// raise (CurrAgent.e. it's not their turn), or the raise is not valid given the state of the round.
func (r *Round) Raise(agent Agent, bet Bet) error {
	if r.CurrAgent != agent {
		return fmt.Errorf("Agent %d cannot Raise. It is Agent %d's turn", agent, r.CurrAgent)
	}
	if !bet.Raises(r.CurrBet) {
		return fmt.Errorf("CurrBet %d %ds cannot be played above %d %ds", bet[0], bet[1], r.CurrBet[0], r.CurrBet[1])
	}
	r.CurrBet = bet
	r.CurrAgent = Agent((int(r.CurrAgent) + 1) % len(r.Dice))
	r.UsedOne = r.UsedOne || bet[1] == 1
	return nil
}

// Exact called by an Agent returns true iff the current CurrBet is exactly equal to the state of the Dice.
// Error if Agent cannot currently make an Exact call.
func (r *Round) Exact(agent Agent) (bool, error) {
	if r.CurrAgent != agent {
		return false, fmt.Errorf("Agent %d cannot Exact. It is Agent %d's turn", agent, r.CurrAgent)
	}
	q, v := r.CurrBet[0], r.CurrBet[1]
	tot := uint(0)
	for _, dice := range r.Dice {
		for _, die := range dice {
			if die == v || (die == 1 && !r.UsedOne) {
				tot++
			}
		}
	}
	return q == tot, nil
}

// Calls called by an Agent returns true iff the current CurrBet is greater than or equal to the state
// of the Dice. Error if Agent cannot currently make an Calls.
func (r *Round) Calls(agent Agent) (bool, error) {
	if r.CurrAgent != agent {
		return false, fmt.Errorf("Agent %d cannot Call. It is Agent %d's turn", agent, r.CurrAgent)
	}
	q, v := r.CurrBet[0], r.CurrBet[1]
	tot := uint(0)
	for _, dice := range r.Dice {
		for _, die := range dice {
			if die == v || (die == 1 && !r.UsedOne) {
				tot++
			}
		}
	}
	return tot >= q, nil
}

func (r *Round) GetDicePerAgent() []uint {
	numDice := make([]uint, len(r.Dice))
	for i, die := range r.Dice {
		numDice[i] = uint(len(die))
	}
	return numDice
}

// InitRound based on the number of Dice for each Agent. Expects CurrAgent < len(numDice)
func InitRound(numDice []uint, i Agent) *Round {
	return &Round{
		Dice:      InitDice(numDice),
		UsedOne:   false,
		CurrBet:   [2]uint{0, 0},
		CurrAgent: i,
	}
}

// InitDice based on the number of Dice
func InitDice(numDice []uint) [][]uint {
	dice := make([][]uint, len(numDice))
	for i, die := range numDice {
		dice[i] = make([]uint, die)
		for j := uint(0); j < die; j++ {
			dice[i][j] = uint(rand.Intn(6) + 1)
		}
	}
	return dice
}

// PlayRound of liar's Dice. Returns the new number of Dice for the next round.
func PlayRound(r Round, actions chan Action, errs chan error) []uint {
	var err error
	var correct bool

	dice := r.GetDicePerAgent()
	for action := range actions {
		switch action.T {
		case Raise:
			err = r.Raise(action.Agent, action.Raise)
		case Call:
			correct, err = r.Calls(action.Agent)
		case Exact:
			correct, err = r.Exact(action.Agent)
		}
		if err != nil {
			errs <- err
			continue
		}
		if action.T == Raise {
			continue
		}

		// Round finished. determine outcome.
		if action.T == Call {
			losingAgent := r.CurrAgent
			if correct {
				losingAgent = Agent((int(r.CurrAgent) - 1) % len(r.Dice))
			}
			dice[losingAgent]--
			break
		} else if action.T == Exact {
			if correct {
				dice[r.CurrAgent] += 2
			}
			break
		}
	}
	return dice
}
