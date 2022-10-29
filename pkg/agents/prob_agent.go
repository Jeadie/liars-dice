package agents

import "github.com/Jeadie/liars-poker/pkg/game"

type ProbAgent struct {
}

func ConstructProbAgent() *ProbAgent {
	return &ProbAgent{}
}

func (h *ProbAgent) Play(r game.Round) game.Action {
	q, v := r.CurrBet[0], r.CurrBet[1]

	if r.CurrBet[1] == 6 {
		return game.Action{
			T:     game.Raise,
			Raise: game.Bet{q + 1, 1},
		}
	} else {
		return game.Action{
			T:     game.Raise,
			Raise: game.Bet{q, v + 1},
		}
	}
}

func (h *ProbAgent) Handle(e game.Event) {

}
