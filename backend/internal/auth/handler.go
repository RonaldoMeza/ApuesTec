package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.service.Register(c.Request.Context(), req, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp, "registration successful")
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
			return
		}
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp, "login successful")
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.service.Refresh(c.Request.Context(), req.RefreshToken, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp, "token refreshed")
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, nil, "logout successful")
}

func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	resp, err := h.service.Me(c.Request.Context(), userID.(string))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp, "user info retrieved")
}

func (h *Handler) GoogleAuth(c *gin.Context) {
	var req GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	resp, err := h.service.GoogleAuth(c.Request.Context(), req, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp, "google authentication successful")
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID.(string), req.CurrentPassword, req.NewPassword); err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, nil, "password changed successfully")
}

func handleServiceError(c *gin.Context, err error) {
	var svcErr *ServiceError
	if errors.As(err, &svcErr) {
		response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
		return
	}
	response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
}
