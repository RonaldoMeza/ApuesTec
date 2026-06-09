package roominvites

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/rooms"
)

type RoomRepository interface {
	FindByID(ctx context.Context, id string) (*rooms.Room, error)
}

type MemberRepository interface {
	GetMemberRole(ctx context.Context, roomID, userID string) (string, error)
	AddMember(ctx context.Context, roomID, userID, role string) error
}

type Service interface {
	CreateInvite(ctx context.Context, roomID, userID string, durationMinutes int) (*InviteResponse, error)
	PreviewInvite(ctx context.Context, code string) (*InvitePreviewResponse, error)
	JoinRoom(ctx context.Context, code, userID string) (*JoinResponse, error)
}

type service struct {
	repo       Repository
	roomRepo   RoomRepository
	memberRepo MemberRepository
}

func NewService(repo Repository, roomRepo RoomRepository, memberRepo MemberRepository) Service {
	return &service{
		repo:       repo,
		roomRepo:   roomRepo,
		memberRepo: memberRepo,
	}
}

func validateDuration(minutes int) error {
	for _, d := range AllowedDurations {
		if d == minutes {
			return nil
		}
	}
	return &ServiceError{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_DURATION",
		Message: fmt.Sprintf("duration must be one of %v", AllowedDurations),
	}
}

func (s *service) CreateInvite(ctx context.Context, roomID, userID string, durationMinutes int) (*InviteResponse, error) {
	if err := validateDuration(durationMinutes); err != nil {
		return nil, err
	}

	role, err := s.memberRepo.GetMemberRole(ctx, roomID, userID)
	if err != nil || role == "" {
		return nil, &ServiceError{
			Status:  http.StatusForbidden,
			Code:    "NOT_MEMBER",
			Message: "you are not a member of this room",
		}
	}

	if role != "OWNER" && role != "MODERATOR" {
		return nil, &ServiceError{
			Status:  http.StatusForbidden,
			Code:    "INSUFFICIENT_PERMISSIONS",
			Message: "only owners and moderators can create invites",
		}
	}

	room, err := s.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return nil, &ServiceError{
			Status:  http.StatusNotFound,
			Code:    "ROOM_NOT_FOUND",
			Message: "room not found",
		}
	}

	if room.Status != "ACTIVE" {
		return nil, &ServiceError{
			Status:  http.StatusBadRequest,
			Code:    "ROOM_NOT_ACTIVE",
			Message: "room is not active",
		}
	}

	if err := s.repo.RevokeActiveInvites(ctx, roomID); err != nil {
		return nil, fmt.Errorf("revoke active invites: %w", err)
	}

	code, err := generateCode()
	if err != nil {
		return nil, fmt.Errorf("generate code: %w", err)
	}

	qrPayload := "https://apuestec.app/join/" + code
	expiresAt := time.Now().Add(time.Duration(durationMinutes) * time.Minute)

	invite, err := s.repo.Create(ctx, roomID, code, qrPayload, userID, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("create invite: %w", err)
	}

	return toInviteResponse(invite), nil
}

func (s *service) PreviewInvite(ctx context.Context, code string) (*InvitePreviewResponse, error) {
	invite, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, &ServiceError{
			Status:  http.StatusNotFound,
			Code:    "INVITE_NOT_FOUND",
			Message: "invite not found",
		}
	}

	room, err := s.roomRepo.FindByID(ctx, invite.RoomID)
	if err != nil {
		return nil, &ServiceError{
			Status:  http.StatusNotFound,
			Code:    "ROOM_NOT_FOUND",
			Message: "room not found",
		}
	}

	isExpired := time.Now().After(invite.ExpiresAt)
	desc := ""
	if room.Description != nil {
		desc = *room.Description
	}

	return &InvitePreviewResponse{
		Code:            invite.Code,
		RoomName:        room.Name,
		RoomDescription: desc,
		ExpiresAt:       invite.ExpiresAt.Format(time.RFC3339),
		IsExpired:       isExpired,
	}, nil
}

func (s *service) JoinRoom(ctx context.Context, code, userID string) (*JoinResponse, error) {
	invite, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, &ServiceError{
			Status:  http.StatusNotFound,
			Code:    "INVITE_NOT_FOUND",
			Message: "invite not found",
		}
	}

	if time.Now().After(invite.ExpiresAt) {
		return nil, &ServiceError{
			Status:  http.StatusGone,
			Code:    "INVITE_EXPIRED",
			Message: "this invite has expired",
		}
	}

	if invite.RevokedAt != nil {
		return nil, &ServiceError{
			Status:  http.StatusGone,
			Code:    "INVITE_REVOKED",
			Message: "this invite has been revoked",
		}
	}

	if invite.UsedAt != nil {
		return nil, &ServiceError{
			Status:  http.StatusGone,
			Code:    "INVITE_USED",
			Message: "this invite has already been used",
		}
	}

	room, err := s.roomRepo.FindByID(ctx, invite.RoomID)
	if err != nil {
		return nil, &ServiceError{
			Status:  http.StatusNotFound,
			Code:    "ROOM_NOT_FOUND",
			Message: "room not found",
		}
	}

	if room.Status != "ACTIVE" {
		return nil, &ServiceError{
			Status:  http.StatusBadRequest,
			Code:    "ROOM_NOT_ACTIVE",
			Message: "room is not active",
		}
	}

	existingRole, err := s.memberRepo.GetMemberRole(ctx, invite.RoomID, userID)
	if err == nil && existingRole != "" {
		return nil, &ServiceError{
			Status:  http.StatusConflict,
			Code:    "ALREADY_MEMBER",
			Message: "you are already a member of this room",
		}
	}

	if err := s.repo.MarkUsed(ctx, invite.ID); err != nil {
		return nil, fmt.Errorf("mark invite used: %w", err)
	}

	if err := s.memberRepo.AddMember(ctx, invite.RoomID, userID, "MEMBER"); err != nil {
		return nil, fmt.Errorf("add member: %w", err)
	}

	return &JoinResponse{
		RoomID:   room.ID,
		RoomName: room.Name,
		Role:     "MEMBER",
	}, nil
}

func toInviteResponse(invite *RoomInvite) *InviteResponse {
	resp := &InviteResponse{
		ID:        invite.ID,
		RoomID:    invite.RoomID,
		Code:      invite.Code,
		CreatedBy: invite.CreatedBy,
		ExpiresAt: invite.ExpiresAt.Format(time.RFC3339),
		CreatedAt: invite.CreatedAt.Format(time.RFC3339),
	}

	if invite.QRPayload != nil {
		resp.QRPayload = *invite.QRPayload
	}

	if invite.UsedAt != nil {
		s := invite.UsedAt.Format(time.RFC3339)
		resp.UsedAt = &s
	}

	if invite.RevokedAt != nil {
		s := invite.RevokedAt.Format(time.RFC3339)
		resp.RevokedAt = &s
	}

	return resp
}
