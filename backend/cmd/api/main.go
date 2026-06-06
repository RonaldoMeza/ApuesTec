package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/api/v1/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "apuestec-backend",
		})
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
