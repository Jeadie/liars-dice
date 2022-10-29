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
	Call             = "bet"
)

type Action struct {
	t     ActionType
	agent Agent
	bet   Bet
}

type Round struct {
	dice    [][]uint
	usedOne bool
	bet     Bet
	i       Agent
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

// Raise a bet on the round by a specific agent. Returns an error if agent cannot currently make a
// raise (i.e. it's not their turn), or the raise is not valid given the state of the round.
func (r *Round) Raise(agent Agent, bet Bet) error {
	if r.i != agent {
		return fmt.Errorf("agent %d cannot Raise. It is agent %d's turn", agent, r.i)
	}
	if !bet.Raises(r.bet) {
		return fmt.Errorf("bet %d %ds cannot be played above %d %ds", bet[0], bet[1], r.bet[0], r.bet[1])
	}
	r.bet = bet
	r.i = Agent((int(r.i) + 1) % len(r.dice))
	r.usedOne = r.usedOne || bet[1] == 1
	return nil
}

// Exact called by an agent returns true iff the current bet is exactly equal to the state of the dice.
// Error if agent cannot currently make an Exact call.
func (r *Round) Exact(agent Agent) (bool, error) {
	if r.i != agent {
		return false, fmt.Errorf("agent %d cannot Exact. It is agent %d's turn", agent, r.i)
	}
	q, v := r.bet[0], r.bet[1]
	tot := uint(0)
	for _, dice := range r.dice {
		for _, die := range dice {
			if die == v || (die == 1 && !r.usedOne) {
				tot++
			}
		}
	}
	return q == tot, nil
}

// Calls called by an agent returns true iff the current bet is greater than or equal to the state
// of the dice. Error if agent cannot currently make an Calls.
func (r *Round) Calls(agent Agent) (bool, error) {
	if r.i != agent {
		return false, fmt.Errorf("agent %d cannot Call. It is agent %d's turn", agent, r.i)
	}
	q, v := r.bet[0], r.bet[1]
	tot := uint(0)
	for _, dice := range r.dice {
		for _, die := range dice {
			if die == v || (die == 1 && !r.usedOne) {
				tot++
			}
		}
	}
	return tot >= q, nil
}

func (r *Round) GetDicePerAgent() []uint {
	numDice := make([]uint, len(r.dice))
	for i, die := range r.dice {
		numDice[i] = uint(len(die))
	}
	return numDice
}

// InitRound based on the number of dice for each agent. Expects i < len(numDice)
func InitRound(numDice []uint, i Agent) *Round {
	return &Round{
		dice:    InitDice(numDice),
		usedOne: false,
		bet:     [2]uint{0, 0},
		i:       i,
	}
}

// InitDice based on the number of dice
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

// PlayRound of liar's dice. Returns the new number of dice for the next round.
func PlayRound(r Round, actions chan Action, errs chan error) []uint {
	var err error
	var correct bool

	dice := r.GetDicePerAgent()
	for action := range actions {
		switch action.t {
		case Raise:
			err = r.Raise(action.agent, action.bet)
		case Call:
			correct, err = r.Calls(action.agent)
		case Exact:
			correct, err = r.Exact(action.agent)
		}
		if err != nil {
			errs <- err
			continue
		}
		if action.t == Raise {
			continue
		}

		// Round finished. determine outcome.
		if action.t == Call {
			losingAgent := r.i
			if correct {
				losingAgent = Agent((int(r.i) - 1) % len(r.dice))
			}
			dice[losingAgent]--
			break
		} else if action.t == Exact {
			if correct {
				dice[r.i] += 2
			}
			break
		}
	}
	return dice
}
