package scoring

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindPredictionsByMatch(ctx context.Context, matchID string) ([]ScoredPrediction, error)
	UpdatePredictionScore(ctx context.Context, id string, isExact, isWinner, isGoalDiff bool, basePts, earlyBonus, streakBonus, total int) error
	LockPrediction(ctx context.Context, id string) error
	CreateScoreEvent(ctx context.Context, userID, matchID, predictionID, eventType string, points int, description *string) error
	DeleteScoreEventsByMatch(ctx context.Context, matchID string) error
	GetUserScore(ctx context.Context, userID string) (*UserScore, error)
	UpsertUserScore(ctx context.Context, score UserScore) error
	RebuildAllUserScores(ctx context.Context) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func scanScoredPrediction(row pgx.Row) (*ScoredPrediction, error) {
	p := &ScoredPrediction{}
	err := row.Scan(
		&p.ID, &p.UserID, &p.MatchID,
		&p.PredictedHomeScore, &p.PredictedAwayScore,
		&p.IsExactScore, &p.IsWinnerCorrect, &p.IsGoalDifferenceCorrect,
		&p.BasePoints, &p.EarlyBonusPoints, &p.StreakBonusPoints, &p.TotalPoints,
		&p.LockedAt, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func scanScoredPredictions(rows pgx.Rows) ([]ScoredPrediction, error) {
	defer rows.Close()
	var predictions []ScoredPrediction
	for rows.Next() {
		var p ScoredPrediction
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.MatchID,
			&p.PredictedHomeScore, &p.PredictedAwayScore,
			&p.IsExactScore, &p.IsWinnerCorrect, &p.IsGoalDifferenceCorrect,
			&p.BasePoints, &p.EarlyBonusPoints, &p.StreakBonusPoints, &p.TotalPoints,
			&p.LockedAt, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan scored prediction row: %w", err)
		}
		predictions = append(predictions, p)
	}
	return predictions, rows.Err()
}

func (r *repository) FindPredictionsByMatch(ctx context.Context, matchID string) ([]ScoredPrediction, error) {
	query := `
		SELECT id, user_id, match_id,
		       predicted_home_score, predicted_away_score,
		       is_exact_score, is_winner_correct, is_goal_difference_correct,
		       base_points, early_bonus_points, streak_bonus_points, total_points,
		       locked_at, created_at
		FROM predictions WHERE match_id = $1 ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, matchID)
	if err != nil {
		return nil, fmt.Errorf("find predictions by match: %w", err)
	}
	return scanScoredPredictions(rows)
}

func (r *repository) UpdatePredictionScore(ctx context.Context, id string, isExact, isWinner, isGoalDiff bool, basePts, earlyBonus, streakBonus, total int) error {
	query := `
		UPDATE predictions SET
			is_exact_score = $1, is_winner_correct = $2, is_goal_difference_correct = $3,
			base_points = $4, early_bonus_points = $5, streak_bonus_points = $6,
			total_points = $7, updated_at = NOW()
		WHERE id = $8
	`
	_, err := r.db.Exec(ctx, query, isExact, isWinner, isGoalDiff, basePts, earlyBonus, streakBonus, total, id)
	if err != nil {
		return fmt.Errorf("update prediction score: %w", err)
	}
	return nil
}

func (r *repository) LockPrediction(ctx context.Context, id string) error {
	query := `UPDATE predictions SET locked_at = COALESCE(locked_at, NOW()), updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("lock prediction: %w", err)
	}
	return nil
}

func (r *repository) CreateScoreEvent(ctx context.Context, userID, matchID, predictionID, eventType string, points int, description *string) error {
	query := `
		INSERT INTO score_events (user_id, match_id, prediction_id, event_type, points, description)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, userID, matchID, predictionID, eventType, points, description)
	if err != nil {
		return fmt.Errorf("create score event: %w", err)
	}
	return nil
}

func (r *repository) DeleteScoreEventsByMatch(ctx context.Context, matchID string) error {
	query := `DELETE FROM score_events WHERE match_id = $1`
	_, err := r.db.Exec(ctx, query, matchID)
	if err != nil {
		return fmt.Errorf("delete score events by match: %w", err)
	}
	return nil
}

func (r *repository) GetUserScore(ctx context.Context, userID string) (*UserScore, error) {
	query := `
		SELECT user_id, total_points, predictions_count, exact_scores_count,
		       winner_correct_count, goal_difference_correct_count,
		       streak_count, last_scored_at, updated_at
		FROM user_scores WHERE user_id = $1
	`
	s := &UserScore{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&s.UserID, &s.TotalPoints, &s.PredictionsCount, &s.ExactScoresCount,
		&s.WinnerCorrectCount, &s.GoalDifferenceCorrectCount,
		&s.StreakCount, &s.LastScoredAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get user score: %w", err)
	}
	return s, nil
}

func (r *repository) UpsertUserScore(ctx context.Context, score UserScore) error {
	query := `
		INSERT INTO user_scores (user_id, total_points, predictions_count, exact_scores_count,
		                         winner_correct_count, goal_difference_correct_count,
		                         streak_count, last_scored_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			total_points = EXCLUDED.total_points,
			predictions_count = EXCLUDED.predictions_count,
			exact_scores_count = EXCLUDED.exact_scores_count,
			winner_correct_count = EXCLUDED.winner_correct_count,
			goal_difference_correct_count = EXCLUDED.goal_difference_correct_count,
			last_scored_at = EXCLUDED.last_scored_at,
			updated_at = NOW()
	`
	_, err := r.db.Exec(ctx, query,
		score.UserID, score.TotalPoints, score.PredictionsCount, score.ExactScoresCount,
		score.WinnerCorrectCount, score.GoalDifferenceCorrectCount,
		score.StreakCount, score.LastScoredAt,
	)
	if err != nil {
		return fmt.Errorf("upsert user score: %w", err)
	}
	return nil
}

func (r *repository) RebuildAllUserScores(ctx context.Context) error {
	query := `
		INSERT INTO user_scores (user_id, total_points, predictions_count, exact_scores_count,
		                         winner_correct_count, goal_difference_correct_count,
		                         streak_count, last_scored_at, updated_at)
		SELECT
			u.id,
			COALESCE(SUM(p.total_points), 0),
			COUNT(p.id) FILTER (WHERE p.total_points > 0),
			COUNT(p.id) FILTER (WHERE p.is_exact_score = true),
			COUNT(p.id) FILTER (WHERE p.is_winner_correct = true),
			COUNT(p.id) FILTER (WHERE p.is_goal_difference_correct = true),
			0,
			MAX(se.created_at),
			NOW()
		FROM users u
		LEFT JOIN predictions p ON p.user_id = u.id
		LEFT JOIN score_events se ON se.user_id = u.id
		GROUP BY u.id
		ON CONFLICT (user_id) DO UPDATE SET
			total_points = EXCLUDED.total_points,
			predictions_count = EXCLUDED.predictions_count,
			exact_scores_count = EXCLUDED.exact_scores_count,
			winner_correct_count = EXCLUDED.winner_correct_count,
			goal_difference_correct_count = EXCLUDED.goal_difference_correct_count,
			last_scored_at = EXCLUDED.last_scored_at,
			updated_at = NOW()
	`
	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("rebuild user scores: %w", err)
	}
	return nil
}
