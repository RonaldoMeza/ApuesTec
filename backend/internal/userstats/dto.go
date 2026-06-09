package userstats

type UserStatsResponse struct {
	TotalPoints                int    `json:"totalPoints"`
	PredictionsCount           int    `json:"predictionsCount"`
	ExactScoresCount           int    `json:"exactScoresCount"`
	WinnerCorrectCount         int    `json:"winnerCorrectCount"`
	GoalDifferenceCorrectCount int    `json:"goalDifferenceCorrectCount"`
	CurrentStreak              int    `json:"currentStreak"`
	BestStreak                 int    `json:"bestStreak"`
	GlobalRank                 int    `json:"globalRank"`
}

type ScoreEventResponse struct {
	ID          string `json:"id"`
	MatchID     string `json:"matchId"`
	EventType   string `json:"eventType"`
	Points      int    `json:"points"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt"`
}
