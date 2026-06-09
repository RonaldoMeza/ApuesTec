package predictions

import "time"

type Prediction struct {
	ID                   string     `json:"id"`
	UserID               string     `json:"userId"`
	MatchID              string     `json:"matchId"`
	PredictedHomeScore   int        `json:"predictedHomeScore"`
	PredictedAwayScore   int        `json:"predictedAwayScore"`
	IsExactScore         bool       `json:"isExactScore"`
	IsWinnerCorrect      bool       `json:"isWinnerCorrect"`
	IsGoalDifferenceCorrect bool    `json:"isGoalDifferenceCorrect"`
	BasePoints           int        `json:"basePoints"`
	EarlyBonusPoints     int        `json:"earlyBonusPoints"`
	StreakBonusPoints    int        `json:"streakBonusPoints"`
	TotalPoints          int        `json:"totalPoints"`
	LockedAt             *time.Time `json:"lockedAt,omitempty"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
}

func (p *Prediction) IsLocked() bool {
	return p.LockedAt != nil
}
