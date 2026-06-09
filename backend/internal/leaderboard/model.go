package leaderboard

type LeaderboardEntry struct {
	UserID       string `json:"userId"`
	FullName     string `json:"fullName"`
	Email        string `json:"email"`
	TotalPoints  int    `json:"totalPoints"`
	PredictionsCount int `json:"predictionsCount"`
	ExactScores  int    `json:"exactScores"`
	WinnerCorrect int   `json:"winnerCorrect"`
	Rank         int    `json:"rank"`
}
