package roles

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Role struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type Repository interface {
	GetByName(ctx context.Context, name string) (*Role, error)
	GetUserRoles(ctx context.Context, userID string) ([]Role, error)
	AssignRole(ctx context.Context, userID, roleID string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `SELECT id, name, description FROM roles WHERE name = $1`
	role := &Role{}
	err := r.db.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		return nil, fmt.Errorf("get role by name: %w", err)
	}
	return role, nil
}

func (r *repository) GetUserRoles(ctx context.Context, userID string) ([]Role, error) {
	query := `
		SELECT r.id, r.name, r.description
		FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.name
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return nil, fmt.Errorf("scan role row: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *repository) AssignRole(ctx context.Context, userID, roleID string) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("assign role: %w", err)
	}
	return nil
}
