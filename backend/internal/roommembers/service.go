package roommembers

import (
	"context"
	"fmt"
	"net/http"
)

var (
	ErrNotMember        = fmt.Errorf("user is not a member of this room")
	ErrOwnerCannotLeave = fmt.Errorf("owner cannot leave the room")
	ErrCannotChangeOwner = fmt.Errorf("cannot change owner role")
	ErrInsufficientPermissions = fmt.Errorf("insufficient permissions")
	ErrTargetNotMember  = fmt.Errorf("target user is not a member")
)

type Service interface {
	ListMembers(ctx context.Context, roomID, userID string) (*MemberListResponse, error)
	ChangeRole(ctx context.Context, roomID, actorID, targetUserID, newRole string) error
	RemoveMember(ctx context.Context, roomID, actorID, targetUserID string) error
	Leave(ctx context.Context, roomID, userID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListMembers(ctx context.Context, roomID, userID string) (*MemberListResponse, error) {
	role, err := s.repo.GetMemberRole(ctx, roomID, userID)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusForbidden, Code: "NOT_MEMBER", Message: ErrNotMember.Error()}
	}
	if role == "" {
		return nil, &ServiceError{Status: http.StatusForbidden, Code: "NOT_MEMBER", Message: ErrNotMember.Error()}
	}

	members, err := s.repo.ListMembers(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("list members service: %w", err)
	}

	if members == nil {
		members = []MemberResponse{}
	}

	return &MemberListResponse{
		Members: members,
		Total:   len(members),
	}, nil
}

func (s *service) ChangeRole(ctx context.Context, roomID, actorID, targetUserID, newRole string) error {
	actorRole, err := s.repo.GetMemberRole(ctx, roomID, actorID)
	if err != nil || actorRole == "" {
		return &ServiceError{Status: http.StatusForbidden, Code: "NOT_MEMBER", Message: ErrNotMember.Error()}
	}

	if actorRole != RoomRoleOwner {
		return &ServiceError{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: ErrInsufficientPermissions.Error()}
	}

	targetRole, err := s.repo.GetMemberRole(ctx, roomID, targetUserID)
	if err != nil || targetRole == "" {
		return &ServiceError{Status: http.StatusNotFound, Code: "TARGET_NOT_MEMBER", Message: ErrTargetNotMember.Error()}
	}

	if targetRole == RoomRoleOwner {
		return &ServiceError{Status: http.StatusForbidden, Code: "CANNOT_CHANGE_OWNER", Message: ErrCannotChangeOwner.Error()}
	}

	if err := s.repo.UpdateRole(ctx, roomID, targetUserID, newRole); err != nil {
		return fmt.Errorf("change role: %w", err)
	}

	return nil
}

func (s *service) RemoveMember(ctx context.Context, roomID, actorID, targetUserID string) error {
	actorRole, err := s.repo.GetMemberRole(ctx, roomID, actorID)
	if err != nil || actorRole == "" {
		return &ServiceError{Status: http.StatusForbidden, Code: "NOT_MEMBER", Message: ErrNotMember.Error()}
	}

	targetRole, err := s.repo.GetMemberRole(ctx, roomID, targetUserID)
	if err != nil || targetRole == "" {
		return &ServiceError{Status: http.StatusNotFound, Code: "TARGET_NOT_MEMBER", Message: ErrTargetNotMember.Error()}
	}

	if actorID == targetUserID {
		return &ServiceError{Status: http.StatusForbidden, Code: "CANNOT_REMOVE_SELF", Message: "cannot remove yourself, use leave instead"}
	}

	switch actorRole {
	case RoomRoleOwner:
		if targetRole == RoomRoleOwner {
			return &ServiceError{Status: http.StatusForbidden, Code: "CANNOT_REMOVE_OWNER", Message: ErrCannotChangeOwner.Error()}
		}
	case RoomRoleModerator:
		if targetRole != RoomRoleMember {
			return &ServiceError{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: ErrInsufficientPermissions.Error()}
		}
	case RoomRoleMember:
		return &ServiceError{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: ErrInsufficientPermissions.Error()}
	}

	if err := s.repo.RemoveMember(ctx, roomID, targetUserID); err != nil {
		return fmt.Errorf("remove member: %w", err)
	}

	return nil
}

func (s *service) Leave(ctx context.Context, roomID, userID string) error {
	role, err := s.repo.GetMemberRole(ctx, roomID, userID)
	if err != nil || role == "" {
		return &ServiceError{Status: http.StatusForbidden, Code: "NOT_MEMBER", Message: ErrNotMember.Error()}
	}

	if role == RoomRoleOwner {
		return &ServiceError{Status: http.StatusForbidden, Code: "OWNER_CANNOT_LEAVE", Message: ErrOwnerCannotLeave.Error()}
	}

	if err := s.repo.RemoveMember(ctx, roomID, userID); err != nil {
		return fmt.Errorf("leave room: %w", err)
	}

	return nil
}
