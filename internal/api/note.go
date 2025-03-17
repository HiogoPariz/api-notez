package api

import (
	"crypto/rand"
	"encoding/hex"
	"hash/fnv"
	"net/http"
	"strconv"
	"time"

	"github.com/HiogoPariz/api-notez/internal/dto"
	"github.com/HiogoPariz/api-notez/internal/integration"
	"github.com/HiogoPariz/api-notez/internal/repository"
	"github.com/gin-gonic/gin"
)

type NoteService struct {
	repo repository.INoteRepository
}

type Notes interface {
	GetNotes(ctx *gin.Context)
	GetNote(ctx *gin.Context)
	DeleteNote(ctx *gin.Context)
	PostNote(ctx *gin.Context)
	GetNoteByUserId(ctx *gin.Context)
}

func createNoteService(repo repository.INoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (service NoteService) GetNotes(ctx *gin.Context) {
	notes_dto, err := service.repo.GetNotes()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	notes := []*dto.NoteObject{}

	for _, note_dto := range notes_dto {
		integration := integration.CreateFileIntegration(note_dto)

		note_res, err := note_dto.EntityToResponse()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		content, err := integration.GetFileContent()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		note_res.Content = content

		notes = append(notes, note_res)
	}

	ctx.JSON(http.StatusOK, notes)
}

func (service NoteService) GetNoteByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
	}

	note_dto, err := service.repo.GetNoteByID(id)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
	}

	note, err := note_dto.EntityToResponse()
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
	}

	notesIntegration := integration.CreateFileIntegration(note_dto)
	content, err := notesIntegration.GetFileContent()
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
	}
	note.Content = content

	ctx.JSON(http.StatusOK, note)
}

func (service NoteService) CreateNote(ctx *gin.Context) {
	body := &dto.NoteRequest{}
	if err := ctx.Bind(body); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
	}

	fileName, err := generateFileName()

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	noteDto := body.RequestToDTO(fileName)

	// Grava na base
	if err := service.repo.CreateNote(noteDto); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Manda pro files
	notesIntegration := integration.CreateFileIntegration(noteDto)
	if err := notesIntegration.CreateFileContent(body.Content, fileName); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "Note created successfully",
		"fileName": fileName,
	})
}

func (service NoteService) GetNoteByUserId(ctx *gin.Context) {
	userIdStr := ctx.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	notes, err := service.repo.GetNoteByUserId(userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, notes)
}

func generateFileName() (string, error) {

	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	randomHex := hex.EncodeToString(randomBytes)

	data := time.Now().String() + randomHex

	hash := fnv.New64a()
	if _, err := hash.Write([]byte(data)); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
