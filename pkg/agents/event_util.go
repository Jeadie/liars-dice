package agents

import (
	"github.com/Jeadie/liars-dice/pkg/game"
	"github.com/rs/zerolog/log"
)

func SendGameStarted(agents []Agent, dice []uint) {
	var e game.Event
	for i, agent := range agents {
		e = game.Event{
			EType: game.GameStart,
			GameStart: &game.GameStartEvent{
				NumDicePerAgent: dice,
				AgentIdx:        i,
			},
		}
		agent.Handle(e)
	}
	//e.GameStart.AgentIdx = -1
	log.Debug().Interface("event", e).Send()
}

func SendGameComplete(agents []Agent, winner game.Agent) {
	var e game.Event
	for _, agent := range agents {
		e = game.Event{
			EType: game.GameComplete,
			GameComplete: &game.GameCompleteEvent{
				Winner: winner,
			},
		}
		agent.Handle(e)
	}
	log.Debug().Interface("event", e).Send()
}

func SendTurnEvent(agents []Agent, act game.Action, ag game.Agent) {
	e := game.Event{
		EType: game.Turn,
		Turn: &game.TurnEvent{
			Action:      act,
			ActionAgent: ag,
		},
	}
	log.Debug().Interface("event", e).Send()
	for _, agx := range agents {
		agx.Handle(e)
	}
}

func SendRoundStarted(agents []Agent, r game.Round) {
	var e game.Event
	for i, agent := range agents {
		e := game.Event{
			EType: game.RoundStart,
			RoundStart: &game.RoundStartEvent{
				DiceRolled: r.Dice[i],
			},
		}
		agent.Handle(e)
	}
	//e.RoundStart.DiceRolled = nil
	log.Debug().Interface("event", e).Send()
}

func SendRoundComplete(agents []Agent, ag game.Agent, diceChange int) {
	e := game.Event{
		EType: game.RoundComplete,
		RoundComplete: &game.RoundCompleteEvent{
			AffectedAgent: ag,
			ChangeInDice:  diceChange,
		},
	}
	log.Debug().Interface("event", e).Send()
	for _, agx := range agents {
		agx.Handle(e)
	}
}
