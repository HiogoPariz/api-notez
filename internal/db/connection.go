package db

import "database/sql"

func NewPostgresStore() (*sql.DB, error) {
	connStr := "host=localhost user=postgres dbname=postgres password=api-notez sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
