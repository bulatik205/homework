package models

type Task struct {
	ID        int64   `json:"id"`
	Subject   string  `json:"subject"`
	Task      string  `json:"task"`
	DateFrom  *string `json:"date_from,omitempty"`
	DateTo    *string `json:"date_to,omitempty"`
	Status    *string `json:"status,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
