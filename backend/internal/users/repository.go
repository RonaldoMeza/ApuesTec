package users

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID                 string     `json:"id"`
	FullName           string     `json:"fullName"`
	Email              string     `json:"email"`
	AvatarURL          *string    `json:"avatarUrl,omitempty"`
	PasswordHash       *string    `json:"-"`
	Status             string     `json:"status"`
	FailedLoginAttempts int       `json:"-"`
	LockedUntil        *time.Time `json:"-"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

type Repository interface {
	Create(ctx context.Context, fullName, email, passwordHash string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	UpdatePassword(ctx context.Context, id, passwordHash string) error
	LockUser(ctx context.Context, id string, until time.Time) error
	ResetFailedAttempts(ctx context.Context, id string) error
	IncrementFailedAttempts(ctx context.Context, id string) (int, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func scanUser(row interface{ Scan(dest ...interface{}) error }) (*User, error) {
	u := &User{}
	err := row.Scan(
		&u.ID, &u.FullName, &u.Email, &u.AvatarURL,
		&u.PasswordHash, &u.Status, &u.FailedLoginAttempts,
		&u.LockedUntil, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) Create(ctx context.Context, fullName, email, passwordHash string) (*User, error) {
	query := `
		INSERT INTO users (full_name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, full_name, email, avatar_url, password_hash, status,
		          failed_login_attempts, locked_until, created_at, updated_at
	`
	u, err := scanUser(r.db.QueryRow(ctx, query, fullName, email, passwordHash))
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, full_name, email, avatar_url, password_hash, status,
		       failed_login_attempts, locked_until, created_at, updated_at
		FROM users WHERE email = $1
	`
	u, err := scanUser(r.db.QueryRow(ctx, query, email))
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return u, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, full_name, email, avatar_url, password_hash, status,
		       failed_login_attempts, locked_until, created_at, updated_at
		FROM users WHERE id = $1
	`
	u, err := scanUser(r.db.QueryRow(ctx, query, id))
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return u, nil
}

func (r *repository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(ctx, query, passwordHash, id)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *repository) LockUser(ctx context.Context, id string, until time.Time) error {
	query := `UPDATE users SET locked_until = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, until, id)
	if err != nil {
		return fmt.Errorf("lock user: %w", err)
	}
	return nil
}

func (r *repository) ResetFailedAttempts(ctx context.Context, id string) error {
	query := `UPDATE users SET failed_login_attempts = 0, locked_until = NULL, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("reset failed attempts: %w", err)
	}
	return nil
}

func (r *repository) IncrementFailedAttempts(ctx context.Context, id string) (int, error) {
	query := `
		UPDATE users SET failed_login_attempts = failed_login_attempts + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING failed_login_attempts
	`
	var attempts int
	err := r.db.QueryRow(ctx, query, id).Scan(&attempts)
	if err != nil {
		return 0, fmt.Errorf("increment failed attempts: %w", err)
	}
	return attempts, nil
}
