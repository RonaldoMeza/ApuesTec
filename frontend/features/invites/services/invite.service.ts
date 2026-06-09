import { apiRequest } from "@/shared/services/api-client";
import type { InvitePreview, JoinResponse } from "@/features/invites/types/invite.types";

export const inviteService = {
  preview(code: string): Promise<InvitePreview> {
    return apiRequest<InvitePreview>(`/invites/${code}`, { skipAuth: true });
  },

  join(code: string): Promise<JoinResponse> {
    return apiRequest<JoinResponse>(`/invites/${code}/join`, { method: "POST" });
  },
};
