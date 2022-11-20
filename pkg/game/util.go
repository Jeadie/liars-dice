package game

import (
	"github.com/Jeadie/liars-dice/pkg/agents"
	"github.com/rs/zerolog/log"
)

func PlayGame(ags []agents.Agent, dice []uint) {
	agents.SendGameStarted(ags, dice)
	round := InitRound(dice, 0)
	winnerIdx, hasWon := WinningPlayer(dice)
	for !hasWon {
		agentIdx, change := PlayRound(round, ags)

		// Changes score from last
		if int(dice[agentIdx])+change <= 0 {
			dice[agentIdx] = 0
		} else {
			dice[agentIdx] = uint(int(dice[agentIdx]) + change)
		}

		winnerIdx, hasWon = WinningPlayer(dice)
		round = InitRound(dice, 0)
	}
	agents.SendGameComplete(ags, winnerIdx)
}

func PlayRound(round *Round, ags []agents.Agent) (Agent, int) {
	agents.SendRoundStarted(ags, *round)

	// TODO: Rotate who starts.
	// Consecutive agent's turn
	for true {
		for i, agent := range ags {
			e := Event{
				EType:     AgentTurn,
				AgentTurn: &AgentTurnEvent{},
			}
			agent.Handle(e)
			log.Debug().Interface("event", e).Send()

			// Ignore evicted players
			if len(round.Dice[i]) == 0 {
				continue
			}
			act := agent.Play(*round)
			agentIdx, changeDice, err := round.PlayTurn(Agent(i), act)
			for err != nil {
				e := Event{EType: InvalidAction, InvalidAction: &InvalidActionEvent{
					InvalidAction: act,
					Err:           err.Error(),
				}}
				agent.Handle(e)
				log.Debug().Interface("event", e).Send()
				act := agent.Play(*round)
				agentIdx, changeDice, err = round.PlayTurn(Agent(i), act)
			}
			agents.SendTurnEvent(ags, act, Agent(i))

			if changeDice != 0 {
				agents.SendRoundComplete(ags, agentIdx, changeDice)
				return agentIdx, changeDice
			}
		}
	}
	return 0, 0
}

// WinningPlayer finds and, if exists, returns the index of the winning player (and whether there was a winning player).
func WinningPlayer(d []uint) (Agent, bool) {
	var winner Agent
	tot := uint(0)
	winnersDice := uint(0)

	for i, u := range d {
		tot += u
		if u != 0 {
			winner = Agent(i)
			winnersDice = u
		}
	}
	// Only winner has dice
	if tot == winnersDice {
		return winner, true
	}
	return winner, false
}
