export interface TeamInfo {
  id: string;
  name: string;
  code: string;
  country: string;
  flagUrl?: string | null;
}

export interface MatchResponse {
  id: string;
  homeTeam?: TeamInfo | null;
  awayTeam?: TeamInfo | null;
  homeTeamId: string;
  awayTeamId: string;
  startsAt: string;
  status: string;
  homeScore?: number | null;
  awayScore?: number | null;
  lockedAt?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface CreateMatchRequest {
  homeTeamId: string;
  awayTeamId: string;
  startsAt: string;
}

export interface UpdateMatchRequest {
  homeTeamId: string;
  awayTeamId: string;
  startsAt: string;
}

export interface UpdateStatusRequest {
  status: string;
}

export interface UpdateResultRequest {
  homeScore: number;
  awayScore: number;
}
