package scoring

import "time"

const (
	EventTypeExactScore         = "EXACT_SCORE"
	EventTypeWinnerCorrect      = "WINNER_CORRECT"
	EventTypeGoalDifference     = "GOAL_DIFFERENCE_CORRECT"
	EventTypeEarlyBonus         = "EARLY_BONUS"
	EventTypeStreakBonus        = "STREAK_BONUS"
)

type ScoreEvent struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	MatchID      string    `json:"matchId"`
	PredictionID string    `json:"predictionId"`
	EventType    string    `json:"eventType"`
	Points       int       `json:"points"`
	Description  *string   `json:"description,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type UserScore struct {
	UserID                     string     `json:"userId"`
	TotalPoints                int        `json:"totalPoints"`
	PredictionsCount           int        `json:"predictionsCount"`
	ExactScoresCount           int        `json:"exactScoresCount"`
	WinnerCorrectCount         int        `json:"winnerCorrectCount"`
	GoalDifferenceCorrectCount int        `json:"goalDifferenceCorrectCount"`
	StreakCount                int        `json:"streakCount"`
	LastScoredAt               *time.Time `json:"lastScoredAt,omitempty"`
	UpdatedAt                  time.Time  `json:"updatedAt"`
}

type ScoredPrediction struct {
	ID                     string
	UserID                 string
	MatchID                string
	PredictedHomeScore     int
	PredictedAwayScore     int
	IsExactScore           bool
	IsWinnerCorrect        bool
	IsGoalDifferenceCorrect bool
	BasePoints             int
	EarlyBonusPoints       int
	StreakBonusPoints      int
	TotalPoints            int
	CreatedAt              time.Time
	LockedAt               *time.Time
}
