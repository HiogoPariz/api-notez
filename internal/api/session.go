package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type SessionService struct {
	Store *redis.Client
}

type Session interface {
	GetSession(ctx *gin.Context)
	CreateSession(ctx *gin.Context)
}

type RedisCache struct {
	Store            *redis.Client
	DefaultCacheTime time.Duration
}

func createRedisStorage() *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisCache{
		Store:            rdb,
		DefaultCacheTime: 60 * time.Minute,
	}
}

func sessionMiddleware(cache *RedisCache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := cache.Store.Get(ctx, "session").Result()
		if err == redis.Nil {
			fmt.Println("No session found.")
		} else if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		} else {
			ctx.Set("session", session)
		}

		ctx.Next()
	}
}

func (service SessionService) GetSession(ctx *gin.Context) {

}
