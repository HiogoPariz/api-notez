package repository

import (
	"database/sql"
	"fmt"

	"github.com/HiogoPariz/api-notez/internal/dto"
	_ "github.com/lib/pq"
)

var NoteService NoteRepository = &RepoService{}

type RepoService struct {
	PostgresRepository
}

type NoteRepository interface {
	CreateNote(*dto.NoteDTO) error
	DeleteNote(int) error
	UpdateNote(*dto.NoteDTO) error
	GetNoteByID(int) (*dto.NoteDTO, error)
	GetNotes() ([]*dto.NoteDTO, error)
}

func (service *RepoService) CreateNote(note *dto.NoteDTO) error {
	query := `
	INSERT INTO note (title, content, active, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)
	`
	resp, err := GetDB(&service.PostgresRepository).Exec(
		query,
		note.Title,
		note.FileName,
		note.Active,
		note.CreatedAt,
		note.UpdatedAt,
	)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (service *RepoService) UpdateNote(note *dto.NoteDTO) error {
	query := `
	UPDATE note 
  SET title = $1,
      content = $2, 
      active = $3,
      created_at = $4,
      updated_at = $5
  WHERE id = $6
  )`
	resp, err := GetDB(&service.PostgresRepository).Exec(
		query,
		note.Title,
		note.FileName,
		note.Active,
		note.CreatedAt,
		note.UpdatedAt,
		note.ID,
	)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (service *RepoService) DeleteNote(id int) error {
	_, err := GetDB(&service.PostgresRepository).Exec(`UPDATE note
	 SET active = false
	 WHERE id = $1
	`, id)

	return err
}

func (service *RepoService) GetNoteByID(id int) (*dto.NoteDTO, error) {
	rows, err := GetDB(
		&service.PostgresRepository,
	).Query("SELECT * FROM note n WHERE n.id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoNote(rows)
	}
	return nil, fmt.Errorf("note.%d not found", id)
}

func (service *RepoService) GetNotes() ([]*dto.NoteDTO, error) {
	rows, err := GetDB(
		&service.PostgresRepository,
	).Query("SELECT * FROM note n WHERE n.active = true")
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", rows)

	notes := []*dto.NoteDTO{}

	for rows.Next() {
		note, err := scanIntoNote(rows)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func scanIntoNote(rows *sql.Rows) (*dto.NoteDTO, error) {
	note := dto.NoteDTO{}
	err := rows.Scan(
		&note.ID,
		&note.Title,
		&note.FileName,
		&note.Active,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &note, err
}
