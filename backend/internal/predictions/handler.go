package predictions

import (
	"errors"
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

func (h *Handler) GetMyPrediction(c *gin.Context) {
	userID := c.GetString("userID")
	matchID := c.Param("id")

	prediction, err := h.service.GetMyPrediction(c.Request.Context(), userID, matchID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, prediction, "prediction retrieved successfully")
}

func (h *Handler) ListMyPredictions(c *gin.Context) {
	userID := c.GetString("userID")

	predictions, err := h.service.ListMyPredictions(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	if predictions == nil {
		predictions = []PredictionResponse{}
	}

	response.OK(c, predictions, "predictions retrieved successfully")
}

func (h *Handler) UpsertPrediction(c *gin.Context) {
	userID := c.GetString("userID")
	matchID := c.Param("id")

	var req CreatePredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	prediction, err := h.service.UpsertPrediction(c.Request.Context(), userID, matchID, req, &ip, &userAgent)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, prediction, "prediction saved successfully")
}

func (h *Handler) DeletePrediction(c *gin.Context) {
	userID := c.GetString("userID")
	matchID := c.Param("id")

	if err := h.service.DeletePrediction(c.Request.Context(), userID, matchID); err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, nil, "prediction deleted successfully")
}

func handleServiceError(c *gin.Context, err error) {
	var svcErr *ServiceError
	if errors.As(err, &svcErr) {
		response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
		return
	}
	response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
}
