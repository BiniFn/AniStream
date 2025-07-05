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

	srv := api.NewServer(env)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
