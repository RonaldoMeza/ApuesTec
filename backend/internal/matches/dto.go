package matches

import (
	"errors"
	"time"
)

var (
	ErrMatchNotFound    = errors.New("match not found")
	ErrSameTeam         = errors.New("home and away teams must be different")
	ErrInvalidStatus    = errors.New("invalid match status")
	ErrScoreNotSet      = errors.New("home and away scores are required")
	ErrMatchNotFinished = errors.New("match is not finished")
	ErrMatchFinished    = errors.New("match is already finished")
)

type CreateMatchRequest struct {
	HomeTeamID string `json:"homeTeamId" binding:"required,uuid"`
	AwayTeamID string `json:"awayTeamId" binding:"required,uuid"`
	StartsAt   string `json:"startsAt" binding:"required"`
}

type UpdateMatchRequest struct {
	HomeTeamID string `json:"homeTeamId" binding:"required,uuid"`
	AwayTeamID string `json:"awayTeamId" binding:"required,uuid"`
	StartsAt   string `json:"startsAt" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=SCHEDULED LOCKED FINISHED CANCELLED"`
}

type UpdateResultRequest struct {
	HomeScore int `json:"homeScore" binding:"required,gte=0"`
	AwayScore int `json:"awayScore" binding:"required,gte=0"`
}

type MatchResponse struct {
	ID        string       `json:"id"`
	HomeTeam  *TeamInfo    `json:"homeTeam,omitempty"`
	AwayTeam  *TeamInfo    `json:"awayTeam,omitempty"`
	HomeTeamID string      `json:"homeTeamId"`
	AwayTeamID string      `json:"awayTeamId"`
	StartsAt  string       `json:"startsAt"`
	Status    string       `json:"status"`
	HomeScore *int         `json:"homeScore,omitempty"`
	AwayScore *int         `json:"awayScore,omitempty"`
	LockedAt  *string      `json:"lockedAt,omitempty"`
	CreatedAt string       `json:"createdAt"`
	UpdatedAt string       `json:"updatedAt"`
}

type TeamInfo struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Code    string  `json:"code"`
	Country string  `json:"country"`
	FlagURL *string `json:"flagUrl,omitempty"`
}

func toTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(timeFormat)
	return &s
}

const timeFormat = "2006-01-02T15:04:05Z"
