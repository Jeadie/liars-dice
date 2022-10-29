package cmd

import (
	"fmt"
	agents2 "github.com/Jeadie/liars-poker/pkg/agents"
	"github.com/Jeadie/liars-poker/pkg/game"
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
		agent.Initialise(*round, round.Dice[i])
	}

	// Consecutive agent's turn
	for true {
		for i, agent := range agents {
			act := agent.Play(*round)
			agentIdx, changeDice, err := round.PlayTurn(game.Agent(i), act)
			for err != nil {
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
	}
	return 0, 0
}
