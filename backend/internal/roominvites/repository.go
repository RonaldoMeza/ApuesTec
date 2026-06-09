package roominvites

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, roomID, code, qrPayload, createdBy string, expiresAt time.Time) (*RoomInvite, error)
	FindByCode(ctx context.Context, code string) (*RoomInvite, error)
	RevokeActiveInvites(ctx context.Context, roomID string) error
	MarkUsed(ctx context.Context, id string) error
	GetActiveInvite(ctx context.Context, roomID string) (*RoomInvite, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func generateCode() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random code: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (r *repository) Create(ctx context.Context, roomID, code, qrPayload, createdBy string, expiresAt time.Time) (*RoomInvite, error) {
	invite := &RoomInvite{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO room_invites (room_id, code, qr_payload, created_by, expires_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, room_id, code, qr_payload, created_by, expires_at, used_at, revoked_at, created_at`,
		roomID, code, qrPayload, createdBy, expiresAt,
	).Scan(&invite.ID, &invite.RoomID, &invite.Code, &invite.QRPayload, &invite.CreatedBy,
		&invite.ExpiresAt, &invite.UsedAt, &invite.RevokedAt, &invite.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create room invite: %w", err)
	}
	return invite, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*RoomInvite, error) {
	invite := &RoomInvite{}
	err := r.db.QueryRow(ctx,
		`SELECT id, room_id, code, qr_payload, created_by, expires_at, used_at, revoked_at, created_at
		 FROM room_invites WHERE code = $1`,
		code,
	).Scan(&invite.ID, &invite.RoomID, &invite.Code, &invite.QRPayload, &invite.CreatedBy,
		&invite.ExpiresAt, &invite.UsedAt, &invite.RevokedAt, &invite.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("find invite by code: %w", err)
	}
	return invite, nil
}

func (r *repository) RevokeActiveInvites(ctx context.Context, roomID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE room_invites
		 SET revoked_at = NOW()
		 WHERE room_id = $1 AND used_at IS NULL AND revoked_at IS NULL AND expires_at > NOW()`,
		roomID,
	)
	if err != nil {
		return fmt.Errorf("revoke active invites: %w", err)
	}
	return nil
}

func (r *repository) MarkUsed(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE room_invites SET used_at = NOW() WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("mark invite used: %w", err)
	}
	return nil
}

func (r *repository) GetActiveInvite(ctx context.Context, roomID string) (*RoomInvite, error) {
	invite := &RoomInvite{}
	err := r.db.QueryRow(ctx,
		`SELECT id, room_id, code, qr_payload, created_by, expires_at, used_at, revoked_at, created_at
		 FROM room_invites
		 WHERE room_id = $1 AND used_at IS NULL AND revoked_at IS NULL AND expires_at > NOW()
		 ORDER BY created_at DESC LIMIT 1`,
		roomID,
	).Scan(&invite.ID, &invite.RoomID, &invite.Code, &invite.QRPayload, &invite.CreatedBy,
		&invite.ExpiresAt, &invite.UsedAt, &invite.RevokedAt, &invite.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get active invite: %w", err)
	}
	return invite, nil
}
