import { apiRequest } from "@/shared/services/api-client";
import type { TeamResponse, CreateTeamRequest, UpdateTeamRequest } from "@/features/teams/types/team.types";

export const teamService = {
  list(): Promise<TeamResponse[]> {
    return apiRequest<TeamResponse[]>("/teams", { skipAuth: true });
  },

  getByID(id: string): Promise<TeamResponse> {
    return apiRequest<TeamResponse>(`/teams/${id}`, { skipAuth: true });
  },

  create(data: CreateTeamRequest): Promise<TeamResponse> {
    return apiRequest<TeamResponse>("/admin/teams", { method: "POST", body: data });
  },

  update(id: string, data: UpdateTeamRequest): Promise<TeamResponse> {
    return apiRequest<TeamResponse>(`/admin/teams/${id}`, { method: "PUT", body: data });
  },

  delete(id: string): Promise<void> {
    return apiRequest<void>(`/admin/teams/${id}`, { method: "DELETE" });
  },
};
