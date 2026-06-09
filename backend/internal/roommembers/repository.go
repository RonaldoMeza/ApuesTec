package roommembers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	ListMembers(ctx context.Context, roomID string) ([]MemberResponse, error)
	GetMemberRole(ctx context.Context, roomID, userID string) (string, error)
	GetMemberByID(ctx context.Context, memberID string) (*struct {
		ID, RoomID, UserID, Role string
	}, error)
	UpdateRole(ctx context.Context, roomID, userID, role string) error
	RemoveMember(ctx context.Context, roomID, userID string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) ListMembers(ctx context.Context, roomID string) ([]MemberResponse, error) {
	query := `
		SELECT rm.id, rm.room_id, rm.user_id, u.full_name, u.email, rm.role,
		       TO_CHAR(rm.joined_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as joined_at
		FROM room_members rm
		JOIN users u ON u.id = rm.user_id
		WHERE rm.room_id = $1
		ORDER BY rm.joined_at ASC
	`
	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()

	var members []MemberResponse
	for rows.Next() {
		var m MemberResponse
		if err := rows.Scan(&m.ID, &m.RoomID, &m.UserID, &m.FullName, &m.Email, &m.Role, &m.JoinedAt); err != nil {
			return nil, fmt.Errorf("scan member row: %w", err)
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *repository) GetMemberRole(ctx context.Context, roomID, userID string) (string, error) {
	query := `SELECT role FROM room_members WHERE room_id = $1 AND user_id = $2`
	var role string
	err := r.db.QueryRow(ctx, query, roomID, userID).Scan(&role)
	if err != nil {
		return "", fmt.Errorf("get member role: %w", err)
	}
	return role, nil
}

func (r *repository) GetMemberByID(ctx context.Context, memberID string) (*struct {
	ID, RoomID, UserID, Role string
}, error) {
	query := `SELECT id, room_id, user_id, role FROM room_members WHERE id = $1`
	m := &struct {
		ID, RoomID, UserID, Role string
	}{}
	err := r.db.QueryRow(ctx, query, memberID).Scan(&m.ID, &m.RoomID, &m.UserID, &m.Role)
	if err != nil {
		return nil, fmt.Errorf("get member by id: %w", err)
	}
	return m, nil
}

func (r *repository) UpdateRole(ctx context.Context, roomID, userID, role string) error {
	query := `UPDATE room_members SET role = $1 WHERE room_id = $2 AND user_id = $3`
	_, err := r.db.Exec(ctx, query, role, roomID, userID)
	if err != nil {
		return fmt.Errorf("update member role: %w", err)
	}
	return nil
}

func (r *repository) RemoveMember(ctx context.Context, roomID, userID string) error {
	query := `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, roomID, userID)
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}
