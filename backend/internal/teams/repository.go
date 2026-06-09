package teams

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, name, code, country string, flagURL *string) (*Team, error)
	FindAll(ctx context.Context) ([]Team, error)
	FindByID(ctx context.Context, id string) (*Team, error)
	Update(ctx context.Context, id, name, code, country string, flagURL *string) (*Team, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func scanTeam(row interface{ Scan(dest ...interface{}) error }) (*Team, error) {
	t := &Team{}
	err := row.Scan(&t.ID, &t.Name, &t.Code, &t.Country, &t.FlagURL, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *repository) Create(ctx context.Context, name, code, country string, flagURL *string) (*Team, error) {
	query := `
		INSERT INTO teams (name, code, country, flag_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, code, country, flag_url, created_at, updated_at
	`
	t, err := scanTeam(r.db.QueryRow(ctx, query, name, code, country, flagURL))
	if err != nil {
		return nil, fmt.Errorf("create team: %w", err)
	}
	return t, nil
}

func (r *repository) FindAll(ctx context.Context) ([]Team, error) {
	query := `SELECT id, name, code, country, flag_url, created_at, updated_at FROM teams ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find all teams: %w", err)
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Code, &t.Country, &t.FlagURL, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan team row: %w", err)
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (r *repository) FindByID(ctx context.Context, id string) (*Team, error) {
	query := `SELECT id, name, code, country, flag_url, created_at, updated_at FROM teams WHERE id = $1`
	t, err := scanTeam(r.db.QueryRow(ctx, query, id))
	if err != nil {
		return nil, fmt.Errorf("find team by id: %w", err)
	}
	return t, nil
}

func (r *repository) Update(ctx context.Context, id, name, code, country string, flagURL *string) (*Team, error) {
	query := `
		UPDATE teams SET name = $1, code = $2, country = $3, flag_url = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, name, code, country, flag_url, created_at, updated_at
	`
	t, err := scanTeam(r.db.QueryRow(ctx, query, name, code, country, flagURL, id))
	if err != nil {
		return nil, fmt.Errorf("update team: %w", err)
	}
	return t, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM teams WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete team: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("team not found")
	}
	return nil
}
