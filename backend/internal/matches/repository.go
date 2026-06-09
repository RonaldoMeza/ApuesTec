package matches

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, homeTeamID, awayTeamID string, startsAt time.Time) (*Match, error)
	FindAll(ctx context.Context) ([]Match, error)
	FindUpcoming(ctx context.Context) ([]Match, error)
	FindFinished(ctx context.Context) ([]Match, error)
	FindByID(ctx context.Context, id string) (*Match, error)
	Update(ctx context.Context, id, homeTeamID, awayTeamID string, startsAt time.Time) (*Match, error)
	UpdateStatus(ctx context.Context, id, status string) (*Match, error)
	UpdateResult(ctx context.Context, id string, homeScore, awayScore int) (*Match, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func scanMatch(row pgx.Row) (*Match, error) {
	m := &Match{}
	err := row.Scan(&m.ID, &m.HomeTeamID, &m.AwayTeamID, &m.StartsAt, &m.Status, &m.HomeScore, &m.AwayScore, &m.LockedAt, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *repository) Create(ctx context.Context, homeTeamID, awayTeamID string, startsAt time.Time) (*Match, error) {
	query := `
		INSERT INTO matches (home_team_id, away_team_id, starts_at)
		VALUES ($1, $2, $3)
		RETURNING id, home_team_id, away_team_id, starts_at, status, home_score, away_score, locked_at, created_at, updated_at
	`
	m, err := scanMatch(r.db.QueryRow(ctx, query, homeTeamID, awayTeamID, startsAt))
	if err != nil {
		return nil, fmt.Errorf("create match: %w", err)
	}
	return m, nil
}

func (r *repository) FindAll(ctx context.Context) ([]Match, error) {
	query := `SELECT id, home_team_id, away_team_id, starts_at, status, home_score, away_score, locked_at, created_at, updated_at FROM matches ORDER BY starts_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find all matches: %w", err)
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.HomeTeamID, &m.AwayTeamID, &m.StartsAt, &m.Status, &m.HomeScore, &m.AwayScore, &m.LockedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan match row: %w", err)
		}
		matches = append(matches, m)
	}
	return matches, nil
}

func (r *repository) FindUpcoming(ctx context.Context) ([]Match, error) {
	query := `SELECT id, home_team_id, away_team_id, starts_at, status, home_score, away_score, locked_at, created_at, updated_at FROM matches WHERE status IN ('SCHEDULED', 'LOCKED') ORDER BY starts_at ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find upcoming matches: %w", err)
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.HomeTeamID, &m.AwayTeamID, &m.StartsAt, &m.Status, &m.HomeScore, &m.AwayScore, &m.LockedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan match row: %w", err)
		}
		matches = append(matches, m)
	}
	return matches, nil
}

func (r *repository) FindFinished(ctx context.Context) ([]Match, error) {
	query := `SELECT id, home_team_id, away_team_id, starts_at, status, home_score, away_score, locked_at, created_at, updated_at FROM matches WHERE status IN ('FINISHED', 'CANCELLED') ORDER BY starts_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find finished matches: %w", err)
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.HomeTeamID, &m.AwayTeamID, &m.StartsAt, &m.Status, &m.HomeScore, &m.AwayScore, &m.LockedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan match row: %w", err)
		}
		matches = append(matches, m)
	}
	return matches, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Match, error) {
	query := `SELECT id, home_team_id, away_team_id, starts_at, status, home_score, away_score, locked_at, created_at, updated_at FROM matches WHERE id = $1`
	m, err := scanMatch(r.db.QueryRow(ctx, query, id))
	if err != nil {
		return nil, fmt.Errorf("find match by id: %w", err)
	}
	return m, nil
}

func (r *repository) Update(ctx context.Context, id, homeTeamID, awayTeamID string, startsAt time.Time) (*Match, error) {
	query := `
		UPDATE matches SET home_team_id = $1, away_team_id = $2, starts_at = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, home_team_id, away_team_id, starts_at, status, home_score, away_score, locked_at, created_at, updated_at
	`
	m, err := scanMatch(r.db.QueryRow(ctx, query, homeTeamID, awayTeamID, startsAt, id))
	if err != nil {
		return nil, fmt.Errorf("update match: %w", err)
	}
	return m, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id, status string) (*Match, error) {
	query := `
		UPDATE matches SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, home_team_id, away_team_id, starts_at, status, home_score, away_score, locked_at, created_at, updated_at
	`
	m, err := scanMatch(r.db.QueryRow(ctx, query, status, id))
	if err != nil {
		return nil, fmt.Errorf("update match status: %w", err)
	}
	return m, nil
}

func (r *repository) UpdateResult(ctx context.Context, id string, homeScore, awayScore int) (*Match, error) {
	query := `
		UPDATE matches SET home_score = $1, away_score = $2, status = 'FINISHED', updated_at = NOW()
		WHERE id = $3
		RETURNING id, home_team_id, away_team_id, starts_at, status, home_score, away_score, locked_at, created_at, updated_at
	`
	m, err := scanMatch(r.db.QueryRow(ctx, query, homeScore, awayScore, id))
	if err != nil {
		return nil, fmt.Errorf("update match result: %w", err)
	}
	return m, nil
}
