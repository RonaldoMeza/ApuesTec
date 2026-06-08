package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	UserAgent *string
	IPAddress *string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type AuthAccount struct {
	ID             string
	UserID         string
	Provider       string
	ProviderUserID string
	ProviderEmail  *string
	CreatedAt      time.Time
}

type Repository interface {
	CreateRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time, userAgent, ipAddress *string) error
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id string) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID string) error
	CreateAuthAccount(ctx context.Context, userID, provider, providerUserID, providerEmail string) error
	FindAuthAccount(ctx context.Context, provider, providerUserID string) (*AuthAccount, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) CreateRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time, userAgent, ipAddress *string) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, userID, tokenHash, userAgent, ipAddress, expiresAt)
	if err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}
	return nil
}

func (r *repository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, user_agent, ip_address::text, expires_at, revoked_at, created_at
		FROM refresh_tokens WHERE token_hash = $1
	`
	t := &RefreshToken{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.UserAgent,
		&t.IPAddress, &t.ExpiresAt, &t.RevokedAt, &t.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find refresh token: %w", err)
	}
	return t, nil
}

func (r *repository) RevokeRefreshToken(ctx context.Context, id string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("refresh token not found or already revoked")
	}
	return nil
}

func (r *repository) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("revoke all user refresh tokens: %w", err)
	}
	return nil
}

func (r *repository) CreateAuthAccount(ctx context.Context, userID, provider, providerUserID, providerEmail string) error {
	query := `
		INSERT INTO auth_accounts (user_id, provider, provider_user_id, provider_email)
		VALUES ($1, $2, $3, NULLIF($4, ''))
		ON CONFLICT (user_id, provider) DO UPDATE SET provider_user_id = EXCLUDED.provider_user_id
	`
	_, err := r.db.Exec(ctx, query, userID, provider, providerUserID, providerEmail)
	if err != nil {
		return fmt.Errorf("create auth account: %w", err)
	}
	return nil
}

func (r *repository) FindAuthAccount(ctx context.Context, provider, providerUserID string) (*AuthAccount, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_email, created_at
		FROM auth_accounts WHERE provider = $1 AND provider_user_id = $2
	`
	a := &AuthAccount{}
	err := r.db.QueryRow(ctx, query, provider, providerUserID).Scan(
		&a.ID, &a.UserID, &a.Provider, &a.ProviderUserID, &a.ProviderEmail, &a.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find auth account: %w", err)
	}
	return a, nil
}
