package matches

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type TeamRepository interface {
	FindByID(ctx context.Context, id string) (*TeamInfo, error)
}

type ScoringFunc func(ctx context.Context, match MatchInfo) error

type MatchInfo struct {
	ID        string
	HomeScore int
	AwayScore int
	Status    string
	StartsAt  time.Time
}

type Service interface {
	Create(ctx context.Context, req CreateMatchRequest) (*MatchResponse, error)
	List(ctx context.Context) ([]MatchResponse, error)
	ListUpcoming(ctx context.Context) ([]MatchResponse, error)
	ListFinished(ctx context.Context) ([]MatchResponse, error)
	GetByID(ctx context.Context, id string) (*MatchResponse, error)
	Update(ctx context.Context, id string, req UpdateMatchRequest) (*MatchResponse, error)
	UpdateStatus(ctx context.Context, id string, req UpdateStatusRequest) (*MatchResponse, error)
	UpdateResult(ctx context.Context, id string, req UpdateResultRequest) (*MatchResponse, error)
	RecalculateScore(ctx context.Context, id string) (*MatchResponse, error)
}

type service struct {
	repo       Repository
	teamRepo   TeamRepository
	scoringFn  ScoringFunc
}

func NewService(repo Repository, teamRepo TeamRepository) Service {
	return &service{repo: repo, teamRepo: teamRepo}
}

func NewServiceWithScoring(repo Repository, teamRepo TeamRepository, scoringFn ScoringFunc) Service {
	return &service{repo: repo, teamRepo: teamRepo, scoringFn: scoringFn}
}

func (s *service) enrichMatch(ctx context.Context, m *Match) *MatchResponse {
	resp := &MatchResponse{
		ID:         m.ID,
		HomeTeamID: m.HomeTeamID,
		AwayTeamID: m.AwayTeamID,
		StartsAt:   m.StartsAt.Format(timeFormat),
		Status:     m.Status,
		HomeScore:  m.HomeScore,
		AwayScore:  m.AwayScore,
		LockedAt:   toTimePtr(m.LockedAt),
		CreatedAt:  m.CreatedAt.Format(timeFormat),
		UpdatedAt:  m.UpdatedAt.Format(timeFormat),
	}

	if homeTeam, err := s.teamRepo.FindByID(ctx, m.HomeTeamID); err == nil {
		resp.HomeTeam = homeTeam
	}
	if awayTeam, err := s.teamRepo.FindByID(ctx, m.AwayTeamID); err == nil {
		resp.AwayTeam = awayTeam
	}

	return resp
}

func (s *service) Create(ctx context.Context, req CreateMatchRequest) (*MatchResponse, error) {
	if req.HomeTeamID == req.AwayTeamID {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "SAME_TEAM", Message: ErrSameTeam.Error()}
	}

	startsAt, err := time.Parse(timeFormat, req.StartsAt)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "INVALID_DATE", Message: "invalid date format, use ISO 8601"}
	}

	match, err := s.repo.Create(ctx, req.HomeTeamID, req.AwayTeamID, startsAt)
	if err != nil {
		return nil, fmt.Errorf("create match: %w", err)
	}

	return s.enrichMatch(ctx, match), nil
}

func (s *service) List(ctx context.Context) ([]MatchResponse, error) {
	matches, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list matches: %w", err)
	}

	resp := make([]MatchResponse, len(matches))
	for i, m := range matches {
		enriched := s.enrichMatch(ctx, &m)
		resp[i] = *enriched
	}
	return resp, nil
}

func (s *service) ListUpcoming(ctx context.Context) ([]MatchResponse, error) {
	matches, err := s.repo.FindUpcoming(ctx)
	if err != nil {
		return nil, fmt.Errorf("list upcoming matches: %w", err)
	}

	resp := make([]MatchResponse, len(matches))
	for i, m := range matches {
		enriched := s.enrichMatch(ctx, &m)
		resp[i] = *enriched
	}
	return resp, nil
}

func (s *service) ListFinished(ctx context.Context) ([]MatchResponse, error) {
	matches, err := s.repo.FindFinished(ctx)
	if err != nil {
		return nil, fmt.Errorf("list finished matches: %w", err)
	}

	resp := make([]MatchResponse, len(matches))
	for i, m := range matches {
		enriched := s.enrichMatch(ctx, &m)
		resp[i] = *enriched
	}
	return resp, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*MatchResponse, error) {
	match, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "MATCH_NOT_FOUND", Message: ErrMatchNotFound.Error()}
	}

	return s.enrichMatch(ctx, match), nil
}

func (s *service) Update(ctx context.Context, id string, req UpdateMatchRequest) (*MatchResponse, error) {
	if req.HomeTeamID == req.AwayTeamID {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "SAME_TEAM", Message: ErrSameTeam.Error()}
	}

	startsAt, err := time.Parse(timeFormat, req.StartsAt)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "INVALID_DATE", Message: "invalid date format, use ISO 8601"}
	}

	match, err := s.repo.Update(ctx, id, req.HomeTeamID, req.AwayTeamID, startsAt)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "MATCH_NOT_FOUND", Message: ErrMatchNotFound.Error()}
	}

	return s.enrichMatch(ctx, match), nil
}

func (s *service) UpdateStatus(ctx context.Context, id string, req UpdateStatusRequest) (*MatchResponse, error) {
	validStatuses := map[string]bool{
		StatusScheduled: true,
		StatusLocked:    true,
		StatusFinished:  true,
		StatusCancelled: true,
	}

	if !validStatuses[req.Status] {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "INVALID_STATUS", Message: ErrInvalidStatus.Error()}
	}

	match, err := s.repo.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "MATCH_NOT_FOUND", Message: ErrMatchNotFound.Error()}
	}

	return s.enrichMatch(ctx, match), nil
}

func (s *service) UpdateResult(ctx context.Context, id string, req UpdateResultRequest) (*MatchResponse, error) {
	match, err := s.repo.UpdateResult(ctx, id, req.HomeScore, req.AwayScore)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "MATCH_NOT_FOUND", Message: ErrMatchNotFound.Error()}
	}

	enriched := s.enrichMatch(ctx, match)

	if s.scoringFn != nil && match.Status == StatusFinished {
		scoringMatch := MatchInfo{
			ID:        match.ID,
			HomeScore: req.HomeScore,
			AwayScore: req.AwayScore,
			Status:    match.Status,
			StartsAt:  match.StartsAt,
		}
		_ = s.scoringFn(ctx, scoringMatch)
	}

	return enriched, nil
}

func (s *service) RecalculateScore(ctx context.Context, id string) (*MatchResponse, error) {
	match, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "MATCH_NOT_FOUND", Message: ErrMatchNotFound.Error()}
	}

	if match.Status != StatusFinished {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "MATCH_NOT_FINISHED", Message: ErrMatchNotFinished.Error()}
	}

	if match.HomeScore == nil || match.AwayScore == nil {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "SCORE_NOT_SET", Message: ErrScoreNotSet.Error()}
	}

	if s.scoringFn != nil {
		scoringMatch := MatchInfo{
			ID:        match.ID,
			HomeScore: *match.HomeScore,
			AwayScore: *match.AwayScore,
			Status:    match.Status,
			StartsAt:  match.StartsAt,
		}
		if err := s.scoringFn(ctx, scoringMatch); err != nil {
			return nil, &ServiceError{Status: http.StatusInternalServerError, Code: "SCORE_CALCULATION_ERROR", Message: "score calculation failed: " + err.Error()}
		}
	}

	return s.enrichMatch(ctx, match), nil
}

type ServiceError struct {
	Status  int
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
