package agents

import "github.com/Jeadie/liars-poker/pkg/game"

type ProbAgent struct {
}

func ConstructProbAgent() *ProbAgent {
	return &ProbAgent{}
}

func (h *ProbAgent) Play(r game.Round) game.Action {
	if r.CurrBet[1] == 6 {
		return game.Action{
			T:     game.Raise,
			Raise: game.Bet{r.CurrBet[0] + 1, 1},
		}
	} else {
		return game.Action{
			T:     game.Raise,
			Raise: game.Bet{r.CurrBet[0], r.CurrBet[1] + 1},
		}
	}
}

func (h *ProbAgent) Initialise(r game.Round, dice []uint) {

}
