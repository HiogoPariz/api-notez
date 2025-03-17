package api

import (
	"database/sql"

	"github.com/HiogoPariz/api-notez/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB) {
	router := gin.Default()
	noteRepo := repository.CreateNoteRepository(db)
	noteService := createNoteService(noteRepo)
	cacheStorage := createRedisStorage()

	router.Use(cors.Default())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessionMiddleware(cacheStorage))

	router.GET("/note", noteService.GetNotes)
	router.GET("/note/:id", noteService.GetNoteByID)
	router.GET("/note/usr/:userId", noteService.GetNoteByUserId)
	router.POST("/note", noteService.CreateNote)

	if err := router.Run(":3000"); err != nil {
		panic(err)
	}
}
