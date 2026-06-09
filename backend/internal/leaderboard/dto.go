package leaderboard

type GlobalLeaderboardResponse struct {
	Total int                 `json:"total"`
	Limit int                 `json:"limit"`
	Entries []LeaderboardEntryResponse `json:"entries"`
}

type LeaderboardEntryResponse struct {
	Rank             int    `json:"rank"`
	UserID           string `json:"userId"`
	FullName         string `json:"fullName"`
	Email            string `json:"email"`
	TotalPoints      int    `json:"totalPoints"`
	PredictionsCount int    `json:"predictionsCount"`
	ExactScores      int    `json:"exactScores"`
	WinnerCorrect    int    `json:"winnerCorrect"`
}

func toLeaderboardResponse(entries []LeaderboardEntry, total, limit int) GlobalLeaderboardResponse {
	resp := GlobalLeaderboardResponse{
		Total:   total,
		Limit:   limit,
		Entries: make([]LeaderboardEntryResponse, len(entries)),
	}
	for i, e := range entries {
		resp.Entries[i] = LeaderboardEntryResponse{
			Rank:             e.Rank,
			UserID:           e.UserID,
			FullName:         e.FullName,
			Email:            e.Email,
			TotalPoints:      e.TotalPoints,
			PredictionsCount: e.PredictionsCount,
			ExactScores:      e.ExactScores,
			WinnerCorrect:    e.WinnerCorrect,
		}
	}
	return resp
}
