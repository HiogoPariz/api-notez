package dto

import (
	"time"

	"github.com/jinzhu/copier"
)

type NoteDTO struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	FileName  string    `json:"file_name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteRequest struct {
	Title    string `json:"title"`
	FileName string `json:"file_name"`
}

type NoteObject struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	FileName string `json:"file_name"`
}

func (note_dto *NoteDTO) EntityToResponse() (*NoteObject, error) {
	note_response := &NoteObject{}

	if err := copier.Copy(note_response, note_dto); err != nil {
		return nil, err
	}

	return note_response, nil
}

func (note_req *NoteRequest) RequestToDTO() *NoteDTO {
	return &NoteDTO{
		Title:     note_req.Title,
		FileName:  note_req.FileName,
		Active:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
