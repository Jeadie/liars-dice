package play

import (
	"github.com/Jeadie/liars-dice/pkg/agents"
	"github.com/Jeadie/liars-dice/pkg/game"
	"github.com/rs/zerolog/log"
)

func PlayGame(ags []agents.Agent, dice []uint) {
	agents.SendGameStarted(ags, dice)
	winnerIdx, hasWon := WinningPlayer(dice)
	for !hasWon {
		round := game.InitRound(dice, 0)
		agentIdx, change := PlayRound(round, ags)

		// Changes score from last
		if int(dice[agentIdx])+change <= 0 {
			dice[agentIdx] = 0
		} else {
			dice[agentIdx] = uint(int(dice[agentIdx]) + change)
		}

		winnerIdx, hasWon = WinningPlayer(dice)
	}
	agents.SendGameComplete(ags, winnerIdx)
}

func PlayRound(round *game.Round, ags []agents.Agent) (game.Agent, int) {
	agents.SendRoundStarted(ags, *round)

	// TODO: Rotate who starts.
	// Consecutive agent's turn
	for true {
		for i, agent := range ags {
			e := game.Event{
				EType:     game.AgentTurn,
				AgentTurn: &game.AgentTurnEvent{},
			}
			agent.Handle(e)
			log.Debug().Interface("event", e).Send()

			// Ignore evicted players
			if len(round.Dice[i]) == 0 {
				log.Debug().Int("evicted player", i).Send()
				round.IncrementCurrentAgent()
				continue
			}
			act := agent.Play(*round)
			agentIdx, changeDice, err := round.PlayTurn(game.Agent(i), act)
			for err != nil {
				e := game.Event{EType: game.InvalidAction, InvalidAction: &game.InvalidActionEvent{
					InvalidAction: act,
					Err:           err.Error(),
				}}
				agent.Handle(e)
				log.Debug().Interface("event", e).Send()
				act = agent.Play(*round)
				agentIdx, changeDice, err = round.PlayTurn(game.Agent(i), act)
			}
			agents.SendTurnEvent(ags, act, game.Agent(i))
			round.IncrementCurrentAgent()

			if changeDice != 0 {
				agents.SendRoundComplete(ags, agentIdx, changeDice)
				return agentIdx, changeDice
			}
		}
	}
	return 0, 0
}

// WinningPlayer finds and, if exists, returns the index of the winning player (and whether there was a winning player).
func WinningPlayer(d []uint) (game.Agent, bool) {
	var winner game.Agent
	tot := uint(0)
	winnersDice := uint(0)

	for i, u := range d {
		tot += u
		if u != 0 {
			winner = game.Agent(i)
			winnersDice = u
		}
	}
	// Only winner has dice
	if tot == winnersDice {
		return winner, true
	}
	return winner, false
}
