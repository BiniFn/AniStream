package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/coeeter/aniways/internal/worker/library"
	"github.com/spf13/cobra"
)

var libraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Library operations",
}

var retryFailedCmd = &cobra.Command{
	Use:   "retry-failed",
	Short: "Retry failed library syncs",
	Run: func(cmd *cobra.Command, args []string) {
		deps, err := initDeps()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing dependencies: %v\n", err)
			os.Exit(1)
		}
		defer deps.Close()

		ctx := context.Background()
		log := deps.Log.With("command", "library-retry-failed")

		library.RetryFailedLibrarySyncs(ctx, deps.Repo, deps.MAL, deps.Anilist, log)
	},
}

func init() {
	rootCmd.AddCommand(libraryCmd)
	libraryCmd.AddCommand(retryFailedCmd)
}