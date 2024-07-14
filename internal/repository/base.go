package repository

import "database/sql"

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresRepository, error) {
	connStr := "host=localhost user=postgres dbname=postgres password=api-notez sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db,
	}, nil
}

func GetDB(repo *PostgresRepository) *sql.DB {
	return repo.db
}
