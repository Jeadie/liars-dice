package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

type TestPaylod struct {
	Foo string `json:"foo"`
}

func SingleGameHandler(w http.ResponseWriter, r *http.Request) {
	x := TestPaylod{Foo: "bar"}
	b, err := json.Marshal(x)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(fmt.Errorf("Bad foo bar"))
		w.Write(b)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(b)
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
