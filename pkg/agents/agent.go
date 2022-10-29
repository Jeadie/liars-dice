package agents

import "github.com/Jeadie/liars-poker/pkg/game"

type Agent interface {
	Play(r game.Round) game.Action
	Handle(e game.Event)
}
