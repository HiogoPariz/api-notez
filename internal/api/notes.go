package api

import (
	"github.com/HiogoPariz/api-notez/internal/dto"
	"github.com/HiogoPariz/api-notez/internal/repository"
	"github.com/gin-gonic/gin"
)

type Notes interface {
	GetNotes(ctx *gin.Context)
	GetNote(ctx *gin.Context)
	DeleteNote(ctx *gin.Context)
	PostNote(ctx *gin.Context)
}

func GetNotes(ctx *gin.Context) {
	notes_dto, err := repository.NoteService.GetNotes()
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
