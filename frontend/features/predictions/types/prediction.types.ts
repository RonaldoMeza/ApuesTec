export interface PredictionResponse {
  id: string;
  matchId: string;
  homeScorePredicted: number;
  awayScorePredicted: number;
  points: number;
  isLocked: boolean;
  canEdit: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreatePredictionRequest {
  homeScorePredicted: number;
  awayScorePredicted: number;
}

export interface PredictionWithMatch extends PredictionResponse {
  match?: {
    id: string;
    homeTeam?: { name: string; code: string; country: string } | null;
    awayTeam?: { name: string; code: string; country: string } | null;
    startsAt: string;
    status: string;
    homeScore?: number | null;
    awayScore?: number | null;
  };
}
