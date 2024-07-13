package main

import (
	"log"

	"github.com/HiogoPariz/api-notez/internal/api"
	"github.com/HiogoPariz/api-notez/internal/migration"
	"github.com/HiogoPariz/api-notez/internal/storage"
)

func main() {
	// Init db connection
	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := migration.Run(store); err != nil {
		log.Fatal(err)
	}

	// Start api server
	server := api.NewAPIServer(":3000", store)
	server.Run()
}
