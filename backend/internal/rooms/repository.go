package rooms

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, ownerID, name string, description *string, visibility string, passwordHash *string, networkPrefix *string) (*Room, error)
	FindByID(ctx context.Context, id string) (*Room, error)
	FindByUserID(ctx context.Context, userID string) ([]Room, error)
	SearchPublicByNetwork(ctx context.Context, networkPrefix, query string, userID string) ([]PublicRoomResponse, error)
	Update(ctx context.Context, id, name string, description *string, visibility *string, passwordHash *string) (*Room, error)
	UpdateVisibility(ctx context.Context, id, visibility string) error
	SetPasswordHash(ctx context.Context, id string, passwordHash *string) error
	Close(ctx context.Context, id string) (*Room, error)
	CountMembers(ctx context.Context, roomID string) (int, error)
	GetMemberRole(ctx context.Context, roomID, userID string) (string, error)
	AddMember(ctx context.Context, roomID, userID, role string) error
	GetLeaderboard(ctx context.Context, roomID string) ([]RoomLeaderboardEntry, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func scanRoom(row pgx.Row) (*Room, error) {
	r := &Room{}
	err := row.Scan(&r.ID, &r.Name, &r.Description, &r.OwnerID, &r.Status, &r.Visibility, &r.PasswordHash, &r.NetworkPrefix, &r.CreatedAt, &r.UpdatedAt, &r.ClosedAt)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *repository) Create(ctx context.Context, ownerID, name string, description *string, visibility string, passwordHash *string, networkPrefix *string) (*Room, error) {
	query := `
		INSERT INTO rooms (name, description, owner_id, visibility, password_hash, network_prefix)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, owner_id, status, visibility, password_hash, network_prefix, created_at, updated_at, closed_at
	`
	room, err := scanRoom(r.db.QueryRow(ctx, query, name, description, ownerID, visibility, passwordHash, networkPrefix))
	if err != nil {
		return nil, fmt.Errorf("create room: %w", err)
	}
	return room, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Room, error) {
	query := `SELECT id, name, description, owner_id, status, visibility, password_hash, network_prefix, created_at, updated_at, closed_at FROM rooms WHERE id = $1`
	room, err := scanRoom(r.db.QueryRow(ctx, query, id))
	if err != nil {
		return nil, fmt.Errorf("find room by id: %w", err)
	}
	return room, nil
}

func (r *repository) FindByUserID(ctx context.Context, userID string) ([]Room, error) {
	query := `
		SELECT r.id, r.name, r.description, r.owner_id, r.status, r.visibility, r.password_hash, r.network_prefix, r.created_at, r.updated_at, r.closed_at
		FROM rooms r
		JOIN room_members rm ON rm.room_id = r.id
		WHERE rm.user_id = $1
		ORDER BY r.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find rooms by user id: %w", err)
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.OwnerID, &room.Status, &room.Visibility, &room.PasswordHash, &room.NetworkPrefix, &room.CreatedAt, &room.UpdatedAt, &room.ClosedAt); err != nil {
			return nil, fmt.Errorf("scan room row: %w", err)
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *repository) SearchPublicByNetwork(ctx context.Context, networkPrefix, query string, userID string) ([]PublicRoomResponse, error) {
	sql := `
		SELECT
			r.id, r.name, COALESCE(r.description, ''), u.full_name,
			(SELECT COUNT(*) FROM room_members WHERE room_id = r.id) AS member_count,
			r.password_hash IS NOT NULL AND r.password_hash != '' AS has_password,
			EXISTS(SELECT 1 FROM room_members WHERE room_id = r.id AND user_id = $3) AS is_member
		FROM rooms r
		JOIN users u ON u.id = r.owner_id
		WHERE r.visibility = 'PUBLIC'
		  AND r.status = 'ACTIVE'
		  AND r.network_prefix = $1
		  AND (r.name ILIKE '%' || $2 || '%' OR $2 = '')
		ORDER BY r.created_at DESC
	`
	rows, err := r.db.Query(ctx, sql, networkPrefix, query, userID)
	if err != nil {
		return nil, fmt.Errorf("search public rooms: %w", err)
	}
	defer rows.Close()

	var rooms []PublicRoomResponse
	for rows.Next() {
		var room PublicRoomResponse
		if err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.OwnerName, &room.MemberCount, &room.HasPassword, &room.IsMember); err != nil {
			return nil, fmt.Errorf("scan public room row: %w", err)
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *repository) Update(ctx context.Context, id, name string, description *string, visibility *string, passwordHash *string) (*Room, error) {
	if visibility != nil {
		if err := r.UpdateVisibility(ctx, id, *visibility); err != nil {
			return nil, err
		}
	}

	if passwordHash != nil {
		if err := r.SetPasswordHash(ctx, id, passwordHash); err != nil {
			return nil, err
		}
	}

	query := `
		UPDATE rooms SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, owner_id, status, visibility, password_hash, network_prefix, created_at, updated_at, closed_at
	`
	room, err := scanRoom(r.db.QueryRow(ctx, query, name, description, id))
	if err != nil {
		return nil, fmt.Errorf("update room: %w", err)
	}

	return room, nil
}

func (r *repository) UpdateVisibility(ctx context.Context, id, visibility string) error {
	_, err := r.db.Exec(ctx, `UPDATE rooms SET visibility = $1, updated_at = NOW() WHERE id = $2`, visibility, id)
	if err != nil {
		return fmt.Errorf("update visibility: %w", err)
	}
	return nil
}

func (r *repository) SetPasswordHash(ctx context.Context, id string, passwordHash *string) error {
	_, err := r.db.Exec(ctx, `UPDATE rooms SET password_hash = $1, updated_at = NOW() WHERE id = $2`, passwordHash, id)
	if err != nil {
		return fmt.Errorf("set password hash: %w", err)
	}
	return nil
}

func (r *repository) Close(ctx context.Context, id string) (*Room, error) {
	query := `
		UPDATE rooms SET status = 'CLOSED', closed_at = NOW(), updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, owner_id, status, visibility, password_hash, network_prefix, created_at, updated_at, closed_at
	`
	room, err := scanRoom(r.db.QueryRow(ctx, query, id))
	if err != nil {
		return nil, fmt.Errorf("close room: %w", err)
	}
	return room, nil
}

func (r *repository) CountMembers(ctx context.Context, roomID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM room_members WHERE room_id = $1`, roomID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count members: %w", err)
	}
	return count, nil
}

