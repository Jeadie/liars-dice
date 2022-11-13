package cmd

import (
	agents2 "github.com/Jeadie/liars-dice/pkg/agents"
	"github.com/Jeadie/liars-dice/pkg/game"
	"github.com/Jeadie/liars-dice/pkg/network"
	"github.com/spf13/cobra"
	"log"
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
				log.Fatal("dial:", err)
				return
			}

			// Listens for Events, sends actions
			events := make(chan game.Event)
			actions := make(chan game.Action)
			go network.Handle(c, events, actions)

			var round *game.Round
			idx := -1
			h := agents2.ConstructHuman()

			for e := range events {
				if e.EType == game.RoundStart {
					// Currently, 0 always starts.
					round = game.InitRound(e.RoundStart.DiceRolled, 0)
				}
				if e.EType == game.GameStart {
					idx = e.GameStart.AgentIdx
				}
				if e.EType == game.Turn {
					// If previous turn was agent beforehand, play round.
					if int(e.Turn.ActionAgent)+1 == idx {
						h.Play(*round)
					}
				}
				h.Handle(e)
			}

		},
	}
)

func init() {
	wsClientCmd.LocalFlags().StringVar(&wsServerAddr, "ws-server", ":8123", "The network address hosting the liars dice websocket game.")
	rootCmd.AddCommand(wsClientCmd)
}
