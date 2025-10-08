package cli

import (
	"context"
	"os"

	"github.com/coeeter/aniways/internal/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "worker",
	Short: "Aniways worker",
	Run: func(cmd *cobra.Command, args []string) {
		daemonCmd.Run(cmd, args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initDeps() (*app.Deps, error) {
	ctx := context.Background()
	return app.InitDeps(ctx, "WORKER")
}

