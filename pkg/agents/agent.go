package agents

import "github.com/Jeadie/liars-poker/pkg/game"

type Agent interface {
	Initialise(r game.Round, dice []uint)
	Play(r game.Round) game.Action
	Handle(e game.Event)
}
