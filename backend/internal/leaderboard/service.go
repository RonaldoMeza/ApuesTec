package leaderboard

import (
	"context"
	"fmt"
)

const defaultLimit = 50

type Service interface {
	GetGlobalLeaderboard(ctx context.Context, limit int) (*GlobalLeaderboardResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetGlobalLeaderboard(ctx context.Context, limit int) (*GlobalLeaderboardResponse, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > 100 {
		limit = 100
	}

	entries, total, err := s.repo.GetGlobalLeaderboard(ctx, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("get global leaderboard: %w", err)
	}

	resp := toLeaderboardResponse(entries, total, limit)
	return &resp, nil
}
