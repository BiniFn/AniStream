package main

import (
	"log"

	"github.com/coeeter/aniways/internal/api"
	"github.com/coeeter/aniways/internal/config"
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

	srv := api.NewServer(env, db)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