func (r *repository) GetMemberRole(ctx context.Context, roomID, userID string) (string, error) {
	var role string
	err := r.db.QueryRow(ctx, `SELECT role FROM room_members WHERE room_id = $1 AND user_id = $2`, roomID, userID).Scan(&role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("get member role: %w", err)
	}
	return role, nil
}

func (r *repository) AddMember(ctx context.Context, roomID, userID, role string) error {
	query := `
		INSERT INTO room_members (room_id, user_id, role)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, roomID, userID, role)
	if err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (r *repository) GetLeaderboard(ctx context.Context, roomID string) ([]RoomLeaderboardEntry, error) {
	query := `
		SELECT
			ROW_NUMBER() OVER (ORDER BY COALESCE(us.total_points, 0) DESC, u.full_name ASC) as rank,
			u.id,
			u.full_name,
			COALESCE(us.total_points, 0),
			COALESCE(us.predictions_count, 0),
			COALESCE(us.exact_scores_count, 0),
			rm.role
		FROM room_members rm
		JOIN users u ON u.id = rm.user_id
		LEFT JOIN user_scores us ON us.user_id = rm.user_id
		WHERE rm.room_id = $1
		ORDER BY COALESCE(us.total_points, 0) DESC, u.full_name ASC
	`
	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, fmt.Errorf("get room leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []RoomLeaderboardEntry
	for rows.Next() {
		var e RoomLeaderboardEntry
		if err := rows.Scan(&e.Rank, &e.UserID, &e.FullName, &e.TotalPoints, &e.PredictionsCount, &e.ExactScoresCount, &e.RoomRole); err != nil {
			return nil, fmt.Errorf("scan leaderboard row: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
