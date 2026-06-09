export interface TeamResponse {
  id: string;
  name: string;
  code: string;
  country: string;
  flagUrl?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface CreateTeamRequest {
  name: string;
  code: string;
  country: string;
  flagUrl?: string | null;
}

export interface UpdateTeamRequest {
  name: string;
  code: string;
  country: string;
  flagUrl?: string | null;
}
