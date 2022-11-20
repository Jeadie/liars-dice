package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Liar's dice server",
		Long: `Run a server that hosts multiple games of liar's dice. 
		`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Warn().Msg("server")
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)
}
