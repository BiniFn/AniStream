package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/coeeter/aniways/internal/api"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/hianime"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/worker"
)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	db, err := config.NewDatabase(env)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	repo := repository.New(db)
	scraper := hianime.NewHianimeScraper()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	count, err := repo.GetCountOfAnimes(ctx)
	if err != nil {
		log.Fatalf("Error counting anime rows: %v", err)
	}

	if count == 0 {
		log.Println("⚙️  no anime in DB—running initial scraper (this will block)…")
		if err := worker.FullSeed(ctx, scraper, repo); err != nil {
			log.Fatalf("🚨 FullSeed failed: %v", err)
		}
		log.Println("✅ initial scrape complete, starting HTTP server")
	} else {
		log.Printf("ℹ️  DB already has %d anime, skipping initial scrape", count)
	}

	go worker.HourlyTask(ctx, scraper, repo)

	srv := api.NewServer(env, repo)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
