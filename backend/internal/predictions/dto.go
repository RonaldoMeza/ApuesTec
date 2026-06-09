package predictions

import "errors"

var (
	ErrPredictionNotFound  = errors.New("prediction not found")
	ErrPredictionExists    = errors.New("prediction already exists")
	ErrMatchNotEditable    = errors.New("match is not editable for predictions")
	ErrMatchLocked         = errors.New("match is locked")
	ErrMatchFinished       = errors.New("match is already finished")
	ErrMatchCancelled      = errors.New("match is cancelled")
	ErrTooCloseToStart     = errors.New("too close to match start time")
	ErrPredictionLocked    = errors.New("prediction is locked and cannot be edited")
	ErrNegativeScore       = errors.New("scores cannot be negative")
)

type CreatePredictionRequest struct {
	HomeScorePredicted int `json:"homeScorePredicted" binding:"required,min=0"`
	AwayScorePredicted int `json:"awayScorePredicted" binding:"required,min=0"`
}

type PredictionResponse struct {
	ID                 string `json:"id"`
	MatchID            string `json:"matchId"`
	HomeScorePredicted int    `json:"homeScorePredicted"`
	AwayScorePredicted int    `json:"awayScorePredicted"`
	Points             int    `json:"points"`
	IsLocked           bool   `json:"isLocked"`
	CanEdit            bool   `json:"canEdit"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

func toPredictionResponse(p *Prediction, canEdit bool) PredictionResponse {
	return PredictionResponse{
		ID:                 p.ID,
		MatchID:            p.MatchID,
		HomeScorePredicted: p.PredictedHomeScore,
		AwayScorePredicted: p.PredictedAwayScore,
		Points:             p.TotalPoints,
		IsLocked:           p.IsLocked(),
		CanEdit:            canEdit,
		CreatedAt:          p.CreatedAt.Format(timeFormat),
		UpdatedAt:          p.UpdatedAt.Format(timeFormat),
	}
}

const timeFormat = "2006-01-02T15:04:05Z"
