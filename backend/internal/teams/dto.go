package teams

import "errors"

var (
	ErrTeamNotFound = errors.New("team not found")
	ErrTeamExists   = errors.New("team code already exists")
)

type CreateTeamRequest struct {
	Name    string  `json:"name" binding:"required,min=1,max=150"`
	Code    string  `json:"code" binding:"required,min=2,max=10"`
	Country string  `json:"country" binding:"required,min=2,max=100"`
	FlagURL *string `json:"flagUrl,omitempty"`
}

type UpdateTeamRequest struct {
	Name    string  `json:"name" binding:"required,min=1,max=150"`
	Code    string  `json:"code" binding:"required,min=2,max=10"`
	Country string  `json:"country" binding:"required,min=2,max=100"`
	FlagURL *string `json:"flagUrl,omitempty"`
}

type TeamResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	Country   string  `json:"country"`
	FlagURL   *string `json:"flagUrl,omitempty"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

func toTeamResponse(t *Team) TeamResponse {
	return TeamResponse{
		ID:        t.ID,
		Name:      t.Name,
		Code:      t.Code,
		Country:   t.Country,
		FlagURL:   t.FlagURL,
		CreatedAt: t.CreatedAt.Format(timeFormat),
		UpdatedAt: t.UpdatedAt.Format(timeFormat),
	}
}

const timeFormat = "2006-01-02T15:04:05Z"
