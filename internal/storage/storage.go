package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/HiogoPariz/api-notez/internal/types"
)

type Page = types.Page

type Store interface {
	CreatePage(*Page) error
	DeletePage(int) error
	UpdatePage(*Page) error
	GetPageByID(int) (*Page, error)
	GetPages() ([]*Page, error)
}

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=api-notez sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		DB: db,
	}, nil
}

func (s *PostgresStore) CreatePage(page *Page) error {
	query := `
	INSERT INTO page (title, content, active, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)
	`
	resp, err := s.DB.Exec(
		query,
		page.Title,
		page.Content,
		page.Active,
		page.CreatedAt,
		page.UpdatedAt,
	)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *PostgresStore) UpdatePage(*Page) error {
	return nil
}

func (s *PostgresStore) DeletePage(id int) error {
	_, err := s.DB.Exec(`UPDATE page
	 SET active = false
	 WHERE id = $1
	`, id)

	return err
}

func (s *PostgresStore) GetPageByID(id int) (*Page, error) {
	rows, err := s.DB.Query("SELECT * FROM page p WHERE p.id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoPage(rows)
	}
	return nil, fmt.Errorf("page %d not found", id)
}

func (s *PostgresStore) GetPages() ([]*Page, error) {
	rows, err := s.DB.Query("SELECT * FROM page p WHERE p.active = true")
	fmt.Printf("%+v\n", rows)
	if err != nil {
		return nil, err
	}

	pages := []*Page{}

	for rows.Next() {
		page, err := scanIntoPage(rows)
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func scanIntoPage(rows *sql.Rows) (*Page, error) {
	page := Page{}
	err := rows.Scan(
		&page.ID,
		&page.Title,
		&page.Content,
		&page.Active,
		&page.CreatedAt,
		&page.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &page, err
}
