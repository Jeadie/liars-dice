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
			defer c.Close()
			if err != nil {
				log.Error().Err(err).Msg("failed to connect to server")
				return
			}

			// Listens for Events, sends actions
			events := make(chan game.Event)
			actions := make(chan game.Action)
			go network.Handle(c, events, actions)

			var round *game.Round
			idx := -1
			n := -1
			h := agents2.ConstructHuman()

			for e := range events {
				h.Handle(e)
				if e.EType == game.RoundStart {
					// Currently, 0 always starts. TODO: fix
					round = game.InitRound(e.RoundStart.DiceRolled, 0)
					if idx == 0 {
						actions <- h.Play(*round)
					}
				}
				if e.EType == game.GameStart {
					idx = e.GameStart.AgentIdx
					n = len(e.GameStart.NumDicePerAgent)
					log.Warn().Int("AgentIdx", idx).Msg("Game starting with AgentIdx")
				}
				if e.EType == game.Turn {
					// If previous turn was agent beforehand, play round.
					if game.Mod(int(e.Turn.ActionAgent)+1, n) == idx {
						actions <- h.Play(*round)
					}
					log.Debug().Interface("tur", e).Int("idx", idx).Int("v", game.Mod(int(e.Turn.ActionAgent)+1, n)).Send()
				}
			}

		},
	}
)

func init() {
	wsClientCmd.LocalFlags().StringVar(&wsServerAddr, "ws-server", ":8321", "The network address hosting the liars dice websocket game.")
	rootCmd.AddCommand(wsClientCmd)
	cobra.OnInitialize(initConfig)
}
