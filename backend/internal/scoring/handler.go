package scoring

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RebuildAllScores(c *gin.Context) {
	if err := h.service.RebuildAllUserScores(c.Request.Context()); err != nil {
		response.Error(c, http.StatusInternalServerError, "REBUILD_ERROR", "failed to rebuild all user scores")
		return
	}

	response.OK(c, gin.H{"message": "all user scores rebuilt successfully"}, "all user scores rebuilt successfully")
}
