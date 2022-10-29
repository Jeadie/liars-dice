package cmd

import (
	"fmt"
	"github.com/Jeadie/liars-poker/pkg/game"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "game",
	Short: "Play a liar's dice",
	Long:  `Play a liar's dice`,
	Run: func(cmd *cobra.Command, args []string) {
		if humanAgent >= uint(len(args)) {
			fmt.Printf("for a round with dice %s per player, 0 <= idx < %d\n", args, len(args))
			os.Exit(1)
		}
		dice := ConvertNumDice(args)
		round := game.InitRound(dice, 0)
		agents := MakeAgents(uint(len(dice)), humanAgent)

		winnerIdx, hasWon := WinningPlayer(dice)
		for !hasWon {
			fmt.Printf("\n-- New Round -- \n")
			agentIdx, change := PlayRound(round, agents)

			// Changes score from last game.
			if int(dice[agentIdx])+change <= 0 {
				dice[agentIdx] = 0
			} else {
				dice[agentIdx] = uint(int(dice[agentIdx]) + change)
			}

			winnerIdx, hasWon = WinningPlayer(dice)
			round = game.InitRound(dice, 0)
		}
		fmt.Printf("\n Player %d wins!\n", winnerIdx)
	},
}

func init() {
	rootCmd.PersistentFlags().UintVar(&humanAgent, "idx", 0, "The index of the user in the order of players.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
