package repository

import (
	"database/sql"
	"fmt"
	"github.com/HiogoPariz/api-notez/internal/dto"
)

type NoteRepository struct {
	DB *sql.DB
}

type INoteRepository interface {
	CreateNote(*dto.NoteDTO) error
	DeleteNote(int) error
	UpdateNote(*dto.NoteDTO) error
	GetNoteByID(int) (*dto.NoteDTO, error)
	GetNotes() ([]*dto.NoteDTO, error)
}

func CreateNoteRepository(db *sql.DB) INoteRepository {
	return &NoteRepository{DB: db}
}

func (repo *NoteRepository) CreateNote(note *dto.NoteDTO) error {
	query := `
	INSERT INTO note (title, file_name, active, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)
	`
	resp, err := repo.DB.Exec(
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

func (repo *NoteRepository) UpdateNote(note *dto.NoteDTO) error {
	query := `
	UPDATE note 
  SET title = $1,
      file_name = $2, 
      active = $3,
      created_at = $4,
      updated_at = $5
  WHERE id = $6
  )`
	resp, err := repo.DB.Exec(
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

func (repo *NoteRepository) DeleteNote(id int) error {
	_, err := repo.DB.Exec(`UPDATE note
	 SET active = false
	 WHERE id = $1
	`, id)

	return err
}

func (repo *NoteRepository) GetNoteByID(id int) (*dto.NoteDTO, error) {
	rows, err := repo.DB.Query("SELECT * FROM note n WHERE n.id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoNote(rows)
	}
	return nil, fmt.Errorf("note.%d not found", id)
}

func (repo *NoteRepository) GetNotes() ([]*dto.NoteDTO, error) {
	rows, err := repo.DB.Query("SELECT * FROM note n WHERE n.active is true")
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


func (repo *NoteRepository) GetNoteByUserId(userId int) (*dto.NoteListById, error) {
	rows, err := repo.DB.Query("SELECT id, title, created_at, updated_at FROM note n WHERE n.active = true AND n.user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []dto.Note

	for rows.Next() {
		var note dto.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	fileNameList := &dto.NoteListById{
		Notes: notes,
	}

	return fileNameList, nil
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
