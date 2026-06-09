package rooms

type CreateRoomRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=120"`
	Description *string `json:"description,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
	Password    *string `json:"password,omitempty" binding:"omitempty,min=4,max=50"`
}

type UpdateRoomRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=120"`
	Description *string `json:"description,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
	Password    *string `json:"password,omitempty" binding:"omitempty,min=4,max=50"`
}

type RoomResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"ownerId"`
	Status      string `json:"status"`
	Visibility  string `json:"visibility"`
	HasPassword bool   `json:"hasPassword"`
	MemberCount int    `json:"memberCount"`
	MyRole      string `json:"myRole"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	ClosedAt    string `json:"closedAt"`
}

type RoomListResponse struct {
	Rooms []RoomResponse `json:"rooms"`
	Total int            `json:"total"`
}

type PublicRoomResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerName   string `json:"ownerName"`
	MemberCount int    `json:"memberCount"`
	HasPassword bool   `json:"hasPassword"`
	IsMember    bool   `json:"isMember"`
}

type SearchPublicRoomsResponse struct {
	Rooms []PublicRoomResponse `json:"rooms"`
	Total int                  `json:"total"`
}

type JoinPublicRoomRequest struct {
	Password string `json:"password,omitempty"`
}

type RoomLeaderboardEntry struct {
	Rank               int    `json:"rank"`
	UserID             string `json:"userId"`
	FullName           string `json:"fullName"`
	TotalPoints        int    `json:"totalPoints"`
	PredictionsCount   int    `json:"predictionsCount"`
	ExactScoresCount   int    `json:"exactScoresCount"`
	RoomRole           string `json:"roomRole"`
}

type RoomLeaderboardResponse struct {
	Entries  []RoomLeaderboardEntry `json:"entries"`
	Total    int                    `json:"total"`
	RoomID   string                 `json:"roomId"`
	RoomName string                 `json:"roomName"`
}

const timeFormat = "2006-01-02T15:04:05Z"
