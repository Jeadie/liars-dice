package game

import (
	"fmt"
	"math/rand"
)

type Bet [2]uint
type Agent int
type ActionType string

const (
	Raise ActionType = "raise"
	Exact ActionType = "exact"
	Call  ActionType = "call"
)

func ParseActionType(x string) (ActionType, error) {
	var actionType ActionType
	if x != string(Call) && x != string(Raise) && x != string(Exact) {
		return actionType, fmt.Errorf("%s is not a valid action", x)
	}
	return ActionType(x), nil
}

type Action struct {
	T     ActionType `json:"T"`
	Raise Bet        `json:"Raise,omitempty"`
}

func (a Action) ToString() string {
	switch a.T {
	case Raise:
		return fmt.Sprintf("raises to %d %ds", a.Raise[0], a.Raise[1])
	case Call:
		return fmt.Sprintf("calls out current bet")
	case Exact:
		return fmt.Sprintf("calls out exactly the current bet")
	default:
		return ""
	}
}

type Round struct {
	Dice      [][]uint `json:"Dice"` // TODO: probably want a *[][]uint
	UsedOne   bool     `json:"UsedOne"`
	CurrBet   Bet      `json:"Bet"`
	CurrAgent Agent    `json:"Agent"`
}

// Raises returns true if the Bet, b can raise the Bet c. A Bet can raise another
// if it calls a higher quantity, zeroth index, or an equal quantity of a higher
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
		return fmt.Errorf("agent %d cannot Raise. It is Agent %d's turn", agent, r.CurrAgent)
	}
	if !bet.Raises(r.CurrBet) {
		return fmt.Errorf("%d %ds does not raise %d %ds", bet[0], bet[1], r.CurrBet[0], r.CurrBet[1])
	}
	r.CurrBet = bet
	r.UsedOne = r.UsedOne || bet[1] == 1
	return nil
}

// Exact called by an Agent returns true iff the current CurrBet is exactly equal to the state of the Dice.
// Error if Agent cannot currently make an Exact call.
func (r *Round) Exact(agent Agent) (bool, error) {
	if r.CurrAgent != agent {
		return false, fmt.Errorf("agent %d cannot Exact. It is Agent %d's turn", agent, r.CurrAgent)
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
// of the Dice. Error if Agent cannot currently make a Calls.
func (r *Round) Calls(agent Agent) (bool, error) {
	if r.CurrAgent != agent {
		return false, fmt.Errorf("agent %d cannot Call. It is Agent %d's turn", agent, r.CurrAgent)
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

// PlayTurn on a board. Returns the number of dice an Agent should gain/lose after that turn. An
// error is returned if the action cannot be played at this stage of the round.
func (r *Round) PlayTurn(agent Agent, action Action) (Agent, int, error) {
	var isExactly bool
	var isAtLeast bool
	var err error

	switch action.T {
	case Raise:
		err = r.Raise(agent, action.Raise)
	case Call:
		isAtLeast, err = r.Calls(agent)
	case Exact:
		isExactly, err = r.Exact(agent)
	default:
		return agent, 0, fmt.Errorf("action, %s, is not valid", action.T)
	}
	if err != nil {
		return 0, 0, err
	}

	if action.T == Raise {
		return 0, 0, nil
	}

	// Round finished. determine outcome.
	if action.T == Call {
		if isAtLeast {
			return r.CurrAgent, -1, nil
		} else {
			return Agent(Mod(int(r.CurrAgent)-1, len(r.Dice))), -1, nil
		}
	} else if action.T == Exact {
		if isExactly {
			return r.CurrAgent, 1, nil
		}
		return r.CurrAgent, -2, nil
	}
	return 0, 0, fmt.Errorf("an unrecognisable error occurred")
}

func Mod(x, m int) int {
	y := x % m
	if y < 0 {
		return y + m
	}
	return y
}

func (r *Round) IncrementCurrentAgent() {
	r.CurrAgent = Agent((int(r.CurrAgent) + 1) % len(r.Dice))
}
