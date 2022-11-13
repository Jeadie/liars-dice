package cmd

import (
	"fmt"
	agents2 "github.com/Jeadie/liars-dice/pkg/agents"
	"github.com/Jeadie/liars-dice/pkg/game"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func MakeAgents(n uint, humanIdx int, wsAgents chan *agents2.WsAgent, numWsAgents uint) []agents2.Agent {
	agents := make([]agents2.Agent, n)
	for i := 0; i < int(n); i++ {
		// If humanIdx == -1, no human agent is used.
		if i == humanIdx {
			agents[i] = agents2.ConstructHuman()

		} else if numWsAgents > 0 {
			numWsAgents--
			x := <-wsAgents
			agents[i] = x
		} else {
			agents[i] = agents2.ConstructProbAgent()
		}
	}
	return agents
}

func PlayRound(round *game.Round, agents []agents2.Agent) (game.Agent, int) {
	game.SendRoundStarted(agents, *round)

	// TODO: Rotate who starts.
	// Consecutive agent's turn
	for true {
		for i, agent := range agents {
			e := game.Event{
				EType:     game.AgentTurn,
				AgentTurn: &game.AgentTurnEvent{},
			}
			agent.Handle(e)
			log.Debug().Interface("event", e).Send()

			// Ignore evicted players
			if len(round.Dice[i]) == 0 {
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
				act := agent.Play(*round)
				agentIdx, changeDice, err = round.PlayTurn(game.Agent(i), act)
			}
			game.SendTurnEvent(agents, act, game.Agent(i))

			if changeDice != 0 {
				game.SendRoundComplete(agents, agentIdx, changeDice)
				return agentIdx, changeDice
			}
		}
	}
	return 0, 0
}

func ConvertNumDice(numDice []string) []uint {
	dice := make([]uint, len(numDice))
	for i, die := range numDice {
		v, err := strconv.Atoi(die)
		if err != nil || v <= 0 {
			fmt.Printf("%s is not a valid, positive, integer", die)
			os.Exit(1)
		}
		dice[i] = uint(v)
	}
	return dice
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
