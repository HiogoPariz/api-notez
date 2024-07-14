package api

import (
	"github.com/HiogoPariz/api-notez/internal/repository"
	"github.com/gin-gonic/gin"
)

func Init(repo *repository.PostgresRepository) {
	router := gin.Default()

	router.GET("/note", GetNotes)

	if err := router.Run(":3000"); err != nil {
		panic(err)
	}
}
