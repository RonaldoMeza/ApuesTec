import { apiRequest } from "@/shared/services/api-client";
import type { PredictionResponse, CreatePredictionRequest } from "@/features/predictions/types/prediction.types";

export const predictionService = {
  listMyPredictions(): Promise<PredictionResponse[]> {
    return apiRequest<PredictionResponse[]>("/predictions/me");
  },

  getMyPrediction(matchId: string): Promise<PredictionResponse> {
    return apiRequest<PredictionResponse>(`/matches/${matchId}/prediction`);
  },

  upsertPrediction(matchId: string, data: CreatePredictionRequest): Promise<PredictionResponse> {
    return apiRequest<PredictionResponse>(`/matches/${matchId}/prediction`, {
      method: "POST",
      body: data,
    });
  },

  deletePrediction(matchId: string): Promise<void> {
    return apiRequest<void>(`/matches/${matchId}/prediction`, {
      method: "DELETE",
    });
  },
};
