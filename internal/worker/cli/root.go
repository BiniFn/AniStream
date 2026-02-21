package cli

import (
	"fmt"
	"os"

	"github.com/coeeter/aniways/internal/app"
	"github.com/spf13/cobra"
)

var deps *app.Deps

func Execute() {
	var rootCmd = &cobra.Command{
		Use:               "worker",
		Short:             "AniStream worker",
		PersistentPreRunE: initDepsOnce,
	}

	rootCmd.AddCommand(daemonCmd, authCmd, scrapeCmd, libraryCmd)

	if len(os.Args) == 1 {
		rootCmd.SetArgs([]string{"daemon"})
	}

	err := rootCmd.Execute()

	if deps != nil {
		deps.Close()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func initDepsOnce(cmd *cobra.Command, args []string) error {
	if deps != nil {
		return nil
	}

	d, err := app.InitDeps(cmd.Context(), "WORKER")
	if err != nil {
		return fmt.Errorf("error initializing dependencies: %w", err)
	}

	deps = d
	return nil
}
