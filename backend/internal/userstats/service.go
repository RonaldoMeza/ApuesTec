package userstats

import (
	"context"
	"fmt"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/leaderboard"
)

type Service interface {
	GetUserStats(ctx context.Context, userID string) (*UserStatsResponse, error)
	ListScoreEvents(ctx context.Context, userID string) ([]ScoreEventResponse, error)
}

type service struct {
	repo          Repository
	leaderboardRepo leaderboard.Repository
}

func NewService(repo Repository, leaderboardRepo leaderboard.Repository) Service {
	return &service{repo: repo, leaderboardRepo: leaderboardRepo}
}

func (s *service) GetUserStats(ctx context.Context, userID string) (*UserStatsResponse, error) {
	score, err := s.repo.GetUserScore(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user score: %w", err)
	}

	stats, err := s.repo.GetPredictionStats(ctx, userID)
	if err != nil {
		stats = &PredictionStats{}
	}

	rank, err := s.leaderboardRepo.GetUserRank(ctx, userID)
	if err != nil {
		rank = 0
	}

	return &UserStatsResponse{
		TotalPoints:                score.TotalPoints,
		PredictionsCount:           score.PredictionsCount,
		ExactScoresCount:           score.ExactScoresCount,
		WinnerCorrectCount:         score.WinnerCorrectCount,
		GoalDifferenceCorrectCount: score.GoalDifferenceCorrectCount,
		CurrentStreak:              stats.CurrentStreak,
		BestStreak:                 stats.BestStreak,
		GlobalRank:                 rank,
	}, nil
}

func (s *service) ListScoreEvents(ctx context.Context, userID string) ([]ScoreEventResponse, error) {
	events, err := s.repo.ListScoreEvents(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list score events: %w", err)
	}

	resp := make([]ScoreEventResponse, len(events))
	for i, e := range events {
		desc := ""
		if e.Description != nil {
			desc = *e.Description
		}
		resp[i] = ScoreEventResponse{
			ID:          e.ID,
			MatchID:     e.MatchID,
			EventType:   e.EventType,
			Points:      e.Points,
			Description: desc,
			CreatedAt:   e.CreatedAt,
		}
	}
	return resp, nil
}
