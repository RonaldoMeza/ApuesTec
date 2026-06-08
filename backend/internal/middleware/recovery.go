package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/response"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
		requestID, _ := ctx.Get(RequestIDKey)
		log.Printf("panic recovered request_id=%v error=%v", requestID, recovered)
		response.Error(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "Unexpected error")
	})
}
