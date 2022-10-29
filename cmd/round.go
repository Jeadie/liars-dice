package cmd

import (
	"fmt"
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

			agents := MakeAgents(uint(len(dice)), humanAgent)
			PlayRound(round, agents)
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
