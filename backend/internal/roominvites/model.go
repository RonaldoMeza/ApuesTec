package roominvites

import "time"

type RoomInvite struct {
	ID        string     `json:"id"`
	RoomID    string     `json:"roomId"`
	Code      string     `json:"code"`
	QRPayload *string    `json:"qrPayload,omitempty"`
	CreatedBy string     `json:"createdBy"`
	ExpiresAt time.Time  `json:"expiresAt"`
	UsedAt    *time.Time `json:"usedAt,omitempty"`
	RevokedAt *time.Time `json:"revokedAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

var AllowedDurations = []int{1, 3, 5, 10, 15, 20}

const DefaultDuration = 5

type ServiceError struct {
	Status  int
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
