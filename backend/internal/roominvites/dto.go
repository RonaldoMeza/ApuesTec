package roominvites

type CreateInviteRequest struct {
	DurationMinutes int `json:"durationMinutes" binding:"required"`
}

type InviteResponse struct {
	ID        string  `json:"id"`
	RoomID    string  `json:"roomId"`
	Code      string  `json:"code"`
	QRPayload string  `json:"qrPayload"`
	CreatedBy string  `json:"createdBy"`
	ExpiresAt string  `json:"expiresAt"`
	UsedAt    *string `json:"usedAt,omitempty"`
	RevokedAt *string `json:"revokedAt,omitempty"`
	CreatedAt string  `json:"createdAt"`
}

type InvitePreviewResponse struct {
	Code            string `json:"code"`
	RoomName        string `json:"roomName"`
	RoomDescription string `json:"roomDescription"`
	ExpiresAt       string `json:"expiresAt"`
	IsExpired       bool   `json:"isExpired"`
}

type JoinResponse struct {
	RoomID   string `json:"roomId"`
	RoomName string `json:"roomName"`
	Role     string `json:"role"`
}
