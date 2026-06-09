package predictions

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, userID, matchID string, homeScore, awayScore int) (*Prediction, error)
	FindByUserAndMatch(ctx context.Context, userID, matchID string) (*Prediction, error)
	FindByUser(ctx context.Context, userID string) ([]Prediction, error)
	Update(ctx context.Context, id string, homeScore, awayScore int) (*Prediction, error)
	Lock(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func scanPrediction(row pgx.Row) (*Prediction, error) {
	p := &Prediction{}
	err := row.Scan(
		&p.ID, &p.UserID, &p.MatchID,
		&p.PredictedHomeScore, &p.PredictedAwayScore,
		&p.IsExactScore, &p.IsWinnerCorrect, &p.IsGoalDifferenceCorrect,
		&p.BasePoints, &p.EarlyBonusPoints, &p.StreakBonusPoints, &p.TotalPoints,
		&p.LockedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func scanPredictions(rows pgx.Rows) ([]Prediction, error) {
	defer rows.Close()
	var predictions []Prediction
	for rows.Next() {
		var p Prediction
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.MatchID,
			&p.PredictedHomeScore, &p.PredictedAwayScore,
			&p.IsExactScore, &p.IsWinnerCorrect, &p.IsGoalDifferenceCorrect,
			&p.BasePoints, &p.EarlyBonusPoints, &p.StreakBonusPoints, &p.TotalPoints,
			&p.LockedAt, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan prediction row: %w", err)
		}
		predictions = append(predictions, p)
	}
	return predictions, rows.Err()
}

func (r *repository) Create(ctx context.Context, userID, matchID string, homeScore, awayScore int) (*Prediction, error) {
	query := `
		INSERT INTO predictions (user_id, match_id, predicted_home_score, predicted_away_score)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, match_id,
		          predicted_home_score, predicted_away_score,
		          is_exact_score, is_winner_correct, is_goal_difference_correct,
		          base_points, early_bonus_points, streak_bonus_points, total_points,
		          locked_at, created_at, updated_at
	`
	p, err := scanPrediction(r.db.QueryRow(ctx, query, userID, matchID, homeScore, awayScore))
	if err != nil {
		return nil, fmt.Errorf("create prediction: %w", err)
	}
	return p, nil
}

func (r *repository) FindByUserAndMatch(ctx context.Context, userID, matchID string) (*Prediction, error) {
	query := `
		SELECT id, user_id, match_id,
		       predicted_home_score, predicted_away_score,
		       is_exact_score, is_winner_correct, is_goal_difference_correct,
		       base_points, early_bonus_points, streak_bonus_points, total_points,
		       locked_at, created_at, updated_at
		FROM predictions WHERE user_id = $1 AND match_id = $2
	`
	p, err := scanPrediction(r.db.QueryRow(ctx, query, userID, matchID))
	if err != nil {
		return nil, fmt.Errorf("find prediction by user and match: %w", err)
	}
	return p, nil
}

func (r *repository) FindByUser(ctx context.Context, userID string) ([]Prediction, error) {
	query := `
		SELECT id, user_id, match_id,
		       predicted_home_score, predicted_away_score,
		       is_exact_score, is_winner_correct, is_goal_difference_correct,
		       base_points, early_bonus_points, streak_bonus_points, total_points,
		       locked_at, created_at, updated_at
		FROM predictions WHERE user_id = $1 ORDER BY updated_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find predictions by user: %w", err)
	}
	return scanPredictions(rows)
}

func (r *repository) Update(ctx context.Context, id string, homeScore, awayScore int) (*Prediction, error) {
	query := `
		UPDATE predictions SET predicted_home_score = $1, predicted_away_score = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, user_id, match_id,
		          predicted_home_score, predicted_away_score,
		          is_exact_score, is_winner_correct, is_goal_difference_correct,
		          base_points, early_bonus_points, streak_bonus_points, total_points,
		          locked_at, created_at, updated_at
	`
	p, err := scanPrediction(r.db.QueryRow(ctx, query, homeScore, awayScore, id))
	if err != nil {
		return nil, fmt.Errorf("update prediction: %w", err)
	}
	return p, nil
}

func (r *repository) Lock(ctx context.Context, id string) error {
	query := `UPDATE predictions SET locked_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("lock prediction: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM predictions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete prediction: %w", err)
	}
	return nil
}
