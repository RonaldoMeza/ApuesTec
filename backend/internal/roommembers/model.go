package roommembers

const (
	RoomRoleOwner      = "OWNER"
	RoomRoleModerator  = "MODERATOR"
	RoomRoleMember     = "MEMBER"
)

type ServiceError struct {
	Status  int
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
