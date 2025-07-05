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

	api.StartServer(env)
}
