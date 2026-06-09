package rooms

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type MemberRepository interface {
	AddMember(ctx context.Context, roomID, userID, role string) error
	GetMemberRole(ctx context.Context, roomID, userID string) (string, error)
}

type Service interface {
	Create(ctx context.Context, userID, name string, description *string, visibility *string, password *string, clientIP string) (*RoomResponse, error)
	ListMyRooms(ctx context.Context, userID string) (*RoomListResponse, error)
	GetByID(ctx context.Context, roomID, userID string) (*RoomResponse, error)
	Update(ctx context.Context, roomID, userID, name string, description *string, visibility *string, password *string) (*RoomResponse, error)
	Close(ctx context.Context, roomID, userID string) (*RoomResponse, error)
	GetLeaderboard(ctx context.Context, roomID, userID string) (*RoomLeaderboardResponse, error)
	SearchPublicRooms(ctx context.Context, networkPrefix, query, userID string) (*SearchPublicRoomsResponse, error)
	JoinPublicRoom(ctx context.Context, roomID, userID, password string) (*RoomResponse, error)
}

type service struct {
	repo       Repository
	memberRepo MemberRepository
}

func NewService(repo Repository, memberRepo MemberRepository) Service {
	return &service{repo: repo, memberRepo: memberRepo}
}

func computeNetworkPrefix(clientIP string) *string {
	idx := strings.LastIndex(clientIP, ".")
	if idx == -1 {
		idx = strings.LastIndex(clientIP, ":")
	}
	if idx == -1 {
		return nil
	}
	prefix := clientIP[:idx]
	return &prefix
}

func (s *service) enrichRoom(ctx context.Context, room *Room, userID string) (*RoomResponse, error) {
	memberCount, err := s.repo.CountMembers(ctx, room.ID)
	if err != nil {
		return nil, fmt.Errorf("count members: %w", err)
	}

	myRole, err := s.memberRepo.GetMemberRole(ctx, room.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("get member role: %w", err)
	}

	desc := ""
	if room.Description != nil {
		desc = *room.Description
	}

	closedAt := ""
	if room.ClosedAt != nil {
		closedAt = room.ClosedAt.Format(timeFormat)
	}

	hasPassword := room.PasswordHash != nil && *room.PasswordHash != ""

	return &RoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Description: desc,
		OwnerID:     room.OwnerID,
		Status:      room.Status,
		Visibility:  room.Visibility,
		HasPassword: hasPassword,
		MemberCount: memberCount,
		MyRole:      myRole,
		CreatedAt:   room.CreatedAt.Format(timeFormat),
		UpdatedAt:   room.UpdatedAt.Format(timeFormat),
		ClosedAt:    closedAt,
	}, nil
}

func (s *service) Create(ctx context.Context, userID, name string, description *string, visibility *string, password *string, clientIP string) (*RoomResponse, error) {
	v := VisibilityPrivate
	if visibility != nil && *visibility == VisibilityPublic {
		v = VisibilityPublic
	}

	var passwordHash *string
	if v == VisibilityPublic && password != nil && *password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		h := string(hash)
		passwordHash = &h
	}

	networkPrefix := computeNetworkPrefix(clientIP)

	room, err := s.repo.Create(ctx, userID, name, description, v, passwordHash, networkPrefix)
	if err != nil {
		return nil, fmt.Errorf("create room: %w", err)
	}

	if err := s.memberRepo.AddMember(ctx, room.ID, userID, RoomRoleOwner); err != nil {
		return nil, fmt.Errorf("add owner member: %w", err)
	}

	return s.enrichRoom(ctx, room, userID)
}

func (s *service) ListMyRooms(ctx context.Context, userID string) (*RoomListResponse, error) {
	rooms, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list my rooms: %w", err)
	}

	responses := make([]RoomResponse, 0, len(rooms))
	for _, r := range rooms {
		enriched, err := s.enrichRoom(ctx, &r, userID)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *enriched)
	}

	return &RoomListResponse{
		Rooms: responses,
		Total: len(responses),
	}, nil
}

func (s *service) GetByID(ctx context.Context, roomID, userID string) (*RoomResponse, error) {
	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "ROOM_NOT_FOUND", Message: "room not found"}
	}

	myRole, err := s.memberRepo.GetMemberRole(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("get member role: %w", err)
	}

	if myRole == "" {
		return nil, &ServiceError{Status: http.StatusForbidden, Code: "NOT_A_MEMBER", Message: "you are not a member of this room"}
	}

	return s.enrichRoom(ctx, room, userID)
}

