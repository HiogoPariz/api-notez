package main

import (
	"log"

	"github.com/HiogoPariz/api-notez/internal/api"
	"github.com/HiogoPariz/api-notez/internal/migration"
	"github.com/HiogoPariz/api-notez/internal/repository"
)

func main() {
	// Init db connection
	repo, err := repository.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	if err := migration.Run(repo); err != nil {
		log.Fatal(err)
	}

	// Start api server
	api.Init(repo)
}
