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
	"github.com/RonaldoMeza/ApuesTec/backend/internal/matches"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/middleware"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/predictions"
	appredis "github.com/RonaldoMeza/ApuesTec/backend/internal/redis"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/response"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/roles"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/teams"
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

	teamRepo := teams.NewRepository(db)
	teamSvc := teams.NewService(teamRepo)
	teamHandler := teams.NewHandler(teamSvc)

	teamInfoRepo := matches.NewTeamInfoRepository(db)
	matchRepo := matches.NewRepository(db)
	matchSvc := matches.NewService(matchRepo, teamInfoRepo)
	matchHandler := matches.NewHandler(matchSvc)

	predictionRepo := predictions.NewRepository(db)
	predictionSvc := predictions.NewService(predictionRepo, matchRepo, auditRepo)
	predictionHandler := predictions.NewHandler(predictionSvc)

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

		adminAuth := appauth.AuthMiddleware(cfg.JWTAccessSecret)
		adminRole := appauth.RequireRole("ADMIN", "SUPER_ADMIN")

		api.GET("/teams", teamHandler.List)
		api.GET("/teams/:id", teamHandler.GetByID)

		api.GET("/matches/upcoming", matchHandler.ListUpcoming)
		api.GET("/matches/finished", matchHandler.ListFinished)
		api.GET("/matches", matchHandler.List)
		api.GET("/matches/:id", matchHandler.GetByID)

		userAuth := appauth.AuthMiddleware(cfg.JWTAccessSecret)
		api.GET("/predictions/me", userAuth, predictionHandler.ListMyPredictions)
		api.GET("/matches/:id/prediction", userAuth, predictionHandler.GetMyPrediction)
		api.POST("/matches/:id/prediction", userAuth, predictionHandler.UpsertPrediction)
		api.DELETE("/matches/:id/prediction", userAuth, predictionHandler.DeletePrediction)

		admin := api.Group("/admin", adminAuth, adminRole)
		{
			admin.POST("/teams", teamHandler.Create)
			admin.PUT("/teams/:id", teamHandler.Update)
			admin.DELETE("/teams/:id", teamHandler.Delete)

			admin.POST("/matches", matchHandler.Create)
			admin.PUT("/matches/:id", matchHandler.Update)
			admin.PATCH("/matches/:id/status", matchHandler.UpdateStatus)
			admin.PATCH("/matches/:id/result", matchHandler.UpdateResult)
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