func (s *service) Update(ctx context.Context, roomID, userID, name string, description *string, visibility *string, password *string) (*RoomResponse, error) {
	role, err := s.memberRepo.GetMemberRole(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("get member role: %w", err)
	}

	if role != RoomRoleOwner {
		return nil, &ServiceError{Status: http.StatusForbidden, Code: "NOT_OWNER", Message: "only the owner can update the room"}
	}

	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "ROOM_NOT_FOUND", Message: "room not found"}
	}

	if room.Status != RoomStatusActive {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "ROOM_CLOSED", Message: "cannot update a closed room"}
	}

	v := room.Visibility
	if visibility != nil {
		v = *visibility
	}

	var passwordHash *string
	if password != nil {
		if *password == "" {
			empty := ""
			passwordHash = &empty
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("hash password: %w", err)
			}
			h := string(hash)
			passwordHash = &h
		}
	}

	if v == VisibilityPublic && (passwordHash == nil || *passwordHash == "") && (room.PasswordHash == nil || *room.PasswordHash == "") {
		return nil, &ServiceError{
			Status:  http.StatusBadRequest,
			Code:    "PASSWORD_REQUIRED",
			Message: "a public room must have a password",
		}
	}

	if v != room.Visibility {
		if err := s.repo.UpdateVisibility(ctx, roomID, v); err != nil {
			return nil, fmt.Errorf("update visibility: %w", err)
		}
	}

	if passwordHash != nil {
		if err := s.repo.SetPasswordHash(ctx, roomID, passwordHash); err != nil {
			return nil, fmt.Errorf("set password: %w", err)
		}
	}

	updated, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "ROOM_NOT_FOUND", Message: "room not found"}
	}

	return s.enrichRoom(ctx, updated, userID)
}

func (s *service) Close(ctx context.Context, roomID, userID string) (*RoomResponse, error) {
	role, err := s.memberRepo.GetMemberRole(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("get member role: %w", err)
	}

	if role != RoomRoleOwner {
		return nil, &ServiceError{Status: http.StatusForbidden, Code: "NOT_OWNER", Message: "only the owner can close the room"}
	}

	closed, err := s.repo.Close(ctx, roomID)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "ROOM_NOT_FOUND", Message: "room not found"}
	}

	return s.enrichRoom(ctx, closed, userID)
}

func (s *service) GetLeaderboard(ctx context.Context, roomID, userID string) (*RoomLeaderboardResponse, error) {
	myRole, err := s.memberRepo.GetMemberRole(ctx, roomID, userID)
	if err != nil || myRole == "" {
		return nil, &ServiceError{Status: http.StatusForbidden, Code: "NOT_A_MEMBER", Message: "you are not a member of this room"}
	}

	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "ROOM_NOT_FOUND", Message: "room not found"}
	}

	entries, err := s.repo.GetLeaderboard(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("get leaderboard: %w", err)
	}

	if entries == nil {
		entries = []RoomLeaderboardEntry{}
	}

	return &RoomLeaderboardResponse{
		Entries:  entries,
		Total:    len(entries),
		RoomID:   room.ID,
		RoomName: room.Name,
	}, nil
}

func (s *service) SearchPublicRooms(ctx context.Context, networkPrefix, query, userID string) (*SearchPublicRoomsResponse, error) {
	rooms, err := s.repo.SearchPublicByNetwork(ctx, networkPrefix, query, userID)
	if err != nil {
		return nil, fmt.Errorf("search public rooms: %w", err)
	}

	if rooms == nil {
		rooms = []PublicRoomResponse{}
	}

	return &SearchPublicRoomsResponse{
		Rooms: rooms,
		Total: len(rooms),
	}, nil
}

func (s *service) JoinPublicRoom(ctx context.Context, roomID, userID, password string) (*RoomResponse, error) {
	existingRole, err := s.memberRepo.GetMemberRole(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("get member role: %w", err)
	}
	if existingRole != "" {
		return nil, &ServiceError{
			Status:  http.StatusConflict,
			Code:    "ALREADY_MEMBER",
			Message: "you are already a member of this room",
		}
	}

	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return nil, &ServiceError{
			Status:  http.StatusNotFound,
			Code:    "ROOM_NOT_FOUND",
			Message: "room not found",
		}
	}

	if room.Status != RoomStatusActive {
		return nil, &ServiceError{
			Status:  http.StatusBadRequest,
			Code:    "ROOM_NOT_ACTIVE",
			Message: "room is not active",
		}
	}

	if room.Visibility != VisibilityPublic {
		return nil, &ServiceError{
			Status:  http.StatusBadRequest,
			Code:    "NOT_PUBLIC",
			Message: "room is not public",
		}
	}

	if room.PasswordHash != nil && *room.PasswordHash != "" {
		if password == "" {
			return nil, &ServiceError{
				Status:  http.StatusBadRequest,
				Code:    "PASSWORD_REQUIRED",
				Message: "this room requires a password",
			}
		}
		if err := bcrypt.CompareHashAndPassword([]byte(*room.PasswordHash), []byte(password)); err != nil {
			return nil, &ServiceError{
				Status:  http.StatusUnauthorized,
				Code:    "WRONG_PASSWORD",
				Message: "incorrect room password",
			}
		}
	}

	if err := s.memberRepo.AddMember(ctx, roomID, userID, RoomRoleMember); err != nil {
		return nil, fmt.Errorf("add member: %w", err)
	}

	return s.enrichRoom(ctx, room, userID)
}
