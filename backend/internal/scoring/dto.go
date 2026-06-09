package scoring

import "errors"

var (
	ErrMatchNotFinished  = errors.New("match is not finished")
	ErrMatchCancelled    = errors.New("match is cancelled, cannot award points")
	ErrScoreNotSet       = errors.New("match does not have a score set")
)

type ScoreEventResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	MatchID     string `json:"matchId"`
	EventType   string `json:"eventType"`
	Points      int    `json:"points"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt"`
}

func toScoreEventResponse(e *ScoreEvent) ScoreEventResponse {
	desc := ""
	if e.Description != nil {
		desc = *e.Description
	}
	return ScoreEventResponse{
		ID:          e.ID,
		UserID:      e.UserID,
		MatchID:     e.MatchID,
		EventType:   e.EventType,
		Points:      e.Points,
		Description: desc,
		CreatedAt:   e.CreatedAt.Format(timeFormat),
	}
}

func toScoreEventResponses(events []ScoreEvent) []ScoreEventResponse {
	resp := make([]ScoreEventResponse, len(events))
	for i, e := range events {
		resp[i] = toScoreEventResponse(&e)
	}
	return resp
}

const timeFormat = "2006-01-02T15:04:05Z"
