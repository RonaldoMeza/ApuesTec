package roommembers

type MemberResponse struct {
	ID        string `json:"id"`
	RoomID    string `json:"roomId"`
	UserID    string `json:"userId"`
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	JoinedAt  string `json:"joinedAt"`
}

type MemberListResponse struct {
	Members []MemberResponse `json:"members"`
	Total   int              `json:"total"`
}

type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=MEMBER MODERATOR"`
}
