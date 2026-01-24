package cli

import (
	"github.com/coeeter/aniways/internal/worker/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Auth operations",
}

var refreshTokensCmd = &cobra.Command{
	Use:   "refresh-tokens",
	Short: "Refresh OAuth tokens",
	Run: func(cmd *cobra.Command, args []string) {
		log := deps.Log.With("command", "auth-refresh-tokens")
		auth.DailyTask(cmd.Context(), deps.Repo, deps.Providers, log)
	},
}

func init() {
	authCmd.AddCommand(refreshTokensCmd)
}
