package api

import (
	"database/sql"
	"strconv"

	"github.com/HiogoPariz/api-notez/internal/dto"
	"github.com/HiogoPariz/api-notez/internal/repository"
	"github.com/gin-gonic/gin"
)

type NoteService struct {
	DB *sql.DB
}

type Notes interface {
	GetNotes(ctx *gin.Context)
	GetNote(ctx *gin.Context)
	DeleteNote(ctx *gin.Context)
	PostNote(ctx *gin.Context)
}

func CreateNoteService(db *sql.DB) *NoteService {
	return &NoteService{DB: db}
}

func (service *NoteService) GetNotes(ctx *gin.Context) {
	repo := repository.CreateNoteRepository(service.DB)
	notes_dto, err := repo.GetNotes()
	if err != nil {
		ctx.AbortWithError(505, err)
		return
	}

	notes := []*dto.NoteObject{}

	for _, note_dto := range notes_dto {
		note_res, err := note_dto.EntityToResponse()
		if err != nil {
			ctx.AbortWithError(505, err)
		}
		notes = append(notes, note_res)
	}

	ctx.JSON(200, notes)
}

func (service *NoteService) GetNoteByID(ctx *gin.Context) {
	repo := repository.CreateNoteRepository(service.DB)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(404, err)
	}

	note_dto, err := repo.GetNoteByID(id)
	if err != nil {
		ctx.AbortWithError(505, err)
	}

	note, err := note_dto.EntityToResponse()
	if err != nil {
		ctx.AbortWithError(505, err)
	}

	ctx.JSON(200, note)
}

func (service *NoteService) CreateNote(ctx *gin.Context) {
	repo := repository.CreateNoteRepository(service.DB)

	body := &dto.NoteRequest{}
	if err := ctx.Bind(body); err != nil {
		ctx.AbortWithError(400, err)
	}
	note_dto := body.RequestToDTO()

	if err := repo.CreateNote(note_dto); err != nil {
		ctx.AbortWithError(505, err)
	}
}
