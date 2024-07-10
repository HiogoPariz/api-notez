package storage

import (
	"database/sql"
	"fmt"

	"github.com/HiogoPariz/api-notez/internal/types"
	_ "github.com/lib/pq"
)

type Page = types.Page

type Store interface {
	Init() error
	CreatePage(*Page) error
	DeletePage(int) error
	UpdatePage(*Page) error
	GetPageByID(int) (*Page, error)
	GetPages() ([]*Page, error)
}

type PostgresStore struct {
	db *sql.DB
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
		db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createPagetable()
}

func (s *PostgresStore) createPagetable() error {
	query := `CREATE TABLE IF NOT EXISTS page (
		id serial primary key,
		title varchar(25),
		content varchar(50),
		active boolean,
		created_at timestamp,
		updated_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreatePage(page *Page) error {
	query := `
	INSERT INTO page (title, content, active, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)
	`
	resp, err := s.db.Exec(
		query,
		page.Title,
		page.Content,
		page.Active,
		page.CreatedAt,
		page.UpdatedAt)

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
	_, err := s.db.Exec(`UPDATE page
	 SET active = false
	 WHERE id = $1
	`, id)

	return err
}

func (s *PostgresStore) GetPageByID(id int) (*Page, error) {
	rows, err := s.db.Query("SELECT * FROM page p WHERE p.id = $1", id)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoPage(rows)
	}
	return nil, fmt.Errorf("page %d not found", id)
}

func (s *PostgresStore) GetPages() ([]*Page, error) {
	rows, err := s.db.Query("SELECT * FROM page p WHERE p.active = true")
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
	err := rows.Scan(&page.ID, &page.Title, &page.Content, &page.Active, &page.CreatedAt, &page.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &page, err
}
