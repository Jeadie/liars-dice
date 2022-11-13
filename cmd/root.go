package cmd

import (
	"fmt"
	agents2 "github.com/Jeadie/liars-dice/pkg/agents"
	"github.com/Jeadie/liars-dice/pkg/game"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"os"
)

var (
	humanAgent   int
	socketAgents uint
	wsAddr       string
	logLevel     string
	rootCmd      = &cobra.Command{
		Use:   "game",
		Short: "Play a liar's dice",
		Long:  `Play a liar's dice`,
		Run: func(cmd *cobra.Command, args []string) {
			if humanAgent >= len(args) {
				fmt.Printf("for a round with dice %s per player, 0 <= idx < %d\n", args, len(args))
				os.Exit(1)
			}
			if socketAgents >= uint(len(args)) {
				fmt.Printf("for a round with dice %s per player, 0 <= ws-agents < %d\n", args, len(args))
				os.Exit(1)
			}

			dice := ConvertNumDice(args)
			round := game.InitRound(dice, 0)
			agents := MakeAgents(uint(len(dice)), humanAgent, agents2.ConstructWsAgents(wsAddr, socketAgents), socketAgents)

			agents2.SendGameStarted(agents, dice)

			winnerIdx, hasWon := WinningPlayer(dice)
			for !hasWon {
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
			agents2.SendGameComplete(agents, winnerIdx)
		},
	}
)

func initConfig() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		rootCmd.PrintErr(err)
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(level)
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().IntVar(&humanAgent, "user-idx", -1, "The index of the user in the order of players. Default to -1; no human user.")
	rootCmd.PersistentFlags().StringVar(&wsAddr, "ws-addr", ":8321", "The network address to wait for users communicating over web sockets.")
	rootCmd.PersistentFlags().UintVar(&socketAgents, "ws-agents", 0, "The number of users to wait for, over websocket, to play in the game.")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", zerolog.WarnLevel.String(), "Level of logging to stderr. Levels: trace, debug, info, warn, error, fatal, panic")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
