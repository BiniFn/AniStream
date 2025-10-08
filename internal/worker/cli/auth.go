package cli

import (
	"context"
	"fmt"
	"os"

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
		deps, err := initDeps()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing dependencies: %v\n", err)
			os.Exit(1)
		}
		defer deps.Close()

		ctx := context.Background()
		log := deps.Log.With("command", "auth-refresh-tokens")

		auth.DailyTask(ctx, deps.Repo, deps.Providers, log)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(refreshTokensCmd)
}