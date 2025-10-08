package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/coeeter/aniways/internal/worker"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run daemon mode",
	Run: func(cmd *cobra.Command, args []string) {
		deps, err := initDeps()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing dependencies: %v\n", err)
			os.Exit(1)
		}
		defer deps.Close()

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		mgr := worker.NewManager(
			deps.Db,
			deps.Repo,
			deps.Scraper,
			deps.MAL,
			deps.Anilist,
			deps.Cache,
			deps.Log.With("component", "worker"),
		)

		if err := mgr.Bootstrap(ctx); err != nil {
			deps.Log.Error("Error in bootstrapping:", "err", err)
			os.Exit(1)
		}

		deps.Log.Info("Starting worker daemon...")
		mgr.StartBackground(ctx, deps.Providers)

		<-ctx.Done()
		deps.Log.Info("Worker daemon stopped")
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}

