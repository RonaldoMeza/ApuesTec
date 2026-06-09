package scoring

import (
	"context"
	"fmt"
	"time"
)

const (
	PointsExactScore     = 5
	PointsWinnerCorrect  = 3
	PointsGoalDifference = 2
	PointsEarlyBonus     = 1
	PointsStreakBonus    = 2
	StreakInterval       = 3
	EarlyThreshold       = 24 * time.Hour
)

type ScoreResult struct {
	PredictionID            string
	UserID                  string
	IsExactScore            bool
	IsWinnerCorrect         bool
	IsGoalDifferenceCorrect bool
	BasePoints              int
	EarlyBonus              int
	StreakBonus             int
	TotalPoints             int
}

type MatchInfo struct {
	ID        string
	HomeScore int
	AwayScore int
	Status    string
	StartsAt  time.Time
}

type Service interface {
	CalculateAndSave(ctx context.Context, match MatchInfo) ([]ScoreResult, error)
	RecalculateMatch(ctx context.Context, match MatchInfo) ([]ScoreResult, error)
	RebuildAllUserScores(ctx context.Context) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CalculateAndSave(ctx context.Context, match MatchInfo) ([]ScoreResult, error) {
	if match.Status == "CANCELLED" {
		return nil, ErrMatchCancelled
	}

	predictions, err := s.repo.FindPredictionsByMatch(ctx, match.ID)
	if err != nil {
		return nil, fmt.Errorf("find predictions: %w", err)
	}

	if len(predictions) == 0 {
		return []ScoreResult{}, nil
	}

	results := make([]ScoreResult, 0, len(predictions))

	for _, pred := range predictions {
		result := s.calculateSingleScore(pred, match)
		results = append(results, result)

		if err := s.repo.UpdatePredictionScore(
			ctx, pred.ID,
			result.IsExactScore, result.IsWinnerCorrect, result.IsGoalDifferenceCorrect,
			result.BasePoints, 0, 0, result.TotalPoints,
		); err != nil {
			return nil, fmt.Errorf("update prediction score: %w", err)
		}

		if err := s.repo.LockPrediction(ctx, pred.ID); err != nil {
			return nil, fmt.Errorf("lock prediction: %w", err)
		}

		if result.IsExactScore {
			desc := fmt.Sprintf("Marcador exacto: %d-%d", match.HomeScore, match.AwayScore)
			if err := s.repo.CreateScoreEvent(ctx, pred.UserID, match.ID, pred.ID, EventTypeExactScore, PointsExactScore, &desc); err != nil {
				return nil, fmt.Errorf("create event: %w", err)
			}
		}
		if result.IsWinnerCorrect && !result.IsExactScore {
			desc := "Ganador/empate correcto"
			if err := s.repo.CreateScoreEvent(ctx, pred.UserID, match.ID, pred.ID, EventTypeWinnerCorrect, PointsWinnerCorrect, &desc); err != nil {
				return nil, fmt.Errorf("create event: %w", err)
			}
		}
		if result.IsGoalDifferenceCorrect && !result.IsExactScore {
			diff := match.HomeScore - match.AwayScore
			desc := fmt.Sprintf("Diferencia de goles correcta: %+d", diff)
			if err := s.repo.CreateScoreEvent(ctx, pred.UserID, match.ID, pred.ID, EventTypeGoalDifference, PointsGoalDifference, &desc); err != nil {
				return nil, fmt.Errorf("create event: %w", err)
			}
		}
		if result.EarlyBonus > 0 {
			desc := "Prediccion anticipada (+24h antes del inicio)"
			if err := s.repo.CreateScoreEvent(ctx, pred.UserID, match.ID, pred.ID, EventTypeEarlyBonus, PointsEarlyBonus, &desc); err != nil {
				return nil, fmt.Errorf("create event: %w", err)
			}
		}
		if result.StreakBonus > 0 {
			desc := fmt.Sprintf("Racha de %d aciertos consecutivos", s.countStreak(pred.UserID, results))
			if err := s.repo.CreateScoreEvent(ctx, pred.UserID, match.ID, pred.ID, EventTypeStreakBonus, result.StreakBonus, &desc); err != nil {
				return nil, fmt.Errorf("create event: %w", err)
			}
		}

		if err := s.updateUserScore(ctx, pred.UserID, result.TotalPoints); err != nil {
			return nil, fmt.Errorf("update user score: %w", err)
		}
	}

	return results, nil
}

func (s *service) calculateSingleScore(pred ScoredPrediction, match MatchInfo) ScoreResult {
	result := ScoreResult{
		PredictionID: pred.ID,
		UserID:       pred.UserID,
	}

	predictedHome := pred.PredictedHomeScore
	predictedAway := pred.PredictedAwayScore
	realHome := match.HomeScore
	realAway := match.AwayScore

	predictedDiff := predictedHome - predictedAway
	realDiff := realHome - realAway

	isExactScore := predictedHome == realHome && predictedAway == realAway
	isWinnerCorrect := isExactScore || (realHome > realAway && predictedHome > predictedAway) ||
		(realHome < realAway && predictedHome < predictedAway) ||
		(realHome == realAway && predictedHome == predictedAway)
	isGoalDiffCorrect := isExactScore || predictedDiff == realDiff

	result.IsExactScore = isExactScore
	result.IsWinnerCorrect = isWinnerCorrect
	result.IsGoalDifferenceCorrect = isGoalDiffCorrect

	if isExactScore {
		result.BasePoints = PointsExactScore + PointsWinnerCorrect + PointsGoalDifference
	} else {
		if isWinnerCorrect {
			result.BasePoints += PointsWinnerCorrect
		}
		if isGoalDiffCorrect {
			result.BasePoints += PointsGoalDifference
		}
	}

	if pred.CreatedAt.Before(match.StartsAt.Add(-EarlyThreshold)) {
		result.EarlyBonus = PointsEarlyBonus
		result.BasePoints += PointsEarlyBonus
	}

	result.TotalPoints = result.BasePoints

	return result
}

func (s *service) countStreak(userID string, results []ScoreResult) int {
	streak := 0
	for i := len(results) - 1; i >= 0; i-- {
		if results[i].UserID != userID {
			continue
		}
		if results[i].IsWinnerCorrect {
			streak++
		} else {
			break
		}
	}
	return streak
}

func (s *service) updateUserScore(ctx context.Context, userID string, pointsToAdd int) error {
	existing, err := s.repo.GetUserScore(ctx, userID)
	if err != nil {
		existing = &UserScore{UserID: userID}
	}

	now := time.Now()
	existing.TotalPoints += pointsToAdd
	existing.PredictionsCount++
	existing.LastScoredAt = &now

	return s.repo.UpsertUserScore(ctx, *existing)
}

func (s *service) RecalculateMatch(ctx context.Context, match MatchInfo) ([]ScoreResult, error) {
	if match.Status == "CANCELLED" {
		return nil, ErrMatchCancelled
	}

	predictions, err := s.repo.FindPredictionsByMatch(ctx, match.ID)
	if err != nil {
		return nil, fmt.Errorf("find predictions: %w", err)
	}

	if len(predictions) == 0 {
		return []ScoreResult{}, nil
	}

	if err := s.repo.DeleteScoreEventsByMatch(ctx, match.ID); err != nil {
		return nil, fmt.Errorf("delete old score events: %w", err)
	}

	for _, pred := range predictions {
		if err := s.repo.UpdatePredictionScore(ctx, pred.ID, false, false, false, 0, 0, 0, 0); err != nil {
			return nil, fmt.Errorf("reset prediction: %w", err)
		}
	}

	if err := s.RebuildAllUserScores(ctx); err != nil {
		return nil, fmt.Errorf("rebuild user scores: %w", err)
	}

	return s.CalculateAndSave(ctx, match)
}

func (s *service) RebuildAllUserScores(ctx context.Context) error {
	return s.repo.RebuildAllUserScores(ctx)
}
