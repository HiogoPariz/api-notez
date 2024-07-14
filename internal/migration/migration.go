package migration

import (
	"embed"

	"github.com/HiogoPariz/api-notez/internal/repository"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func Run(repo *repository.PostgresRepository) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(repository.GetDB(repo), "."); err != nil {
		return err
	}
	return nil
}
