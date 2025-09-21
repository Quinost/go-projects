package models

type Status string

const (
	CreatedStatus Status = "Created"
	UpdatedStatus Status = "Updated"
	ErrorStatus Status = "Error"
)

type ItemCreateDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Response struct {
	Info   string `json:"info,omitempty"`
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}

func NewResponse(info string, status Status) *Response {
	return &Response{
		Info:   info,
		Status: status,
	}
}

func ErrorResponse(err string) *Response {
	return &Response{
		Status: ErrorStatus,
		Error: err,
	}
}
