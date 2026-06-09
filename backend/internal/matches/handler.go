package matches

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

func (h *Handler) List(c *gin.Context) {
	matches, err := h.service.List(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, matches, "matches retrieved successfully")
}

func (h *Handler) ListUpcoming(c *gin.Context) {
	matches, err := h.service.ListUpcoming(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, matches, "upcoming matches retrieved successfully")
}

func (h *Handler) ListFinished(c *gin.Context) {
	matches, err := h.service.ListFinished(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, matches, "finished matches retrieved successfully")
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	match, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, match, "match retrieved successfully")
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	match, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, match, "match created successfully")
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	match, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, match, "match updated successfully")
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	match, err := h.service.UpdateStatus(c.Request.Context(), id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, match, "match status updated successfully")
}

func (h *Handler) UpdateResult(c *gin.Context) {
	id := c.Param("id")
	var req UpdateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	match, err := h.service.UpdateResult(c.Request.Context(), id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, match, "match result updated successfully")
}

func handleServiceError(c *gin.Context, err error) {
	var svcErr *ServiceError
	if errors.As(err, &svcErr) {
		response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
		return
	}
	response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
}
