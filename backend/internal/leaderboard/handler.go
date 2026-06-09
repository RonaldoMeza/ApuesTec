package leaderboard

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetGlobalLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	leaderboard, err := h.service.GetGlobalLeaderboard(c.Request.Context(), limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "LEADERBOARD_ERROR", "failed to retrieve leaderboard")
		return
	}

	response.OK(c, leaderboard, "global leaderboard retrieved successfully")
}
