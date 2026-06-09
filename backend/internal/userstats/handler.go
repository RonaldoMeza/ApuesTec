package userstats

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

func (h *Handler) GetMyStats(c *gin.Context) {
	userID := c.GetString("userID")

	stats, err := h.service.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "STATS_ERROR", "failed to retrieve user stats")
		return
	}

	response.OK(c, stats, "user stats retrieved successfully")
}

func (h *Handler) ListMyScoreEvents(c *gin.Context) {
	userID := c.GetString("userID")

	events, err := h.service.ListScoreEvents(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SCORE_EVENTS_ERROR", "failed to retrieve score events")
		return
	}

	if events == nil {
		events = []ScoreEventResponse{}
	}

	response.OK(c, events, "score events retrieved successfully")
}
