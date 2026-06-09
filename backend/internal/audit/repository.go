package audit

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ActionUserRegistered    = "USER_REGISTERED"
	ActionUserLoggedIn      = "USER_LOGGED_IN"
	ActionLoginFailed       = "LOGIN_FAILED"
	ActionUserLocked        = "USER_LOCKED"
	ActionTokenRefreshed    = "TOKEN_REFRESHED"
	ActionUserLoggedOut     = "USER_LOGGED_OUT"
	ActionPasswordChanged   = "PASSWORD_CHANGED"
	ActionGoogleAuthSuccess = "GOOGLE_AUTH_SUCCESS"
	ActionGoogleAuthFailed  = "GOOGLE_AUTH_FAILED"

	ActionPredictionCreated = "PREDICTION_CREATED"
	ActionPredictionUpdated = "PREDICTION_UPDATED"

	ActionMatchResultUpdated = "MATCH_RESULT_UPDATED"
	ActionScoreCalculated    = "SCORE_CALCULATED"

	RoomCreated        = "ROOM_CREATED"
	RoomUpdated        = "ROOM_UPDATED"
	RoomClosed         = "ROOM_CLOSED"
	RoomMemberJoined   = "ROOM_MEMBER_JOINED"
	RoomMemberLeft     = "ROOM_MEMBER_LEFT"
	RoomMemberRemoved  = "ROOM_MEMBER_REMOVED"
	RoomRoleChanged    = "ROOM_ROLE_CHANGED"
	RoomInviteCreated  = "ROOM_INVITE_CREATED"
	RoomInviteUsed     = "ROOM_INVITE_USED"
	RoomInviteRevoked  = "ROOM_INVITE_REVOKED"
)

type Repository interface {
	Log(ctx context.Context, userID *string, action, entity, entityID string, ipAddress, userAgent *string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Log(ctx context.Context, userID *string, action, entity, entityID string, ipAddress, userAgent *string) error {
	query := `
		INSERT INTO audit_logs (user_id, action, entity, entity_id, ip_address, user_agent)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6)
	`
	_, err := r.db.Exec(ctx, query, userID, action, entity, entityID, ipAddress, userAgent)
	if err != nil {
		return fmt.Errorf("audit log: %w", err)
	}
	return nil
}
