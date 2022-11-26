package cmd

import (
	"encoding/json"
	"github.com/Jeadie/liars-dice/pkg/agents"
	"github.com/Jeadie/liars-dice/pkg/network"
	"github.com/Jeadie/liars-dice/pkg/play"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

func SingleGameHandler(w http.ResponseWriter, r *http.Request) {
	e := hlog.FromRequest(r).Debug()
	defer e.Send()

	conn, err := network.UpgradeToWebsockets(w, r)

	if err != nil {
		w.WriteHeader(http.StatusTeapot)
		b, _ := json.Marshal(err)
		w.Write(b)
		e = e.Err(err)
	} else {
		agent := agents.CreateWsAgent(conn)
		go play.PlayGame(
			[]agents.Agent{agent, agents.ConstructProbAgent(), agents.ConstructProbAgent()},
			[]uint{3, 3, 3},
		)
		w.WriteHeader(http.StatusOK)
	}
}

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Liar's dice server",
		Long:  `Run a server that hosts multiple games of liar's dice.`,
		Run: func(cmd *cobra.Command, args []string) {
			r := mux.NewRouter()
			r.HandleFunc("/single", SingleGameHandler)
			srv := &http.Server{
				Handler:      r,
				Addr:         wsAddr,
				WriteTimeout: 30 * time.Second,
				ReadTimeout:  30 * time.Second,
			}
			err := srv.ListenAndServe()
			if err != nil {
				log.Error().Err(err).Msg("Problem running server")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)
}
