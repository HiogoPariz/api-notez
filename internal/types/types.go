package types

import (
	"time"
)

type CreatePageRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Page struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewPage(title, content string) *Page {
	return &Page{
		Title:     title,
		Content:   content,
		Active:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
