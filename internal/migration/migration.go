package migration

import (
	"embed"

	"github.com/pressly/goose/v3"

	"github.com/HiogoPariz/api-notez/internal/storage"
)

//go:embed *.sql
var embedMigrations embed.FS

func Run(storage *storage.PostgresStore) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(storage.DB, "."); err != nil {
		return err
	}
	return nil
}
