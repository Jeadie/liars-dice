package cmd

import (
	"fmt"
	agents2 "github.com/Jeadie/liars-poker/pkg/agents"
	"github.com/Jeadie/liars-poker/pkg/game"
	"os"
	"strconv"
)

func MakeAgents(n uint, humanIdx uint) []agents2.Agent {
	agents := make([]agents2.Agent, n)
	for i := uint(0); i < n; i++ {
		if i == humanIdx {
			agents[i] = agents2.ConstructHuman()
		} else {
			agents[i] = agents2.ConstructProbAgent()
		}
	}
	return agents
}

func PlayRound(round *game.Round, agents []agents2.Agent) (game.Agent, int) {
	for i, agent := range agents {
		agent.Handle(game.Event{
			EType: game.RoundStart,
			RoundStart: game.RoundStartEvent{
				DiceRolled: round.Dice[i],
			},
		})
	}

	// Consecutive agent's turn
	for true {
		for i, agent := range agents {
			// Ignore evicted players
			if len(round.Dice[i]) == 0 {
				continue
			}
			act := agent.Play(*round)
			agentIdx, changeDice, err := round.PlayTurn(game.Agent(i), act)
			for err != nil {
				fmt.Printf("%s\n", err)
				act := agent.Play(*round)
				agentIdx, changeDice, err = round.PlayTurn(game.Agent(i), act)
			}
			fmt.Printf("Player %d %s\n", i, act.ToString())

			if changeDice != 0 {
				if changeDice > 0 {
					fmt.Printf("Player %d gains %d dice\n", agentIdx, changeDice)
				} else {
					fmt.Printf("Player %d loses %d dice(s)\n", agentIdx, -1*changeDice)
				}
				return agentIdx, changeDice
			}
		}
		fmt.Println()
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
func WinningPlayer(d []uint) (uint, bool) {
	var winner uint
	tot := uint(0)
	winnersDice := uint(0)

	for i, u := range d {
		tot += u
		if u != 0 {
			winner = uint(i)
			winnersDice = u
		}
	}
	// Only winner has dice
	if tot == winnersDice {
		return winner, true
	}
	return winner, false
}
