package teams

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Service interface {
	Create(ctx context.Context, req CreateTeamRequest) (*TeamResponse, error)
	List(ctx context.Context) ([]TeamResponse, error)
	GetByID(ctx context.Context, id string) (*TeamResponse, error)
	Update(ctx context.Context, id string, req UpdateTeamRequest) (*TeamResponse, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateTeamRequest) (*TeamResponse, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	team, err := s.repo.Create(ctx, req.Name, code, req.Country, req.FlagURL)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, &ServiceError{Status: http.StatusConflict, Code: "TEAM_EXISTS", Message: ErrTeamExists.Error()}
		}
		return nil, fmt.Errorf("create team: %w", err)
	}

	resp := toTeamResponse(team)
	return &resp, nil
}

func (s *service) List(ctx context.Context) ([]TeamResponse, error) {
	teams, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}

	resp := make([]TeamResponse, len(teams))
	for i, t := range teams {
		resp[i] = toTeamResponse(&t)
	}
	return resp, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*TeamResponse, error) {
	team, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "TEAM_NOT_FOUND", Message: ErrTeamNotFound.Error()}
	}

	resp := toTeamResponse(team)
	return &resp, nil
}

func (s *service) Update(ctx context.Context, id string, req UpdateTeamRequest) (*TeamResponse, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	team, err := s.repo.Update(ctx, id, req.Name, code, req.Country, req.FlagURL)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, &ServiceError{Status: http.StatusConflict, Code: "TEAM_EXISTS", Message: ErrTeamExists.Error()}
		}
		return nil, fmt.Errorf("update team: %w", err)
	}

	resp := toTeamResponse(team)
	return &resp, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return &ServiceError{Status: http.StatusNotFound, Code: "TEAM_NOT_FOUND", Message: ErrTeamNotFound.Error()}
	}
	return nil
}

type ServiceError struct {
	Status  int
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
