import { apiRequest } from "@/shared/services/api-client";
import type {
  Room,
  RoomListResponse,
  CreateRoomRequest,
  UpdateRoomRequest,
  MemberListResponse,
  RoomLeaderboardResponse,
  RoomInvite,
  SearchPublicRoomsResponse,
  JoinPublicRoomRequest,
} from "@/features/rooms/types/room.types";

export const roomService = {
  listMyRooms(): Promise<RoomListResponse> {
    return apiRequest<RoomListResponse>("/rooms");
  },

  create(data: CreateRoomRequest): Promise<Room> {
    return apiRequest<Room>("/rooms", {
      method: "POST",
      body: data,
    });
  },

  getById(id: string): Promise<Room> {
    return apiRequest<Room>(`/rooms/${id}`);
  },

  update(id: string, data: UpdateRoomRequest): Promise<Room> {
    return apiRequest<Room>(`/rooms/${id}`, {
      method: "PUT",
      body: data,
    });
  },

  close(id: string): Promise<Room> {
    return apiRequest<Room>(`/rooms/${id}/close`, {
      method: "PATCH",
    });
  },

  searchPublic(q = ""): Promise<SearchPublicRoomsResponse> {
    return apiRequest<SearchPublicRoomsResponse>(
      `/rooms/public${q ? `?q=${encodeURIComponent(q)}` : ""}`
    );
  },

  joinPublic(id: string, data?: JoinPublicRoomRequest): Promise<Room> {
    return apiRequest<Room>(`/rooms/${id}/join`, {
      method: "POST",
      body: data || {},
    });
  },

  getMembers(roomId: string): Promise<MemberListResponse> {
    return apiRequest<MemberListResponse>(`/rooms/${roomId}/members`);
  },

  changeRole(roomId: string, userId: string, role: string): Promise<void> {
    return apiRequest<void>(`/rooms/${roomId}/members/${userId}/role`, {
      method: "PATCH",
      body: { role },
    });
  },

  removeMember(roomId: string, userId: string): Promise<void> {
    return apiRequest<void>(`/rooms/${roomId}/members/${userId}`, {
      method: "DELETE",
    });
  },

  leave(roomId: string): Promise<void> {
    return apiRequest<void>(`/rooms/${roomId}/leave`, {
      method: "POST",
    });
  },

  getLeaderboard(roomId: string): Promise<RoomLeaderboardResponse> {
    return apiRequest<RoomLeaderboardResponse>(`/rooms/${roomId}/leaderboard`);
  },

  createInvite(roomId: string, durationMinutes: number): Promise<RoomInvite> {
    return apiRequest<RoomInvite>(`/rooms/${roomId}/invites`, {
      method: "POST",
      body: { durationMinutes },
    });
  },
};
