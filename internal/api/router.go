package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB) {
	router := gin.Default()
	noteService := CreateNoteService(db)
	router.GET("/note", noteService.GetNotes)
	router.GET("/note/:id", noteService.GetNoteByID)
	router.POST("/note", noteService.CreateNote)

	if err := router.Run(":3000"); err != nil {
		panic(err)
	}
}
