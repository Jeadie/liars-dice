package agents

import (
	"fmt"
	"github.com/Jeadie/liars-dice/pkg/game"
)

type ProbAgent struct {
	selfIdx      int
	agentNumDice []uint
	roundDice    []uint
}

func ConstructProbAgent() *ProbAgent {
	return &ProbAgent{}
}

func Sum(x []uint) uint {
	tot := uint(0)
	for _, xx := range x {
		tot += xx
	}
	return tot
}

func Count(x []uint, v uint) uint {
	tot := uint(0)
	for _, xx := range x {
		if xx == v {
			tot += 1
		}
	}
	return tot
}

func (h *ProbAgent) CalculateOptimalBet(current game.Bet) game.Bet {
	freq := make(map[uint]uint, 6)
	for _, die := range h.roundDice {
		freq[die]++
	}

	maxV := uint(0)
	maxQ := uint(0)
	for k, v := range freq {
		if v > maxQ {
			maxQ = v
			maxV = k
		}
	}

	q, v := current[0], current[1]
	if maxV > v {
		return game.Bet{q, maxV}
	} else {
		return game.Bet{q + 1, maxV}
	}
}

func (h *ProbAgent) Play(r game.Round) game.Action {
	// TODO: consider the affects of 1's

	qCurr, vCurr := r.CurrBet[0], r.CurrBet[1]

	numDice := Sum(h.agentNumDice)

	eqFromAgent := Count(h.roundDice, vCurr)
	eqFromOthers := float32(numDice-uint(len(h.roundDice))) / 6.0
	neededQFromOthers := float32(qCurr) - float32(eqFromAgent)

	// fmt.Printf("[DEBUG]- Hand=%d, |Dice|=%d, v=%d, E[v]=%f, QFromOthers(needed)=%f\n", h.roundDice, numDice, vCurr, eq, neededQFromOthers)

	if neededQFromOthers == eqFromOthers {
		return game.Action{
			T: game.Exact,
		}
	} else if neededQFromOthers > eqFromOthers {
		return game.Action{
			T: game.Call,
		}
	} else {
		return game.Action{
			T:     game.Raise,
			Raise: h.CalculateOptimalBet(r.CurrBet),
		}
	}
}

func (h *ProbAgent) Handle(e game.Event) {
	switch e.EType {
	case game.Turn:
		// action, agent := e.Turn.Action, e.Turn.ActionAgent
		// Simple ProbAgent does not consider how other agents play. Just
		// information of their own hand.
		break

	case game.RoundComplete:
		change, agentIdx := e.RoundComplete.ChangeInDice, e.RoundComplete.AffectedAgent
		if int(h.agentNumDice[agentIdx])+change <= 0 {
			h.agentNumDice[agentIdx] = 0
		} else {
			h.agentNumDice[agentIdx] = uint(int(h.agentNumDice[agentIdx]) + change)
		}

	case game.AgentRemoved:
		// Don't need this. RoundComplete is sufficient
		break

	case game.InvalidAction:
		// Do not expect to receive invalid Actions.
		break

	case game.RoundStart:
		h.roundDice = e.RoundStart.DiceRolled

	case game.GameStart:
		h.agentNumDice = e.GameStart.NumDicePerAgent
		h.selfIdx = e.GameStart.AgentIdx

	default:
		fmt.Printf("Unknown event %s\n", e.EType)
	}
}
