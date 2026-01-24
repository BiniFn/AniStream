package cli

import (
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
		log := deps.Log.With("command", "library-retry-failed")

		library.RetryFailedLibrarySyncs(cmd.Context(), deps.Repo, deps.MAL, deps.Anilist, log)
	},
}

func init() {
	libraryCmd.AddCommand(retryFailedCmd)
}
