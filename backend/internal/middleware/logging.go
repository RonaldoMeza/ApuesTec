package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		requestID, _ := ctx.Get(RequestIDKey)
		log.Printf(
			"request_id=%v method=%s path=%s status=%d latency=%s client_ip=%s",
			requestID,
			ctx.Request.Method,
			ctx.Request.URL.Path,
			ctx.Writer.Status(),
			time.Since(start),
			ctx.ClientIP(),
		)
	}
}
