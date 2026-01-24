package cli

import (
	"github.com/coeeter/aniways/internal/worker/scraper"
	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape anime data",
}

var recentlyUpdatedCmd = &cobra.Command{
	Use:   "recently-updated",
	Short: "Scrape recently updated anime",
	Run: func(cmd *cobra.Command, args []string) {
		log := deps.Log.With("command", "scrape-recently-updated")

		scraper.HourlyTask(cmd.Context(), deps.Scraper, deps.Repo, deps.Cache, log)
	},
}

var fullSeedCmd = &cobra.Command{
	Use:   "full-seed",
	Short: "Full database seed",
	Run: func(cmd *cobra.Command, args []string) {
		log := deps.Log.With("command", "scrape-full-seed")

		if err := scraper.FullSeed(cmd.Context(), deps.Scraper, deps.Repo, log); err != nil {
			log.Error("Full seed failed", "err", err)
			return
		}
	},
}

var allRecentlyUpdatedCmd = &cobra.Command{
	Use:   "all-recently-updated",
	Short: "Scrape all recently updated pages",
	Run: func(cmd *cobra.Command, args []string) {
		log := deps.Log.With("command", "scrape-all-recently-updated")

		if err := scraper.ScrapeAllRecentlyUpdated(cmd.Context(), deps.Scraper, deps.Repo, log); err != nil {
			log.Error("All recently updated scrape failed", "err", err)
			return
		}
	},
}

func init() {
	scrapeCmd.AddCommand(recentlyUpdatedCmd)
	scrapeCmd.AddCommand(allRecentlyUpdatedCmd)
	scrapeCmd.AddCommand(fullSeedCmd)
}
