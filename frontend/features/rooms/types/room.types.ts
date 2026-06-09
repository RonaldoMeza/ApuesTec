export interface Room {
  id: string
  name: string
  description: string
  ownerId: string
  status: 'ACTIVE' | 'CLOSED'
  visibility: 'PUBLIC' | 'PRIVATE'
  hasPassword: boolean
  memberCount: number
  myRole: 'OWNER' | 'MODERATOR' | 'MEMBER'
  createdAt: string
  updatedAt: string
  closedAt?: string
}

export interface RoomListResponse {
  rooms: Room[]
  total: number
}

export interface CreateRoomRequest {
  name: string
  description?: string
  visibility?: 'PUBLIC' | 'PRIVATE'
  password?: string
}

export interface UpdateRoomRequest {
  name: string
  description?: string
  visibility?: 'PUBLIC' | 'PRIVATE'
  password?: string
}

export interface PublicRoom {
  id: string
  name: string
  description: string
  ownerName: string
  memberCount: number
  hasPassword: boolean
  isMember: boolean
}

export interface SearchPublicRoomsResponse {
  rooms: PublicRoom[]
  total: number
}

export interface JoinPublicRoomRequest {
  password?: string
}

export interface RoomMember {
  id: string
  roomId: string
  userId: string
  fullName: string
  email: string
  role: 'OWNER' | 'MODERATOR' | 'MEMBER'
  joinedAt: string
}

export interface MemberListResponse {
  members: RoomMember[]
  total: number
}

export interface ChangeRoleRequest {
  role: 'MEMBER' | 'MODERATOR'
}

export interface RoomLeaderboardEntry {
  rank: number
  userId: string
  fullName: string
  totalPoints: number
  predictionsCount: number
  exactScoresCount: number
  roomRole: string
}

export interface RoomLeaderboardResponse {
  entries: RoomLeaderboardEntry[]
  total: number
  roomId: string
  roomName: string
}

export interface RoomInvite {
  id: string
  roomId: string
  code: string
  qrPayload: string
  createdBy: string
  expiresAt: string
  usedAt?: string
  revokedAt?: string
  createdAt: string
}

export interface CreateInviteRequest {
  durationMinutes: number
}
