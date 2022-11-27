package cmd

import (
	agents2 "github.com/Jeadie/liars-dice/pkg/agents"
	"github.com/Jeadie/liars-dice/pkg/game"
	"github.com/Jeadie/liars-dice/pkg/network"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	wsServerAddr string
	wsClientCmd  = &cobra.Command{
		Use:   "ws-client",
		Short: "Play a game of liar's dice",
		Long:  `Play a game of liar's dice`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create connection to server
			c, err := network.ConnectToServer(wsServerAddr)
			if err != nil {
				log.Error().Err(err).Msg("failed to connect to server")
				return
			}
			defer c.Close()

			// Listens for Events, sends actions
			events := make(chan game.Event)
			actions := make(chan game.Action)
			go network.Handle(c, events, actions)

			var round *game.Round
			h := agents2.ConstructHuman()

			for e := range events {
				h.Handle(e)
				if e.EType == game.RoundStart {
					// Currently, 0 always starts. TODO: fix
					round = game.InitRound(e.RoundStart.DiceRolled, 0)
				}
				if e.EType == game.AgentTurn {
					actions <- h.Play(*round)
				}

				if e.EType == game.GameComplete {
					close(events)
					close(actions)
					return
				}
			}

		},
	}
)

func init() {
	wsClientCmd.Flags().StringVar(&wsServerAddr, "ws-server", "ws://0.0.0.0:8321", "The network address hosting the liars dice websocket game.")
	rootCmd.AddCommand(wsClientCmd)
	cobra.OnInitialize(initConfig)
}
