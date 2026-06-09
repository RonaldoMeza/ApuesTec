package roommembers

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

func (h *Handler) ListMembers(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")

	members, err := h.service.ListMembers(c.Request.Context(), roomID, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, members, "members retrieved successfully")
}

func (h *Handler) ChangeRole(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")
	targetUserID := c.Param("userId")

	var req ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.service.ChangeRole(c.Request.Context(), roomID, userID, targetUserID, req.Role); err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, nil, "role changed successfully")
}

func (h *Handler) RemoveMember(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")
	targetUserID := c.Param("userId")

	if err := h.service.RemoveMember(c.Request.Context(), roomID, userID, targetUserID); err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, nil, "member removed successfully")
}

func (h *Handler) Leave(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")

	if err := h.service.Leave(c.Request.Context(), roomID, userID); err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, nil, "you left the room successfully")
}

func handleServiceError(c *gin.Context, err error) {
	var svcErr *ServiceError
	if errors.As(err, &svcErr) {
		response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
		return
	}
	response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
}
