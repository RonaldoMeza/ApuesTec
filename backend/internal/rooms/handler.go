package rooms

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

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("userID")
	clientIP := c.ClientIP()
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	result, err := h.service.Create(c.Request.Context(), userID, req.Name, req.Description, req.Visibility, req.Password, clientIP)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "ROOMS_ERROR", err.Error())
		return
	}
	response.OK(c, result, "room created successfully")
}

func (h *Handler) ListMyRooms(c *gin.Context) {
	userID := c.GetString("userID")

	result, err := h.service.ListMyRooms(c.Request.Context(), userID)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "ROOMS_ERROR", err.Error())
		return
	}
	response.OK(c, result, "rooms retrieved successfully")
}

func (h *Handler) SearchPublic(c *gin.Context) {
	userID := c.GetString("userID")
	clientIP := c.ClientIP()
	networkPrefix := computeNetworkPrefix(clientIP)
	if networkPrefix == nil {
		response.Error(c, http.StatusBadRequest, "INVALID_IP", "could not determine network")
		return
	}

	q := c.DefaultQuery("q", "")

	result, err := h.service.SearchPublicRooms(c.Request.Context(), *networkPrefix, q, userID)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "ROOMS_ERROR", err.Error())
		return
	}
	response.OK(c, result, "public rooms retrieved successfully")
}

func (h *Handler) GetByID(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")

	result, err := h.service.GetByID(c.Request.Context(), roomID, userID)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "ROOMS_ERROR", err.Error())
		return
	}
	response.OK(c, result, "room retrieved successfully")
}

func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")
	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	result, err := h.service.Update(c.Request.Context(), roomID, userID, req.Name, req.Description, req.Visibility, req.Password)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "ROOMS_ERROR", err.Error())
		return
	}
	response.OK(c, result, "room updated successfully")
}

func (h *Handler) Close(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")

	result, err := h.service.Close(c.Request.Context(), roomID, userID)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "ROOMS_ERROR", err.Error())
		return
	}
	response.OK(c, result, "room closed successfully")
}

func (h *Handler) JoinPublic(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")
	var req JoinPublicRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	result, err := h.service.JoinPublicRoom(c.Request.Context(), roomID, userID, req.Password)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "ROOMS_ERROR", err.Error())
		return
	}
	response.OK(c, result, "joined room successfully")
}

func (h *Handler) GetLeaderboard(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")

	result, err := h.service.GetLeaderboard(c.Request.Context(), roomID, userID)
	if err != nil {
		var svcErr *ServiceError
		if errors.As(err, &svcErr) {
			response.Error(c, svcErr.Status, svcErr.Code, svcErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, "LEADERBOARD_ERROR", err.Error())
		return
	}
	response.OK(c, result, "room leaderboard retrieved successfully")
}
