package roominvites

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

func (h *Handler) CreateInvite(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "room id is required")
		return
	}

	var req CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.service.CreateInvite(c.Request.Context(), roomID, userID, req.DurationMinutes)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.OK(c, resp, "invite created successfully")
}

func (h *Handler) PreviewInvite(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "code is required")
		return
	}

	resp, err := h.service.PreviewInvite(c.Request.Context(), code)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.OK(c, resp, "invite preview retrieved")
}

func (h *Handler) JoinRoom(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	code := c.Param("code")
	if code == "" {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "code is required")
		return
	}

	resp, err := h.service.JoinRoom(c.Request.Context(), code, userID)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.OK(c, resp, "joined room successfully")
}
