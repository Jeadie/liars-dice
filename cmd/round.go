package cmd

import (
	"fmt"
	agents2 "github.com/Jeadie/liars-dice/pkg/agents"
	"github.com/Jeadie/liars-dice/pkg/game"
	"github.com/spf13/cobra"
	"os"
)

var (
	roundCmd = &cobra.Command{
		Use:   "round",
		Short: "Play a round of Liar's dice",
		Long:  `Play a round of Liar's dice`,
		Run: func(cmd *cobra.Command, args []string) {
			if humanAgent >= len(args) {
				fmt.Printf("for a round with dice %s per player, 0 <= idx < %d\n", args, len(args))
				os.Exit(1)
			}
			dice := ConvertNumDice(args)
			round := game.InitRound(dice, 0)
			agents := MakeAgents(uint(len(dice)), humanAgent, agents2.ConstructWsAgents(wsAddr, socketAgents), socketAgents)

			agents2.SendGameStarted(agents, dice)
			PlayRound(round, agents)
		},
	}
)

func init() {
	rootCmd.AddCommand(roundCmd)
}
