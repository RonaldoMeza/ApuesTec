package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/config"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/database"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/middleware"
	appredis "github.com/RonaldoMeza/ApuesTec/backend/internal/redis"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/response"
)

func NewRouter(cfg *config.Config, db *pgxpool.Pool, redisClient *redis.Client) *gin.Engine {
	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.CORS(cfg.CORSAllowedOrigins),
		middleware.SecurityHeaders(),
	)

	api := router.Group(cfg.APIPrefix)
	{
		api.GET("/health", health(cfg))
		api.GET("/health/dependencies", dependenciesHealth(db, redisClient))
	}

	return router
}

func health(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"service": cfg.AppName,
			"status":  "ok",
		})
	}
}

func dependenciesHealth(db *pgxpool.Pool, redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), 3*time.Second)
		defer cancel()

		status := map[string]string{
			"postgres": "ok",
			"redis":    "ok",
		}

		if err := database.Ping(checkCtx, db); err != nil {
			status["postgres"] = "unavailable"
		}

		if err := appredis.Ping(checkCtx, redisClient); err != nil {
			status["redis"] = "unavailable"
		}

		if status["postgres"] != "ok" || status["redis"] != "ok" {
			response.Error(ctx, http.StatusServiceUnavailable, "DEPENDENCY_UNAVAILABLE", "One or more dependencies are unavailable")
			return
		}

		response.OK(ctx, status, "dependencies ok")
	}
}
