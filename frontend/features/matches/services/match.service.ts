import { apiRequest } from "@/shared/services/api-client";
import type {
  MatchResponse,
  CreateMatchRequest,
  UpdateStatusRequest,
  UpdateResultRequest,
} from "@/features/matches/types/match.types";

export const matchService = {
  list(): Promise<MatchResponse[]> {
    return apiRequest<MatchResponse[]>("/matches", { skipAuth: true });
  },

  listUpcoming(): Promise<MatchResponse[]> {
    return apiRequest<MatchResponse[]>("/matches/upcoming", { skipAuth: true });
  },

  listFinished(): Promise<MatchResponse[]> {
    return apiRequest<MatchResponse[]>("/matches/finished", { skipAuth: true });
  },

  getByID(id: string): Promise<MatchResponse> {
    return apiRequest<MatchResponse>(`/matches/${id}`, { skipAuth: true });
  },

  create(data: CreateMatchRequest): Promise<MatchResponse> {
    return apiRequest<MatchResponse>("/admin/matches", { method: "POST", body: data });
  },

  updateStatus(id: string, data: UpdateStatusRequest): Promise<MatchResponse> {
    return apiRequest<MatchResponse>(`/admin/matches/${id}/status`, { method: "PATCH", body: data });
  },

  updateResult(id: string, data: UpdateResultRequest): Promise<MatchResponse> {
    return apiRequest<MatchResponse>(`/admin/matches/${id}/result`, { method: "PATCH", body: data });
  },
};
