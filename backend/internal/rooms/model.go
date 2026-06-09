package rooms

import "time"

const (
	RoomStatusActive = "ACTIVE"
	RoomStatusClosed = "CLOSED"
)

const (
	RoomRoleOwner     = "OWNER"
	RoomRoleModerator = "MODERATOR"
	RoomRoleMember    = "MEMBER"
)

const (
	VisibilityPublic  = "PUBLIC"
	VisibilityPrivate = "PRIVATE"
)

type Room struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	OwnerID       string     `json:"ownerId"`
	Status        string     `json:"status"`
	Visibility    string     `json:"visibility"`
	PasswordHash  *string    `json:"-"`
	NetworkPrefix *string    `json:"networkPrefix,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	ClosedAt      *time.Time `json:"closedAt,omitempty"`
}

type RoomMember struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"roomId"`
	UserID    string    `json:"userId"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joinedAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ServiceError struct {
	Status  int
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
