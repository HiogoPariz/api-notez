package api

import (
	"database/sql"
	"encoding/hex"
	"hash/fnv"
	"strconv"
	"time"

	"github.com/HiogoPariz/api-notez/internal/dto"
	"github.com/HiogoPariz/api-notez/internal/integration"
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

	notesIntegration := integration.CreateFileIntegration(note_dto)
	content, err := notesIntegration.GetFileContent()
	if err != nil {
		ctx.AbortWithError(505, err)
	}
	note.Content = content

	ctx.JSON(200, note)
}

func (service *NoteService) CreateNote(ctx *gin.Context) {
	repo := repository.CreateNoteRepository(service.DB)
	hash := fnv.New64a()

	body := &dto.NoteRequest{}
	if err := ctx.Bind(body); err != nil {
		ctx.AbortWithError(400, err)
	}

	// Cria hash do file name, por enquanto ta com a data
	if _, err := hash.Write([]byte(time.Now().String())); err != nil {
		ctx.AbortWithError(400, err)
	}

	if _, err := hash.Write([]byte("note/json")); err != nil {
		ctx.AbortWithError(400, err)
	}
	file_name := hex.EncodeToString(hash.Sum(nil))
	note_dto := body.RequestToDTO(file_name)

	// Grava na base
	if err := repo.CreateNote(note_dto); err != nil {
		ctx.AbortWithError(505, err)
	}

	// Manda pro files
	notesIntegration := integration.CreateFileIntegration(note_dto)
	if err := notesIntegration.CreateFileContent(body.Content, file_name); err != nil {
		ctx.AbortWithError(505, err)
	}
}
