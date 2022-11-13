package game

import (
	agents2 "github.com/Jeadie/liars-dice/pkg/agents"
	"github.com/rs/zerolog/log"
)

func SendGameStarted(agents []agents2.Agent, dice []uint) {
	var e Event
	for i, agent := range agents {
		e = Event{
			EType: GameStart,
			GameStart: &GameStartEvent{
				NumDicePerAgent: dice,
				AgentIdx:        i,
			},
		}
		agent.Handle(e)
	}
	e.GameStart.AgentIdx = -1
	log.Debug().Interface("event", e).Send()
}

func SendGameComplete(agents []agents2.Agent, winner Agent) {
	var e Event
	for _, agent := range agents {
		e = Event{
			EType: GameComplete,
			GameComplete: &GameCompleteEvent{
				Winner: winner,
			},
		}
		agent.Handle(e)
	}
	log.Debug().Interface("event", e).Send()
}

func SendTurnEvent(agents []agents2.Agent, act Action, ag Agent) {
	e := Event{
		EType: Turn,
		Turn: &TurnEvent{
			Action:      act,
			ActionAgent: ag,
		},
	}
	log.Debug().Interface("event", e).Send()
	for _, agx := range agents {
		agx.Handle(e)
	}
}

func SendRoundStarted(agents []agents2.Agent, r Round) {
	var e Event
	for i, agent := range agents {
		e := Event{
			EType: RoundStart,
			RoundStart: &RoundStartEvent{
				DiceRolled: r.Dice[i],
			},
		}
		agent.Handle(e)
	}
	e.RoundStart.DiceRolled = nil
	log.Debug().Interface("event", e).Send()
}

func SendRoundComplete(agents []agents2.Agent, ag Agent, diceChange int) {
	e := Event{
		EType: RoundComplete,
		RoundComplete: &RoundCompleteEvent{
			AffectedAgent: ag,
			ChangeInDice:  diceChange,
		},
	}
	log.Debug().Interface("event", e).Send()
	for _, agx := range agents {
		agx.Handle(e)
	}
}
