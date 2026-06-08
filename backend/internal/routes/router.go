package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/audit"
	appauth "github.com/RonaldoMeza/ApuesTec/backend/internal/auth"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/config"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/database"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/middleware"
	appredis "github.com/RonaldoMeza/ApuesTec/backend/internal/redis"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/response"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/roles"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/users"
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

	userRepo := users.NewRepository(db)
	authRepo := appauth.NewRepository(db)
	roleRepo := roles.NewRepository(db)
	auditRepo := audit.NewRepository(db)

	authSvc := appauth.NewService(userRepo, authRepo, roleRepo, auditRepo, cfg)
	authHandler := appauth.NewHandler(authSvc)

	api := router.Group(cfg.APIPrefix)
	{
		api.GET("/health", health(cfg))
		api.GET("/health/dependencies", dependenciesHealth(db, redisClient))

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/google", authHandler.GoogleAuth)

			auth.GET("/me", appauth.AuthMiddleware(cfg.JWTAccessSecret), authHandler.Me)
			auth.POST("/change-password", appauth.AuthMiddleware(cfg.JWTAccessSecret), authHandler.ChangePassword)
		}
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
