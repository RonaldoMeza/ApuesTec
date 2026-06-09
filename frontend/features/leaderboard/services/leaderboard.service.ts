import { apiRequest } from "@/shared/services/api-client";
import type { GlobalLeaderboardResponse, UserStats, ScoreEvent } from "@/features/leaderboard/types/leaderboard.types";

export const leaderboardService = {
  getGlobalLeaderboard(limit = 50): Promise<GlobalLeaderboardResponse> {
    return apiRequest<GlobalLeaderboardResponse>(`/leaderboard/global?limit=${limit}`, { skipAuth: true });
  },

  getMyStats(): Promise<UserStats> {
    return apiRequest<UserStats>("/users/me/stats");
  },

  getMyScoreEvents(): Promise<ScoreEvent[]> {
    return apiRequest<ScoreEvent[]>("/users/me/score-events");
  },
};
