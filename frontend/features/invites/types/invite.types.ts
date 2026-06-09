export interface InvitePreview {
  code: string
  roomName: string
  roomDescription: string
  expiresAt: string
  isExpired: boolean
}

export interface JoinResponse {
  roomId: string
  roomName: string
  role: string
}
