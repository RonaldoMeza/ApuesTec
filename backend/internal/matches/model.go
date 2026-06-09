package matches

import "time"

const (
	StatusScheduled = "SCHEDULED"
	StatusLocked    = "LOCKED"
	StatusFinished  = "FINISHED"
	StatusCancelled = "CANCELLED"
)

type Match struct {
	ID          string     `json:"id"`
	HomeTeamID  string     `json:"homeTeamId"`
	AwayTeamID  string     `json:"awayTeamId"`
	StartsAt    time.Time  `json:"startsAt"`
	Status      string     `json:"status"`
	HomeScore   *int       `json:"homeScore,omitempty"`
	AwayScore   *int       `json:"awayScore,omitempty"`
	LockedAt    *time.Time `json:"lockedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
