package userstats

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetUserScore(ctx context.Context, userID string) (*UserScore, error)
	GetPredictionStats(ctx context.Context, userID string) (*PredictionStats, error)
	ListScoreEvents(ctx context.Context, userID string) ([]ScoreEventRow, error)
}

type UserScore struct {
	TotalPoints                int
	PredictionsCount           int
	ExactScoresCount           int
	WinnerCorrectCount         int
	GoalDifferenceCorrectCount int
}

type PredictionStats struct {
	CurrentStreak int
	BestStreak    int
}

type ScoreEventRow struct {
	ID          string
	MatchID     string
	EventType   string
	Points      int
	Description *string
	CreatedAt   string
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserScore(ctx context.Context, userID string) (*UserScore, error) {
	query := `
		SELECT
			COALESCE(us.total_points, 0),
			COALESCE(us.predictions_count, 0),
			COALESCE(us.exact_scores_count, 0),
			COALESCE(us.winner_correct_count, 0),
			COALESCE(us.goal_difference_correct_count, 0)
		FROM user_scores us WHERE us.user_id = $1
	`
	s := &UserScore{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&s.TotalPoints, &s.PredictionsCount,
		&s.ExactScoresCount, &s.WinnerCorrectCount, &s.GoalDifferenceCorrectCount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &UserScore{}, nil
		}
		return nil, fmt.Errorf("get user score: %w", err)
	}
	return s, nil
}

func (r *repository) GetPredictionStats(ctx context.Context, userID string) (*PredictionStats, error) {
	query := `
		WITH ordered_desc AS (
			SELECT p.is_winner_correct,
			       ROW_NUMBER() OVER (ORDER BY m.starts_at DESC) as rn
			FROM predictions p
			JOIN matches m ON m.id = p.match_id
			WHERE p.user_id = $1 AND m.status = 'FINISHED'
		)
		SELECT
			COALESCE((
				SELECT COUNT(*) FROM ordered_desc
				WHERE is_winner_correct = true AND rn < (
					SELECT COALESCE(MIN(rn), 99999999) FROM ordered_desc WHERE is_winner_correct = false
				)
			), 0) as current_streak,
			COALESCE((
				SELECT MAX(cnt) FROM (
					SELECT COUNT(*) as cnt
					FROM (
						SELECT rn, rn - ROW_NUMBER() OVER (ORDER BY rn) as grp
						FROM ordered_desc WHERE is_winner_correct = true
					) t
					GROUP BY grp
				) streaks
			), 0) as best_streak
	`
	s := &PredictionStats{}
	err := r.db.QueryRow(ctx, query, userID).Scan(&s.CurrentStreak, &s.BestStreak)
	if err != nil {
		return &PredictionStats{}, nil
	}
	return s, nil
}

func (r *repository) ListScoreEvents(ctx context.Context, userID string) ([]ScoreEventRow, error) {
	query := `
		SELECT id, match_id, event_type, points, description,
		       TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at
		FROM score_events WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list score events: %w", err)
	}
	defer rows.Close()

	var events []ScoreEventRow
	for rows.Next() {
		var e ScoreEventRow
		if err := rows.Scan(&e.ID, &e.MatchID, &e.EventType, &e.Points, &e.Description, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan score event: %w", err)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
