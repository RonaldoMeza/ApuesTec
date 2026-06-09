export interface LeaderboardEntry {
  rank: number;
  userId: string;
  fullName: string;
  email: string;
  totalPoints: number;
  predictionsCount: number;
  exactScores: number;
  winnerCorrect: number;
}

export interface GlobalLeaderboardResponse {
  total: number;
  limit: number;
  entries: LeaderboardEntry[];
}

export interface UserStats {
  totalPoints: number;
  predictionsCount: number;
  exactScoresCount: number;
  winnerCorrectCount: number;
  goalDifferenceCorrectCount: number;
  currentStreak: number;
  bestStreak: number;
  globalRank: number;
}

export interface ScoreEvent {
  id: string;
  matchId: string;
  eventType: string;
  points: number;
  description: string;
  createdAt: string;
}
