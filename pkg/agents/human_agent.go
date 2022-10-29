package agents

import "github.com/Jeadie/liars-poker/pkg/game"

type HumanAgent struct {
	dice []uint
}

func ConstructHuman() *HumanAgent {
	return &HumanAgent{}
}

func (h *HumanAgent) Play(r game.Round) game.Action {
	return game.Action{
		T: game.Exact,
	}
}

func (h *HumanAgent) Initialise(r game.Round, dice []uint) {
	h.dice = dice
}
