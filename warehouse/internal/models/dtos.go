package models

type Status string

const (
	CreatedStatus Status = "Created"
	UpdatedStatus Status = "Updated"
	DeletedStatus Status = "Deleted"
	ErrorStatus   Status = "Error"
)

type ItemCreateDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LoginDto struct {
	Username string `json:"username" validator:"required"`
	Password string `json:"password" validator:"required"`
}

type Response struct {
	Info   string `json:"info,omitempty"`
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}
