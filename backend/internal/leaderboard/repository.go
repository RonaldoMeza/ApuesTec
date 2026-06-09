package leaderboard

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetGlobalLeaderboard(ctx context.Context, limit, offset int) ([]LeaderboardEntry, int, error)
	GetUserRank(ctx context.Context, userID string) (int, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetGlobalLeaderboard(ctx context.Context, limit, offset int) ([]LeaderboardEntry, int, error) {
	countQuery := `SELECT COUNT(*) FROM user_scores`
	var total int
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count leaderboard: %w", err)
	}

	query := `
		SELECT
			u.id,
			u.full_name,
			u.email,
			COALESCE(us.total_points, 0),
			COALESCE(us.predictions_count, 0),
			COALESCE(us.exact_scores_count, 0),
			COALESCE(us.winner_correct_count, 0),
			ROW_NUMBER() OVER (ORDER BY COALESCE(us.total_points, 0) DESC, u.full_name ASC)
		FROM users u
		LEFT JOIN user_scores us ON us.user_id = u.id
		ORDER BY COALESCE(us.total_points, 0) DESC, u.full_name ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []LeaderboardEntry
	for rows.Next() {
		var e LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.FullName, &e.Email, &e.TotalPoints, &e.PredictionsCount, &e.ExactScores, &e.WinnerCorrect, &e.Rank); err != nil {
			return nil, 0, fmt.Errorf("scan leaderboard row: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, total, rows.Err()
}

func (r *repository) GetUserRank(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT rank FROM (
			SELECT u.id, ROW_NUMBER() OVER (ORDER BY COALESCE(us.total_points, 0) DESC, u.full_name ASC) as rank
			FROM users u
			LEFT JOIN user_scores us ON us.user_id = u.id
		) ranked WHERE id = $1
	`
	var rank int
	err := r.db.QueryRow(ctx, query, userID).Scan(&rank)
	if err != nil {
		return 0, fmt.Errorf("get user rank: %w", err)
	}
	return rank, nil
}
