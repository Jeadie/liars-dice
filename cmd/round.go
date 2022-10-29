package cmd

import (
	"fmt"
	agents2 "github.com/Jeadie/liars-poker/pkg/agents"
	"github.com/Jeadie/liars-poker/pkg/game"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var (
	humanAgent uint
	roundCmd   = &cobra.Command{
		Use:   "round",
		Short: "Play a round of Liar's dice",
		Long:  `Play a round of Liar's dice`,
		Run: func(cmd *cobra.Command, args []string) {
			if humanAgent >= uint(len(args)) {
				fmt.Printf("for a round with dice %s per player, 0 <= idx < %d", args, len(args))
				os.Exit(1)
			}
			dice := ConvertNumDice(args)
			round := game.InitRound(dice, 0)

			// Create and initialise agents
			agents := make([]agents2.Agent, len(dice))
			for i := 0; i < len(dice); i++ {
				if uint(i) == humanAgent {
					agents[i] = agents2.ConstructHuman()
				} else {
					agents[i] = agents2.ConstructProbAgent()
				}
				agents[i].Initialise(*round, round.Dice[i])
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
						os.Exit(0)
					}
				}
			}
		},
	}
)

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

func init() {
	rootCmd.AddCommand(roundCmd)
	roundCmd.PersistentFlags().UintVar(&humanAgent, "idx", 0, "The index of the user in the order of players.")
}
