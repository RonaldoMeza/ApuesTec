package users

type UserResponse struct {
	ID        string   `json:"id"`
	FullName  string   `json:"fullName"`
	Email     string   `json:"email"`
	AvatarURL *string  `json:"avatarUrl,omitempty"`
	Status    string   `json:"status"`
	Roles     []string `json:"roles"`
	CreatedAt string   `json:"createdAt"`
}
