package predictions

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/audit"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/matches"
)

type MatchRepository interface {
	FindByID(ctx context.Context, id string) (*matches.Match, error)
}

type Service interface {
	GetMyPrediction(ctx context.Context, userID, matchID string) (*PredictionResponse, error)
	ListMyPredictions(ctx context.Context, userID string) ([]PredictionResponse, error)
	UpsertPrediction(ctx context.Context, userID, matchID string, req CreatePredictionRequest, ipAddress, userAgent *string) (*PredictionResponse, error)
	DeletePrediction(ctx context.Context, userID, matchID string) error
}

type service struct {
	repo      Repository
	matchRepo MatchRepository
	auditRepo audit.Repository
}

func NewService(repo Repository, matchRepo MatchRepository, auditRepo audit.Repository) Service {
	return &service{repo: repo, matchRepo: matchRepo, auditRepo: auditRepo}
}

func (s *service) GetMyPrediction(ctx context.Context, userID, matchID string) (*PredictionResponse, error) {
	prediction, err := s.repo.FindByUserAndMatch(ctx, userID, matchID)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "PREDICTION_NOT_FOUND", Message: ErrPredictionNotFound.Error()}
	}

	match, err := s.matchRepo.FindByID(ctx, matchID)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "MATCH_NOT_FOUND", Message: "match not found"}
	}

	canEdit := s.canEdit(match, prediction)
	resp := toPredictionResponse(prediction, canEdit)
	return &resp, nil
}

func (s *service) ListMyPredictions(ctx context.Context, userID string) ([]PredictionResponse, error) {
	predictions, err := s.repo.FindByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list my predictions: %w", err)
	}

	resp := make([]PredictionResponse, len(predictions))
	for i, p := range predictions {
		match, matchErr := s.matchRepo.FindByID(ctx, p.MatchID)
		canEdit := false
		if matchErr == nil {
			canEdit = s.canEdit(match, &p)
		}
		resp[i] = toPredictionResponse(&p, canEdit)
	}
	return resp, nil
}

func (s *service) UpsertPrediction(ctx context.Context, userID, matchID string, req CreatePredictionRequest, ipAddress, userAgent *string) (*PredictionResponse, error) {
	if req.HomeScorePredicted < 0 || req.AwayScorePredicted < 0 {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "NEGATIVE_SCORE", Message: ErrNegativeScore.Error()}
	}

	match, err := s.matchRepo.FindByID(ctx, matchID)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "MATCH_NOT_FOUND", Message: "match not found"}
	}

	if err := s.validateMatchForPrediction(match); err != nil {
		return nil, err
	}

	existingPrediction, findErr := s.repo.FindByUserAndMatch(ctx, userID, matchID)
	if findErr == nil && existingPrediction != nil {
		if existingPrediction.IsLocked() {
			return nil, &ServiceError{Status: http.StatusForbidden, Code: "PREDICTION_LOCKED", Message: ErrPredictionLocked.Error()}
		}

		updated, err := s.repo.Update(ctx, existingPrediction.ID, req.HomeScorePredicted, req.AwayScorePredicted)
		if err != nil {
			return nil, fmt.Errorf("update prediction: %w", err)
		}

		_ = s.auditRepo.Log(ctx, &userID, audit.ActionPredictionUpdated, "prediction", updated.ID, ipAddress, userAgent)

		canEdit := s.canEdit(match, updated)
		resp := toPredictionResponse(updated, canEdit)
		return &resp, nil
	}

	created, err := s.repo.Create(ctx, userID, matchID, req.HomeScorePredicted, req.AwayScorePredicted)
	if err != nil {
		return nil, fmt.Errorf("create prediction: %w", err)
	}

	_ = s.auditRepo.Log(ctx, &userID, audit.ActionPredictionCreated, "prediction", created.ID, ipAddress, userAgent)

	canEdit := s.canEdit(match, created)
	resp := toPredictionResponse(created, canEdit)
	return &resp, nil
}

func (s *service) DeletePrediction(ctx context.Context, userID, matchID string) error {
	prediction, err := s.repo.FindByUserAndMatch(ctx, userID, matchID)
	if err != nil {
		return &ServiceError{Status: http.StatusNotFound, Code: "PREDICTION_NOT_FOUND", Message: ErrPredictionNotFound.Error()}
	}

	if prediction.IsLocked() {
		return &ServiceError{Status: http.StatusForbidden, Code: "PREDICTION_LOCKED", Message: ErrPredictionLocked.Error()}
	}

	match, err := s.matchRepo.FindByID(ctx, matchID)
	if err != nil {
		return &ServiceError{Status: http.StatusNotFound, Code: "MATCH_NOT_FOUND", Message: "match not found"}
	}

	if err := s.validateMatchForPrediction(match); err != nil {
		return err
	}

	return s.repo.Delete(ctx, prediction.ID)
}

func (s *service) validateMatchForPrediction(match *matches.Match) *ServiceError {
	switch match.Status {
	case matches.StatusFinished:
		return &ServiceError{Status: http.StatusForbidden, Code: "MATCH_FINISHED", Message: ErrMatchFinished.Error()}
	case matches.StatusCancelled:
		return &ServiceError{Status: http.StatusForbidden, Code: "MATCH_CANCELLED", Message: ErrMatchCancelled.Error()}
	case matches.StatusLocked:
		return &ServiceError{Status: http.StatusForbidden, Code: "MATCH_LOCKED", Message: ErrMatchLocked.Error()}
	}

	if time.Until(match.StartsAt) <= 10*time.Minute {
		return &ServiceError{Status: http.StatusForbidden, Code: "TOO_CLOSE_TO_START", Message: ErrTooCloseToStart.Error()}
	}

	return nil
}

func (s *service) canEdit(match *matches.Match, prediction *Prediction) bool {
	if prediction.IsLocked() {
		return false
	}

	if match.Status == matches.StatusFinished || match.Status == matches.StatusCancelled {
		return false
	}

	if match.Status == matches.StatusLocked {
		return false
	}

	if time.Until(match.StartsAt) <= 10*time.Minute {
		return false
	}

	return true
}

type ServiceError struct {
	Status  int
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
