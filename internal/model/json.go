package model

type Response[T any] struct {
	Success bool   `json:"success"`
	Body    T      `json:"body,omitempty"`
	Code    int    `json:"code,omitempty"`
	Error   string `json:"error,omitempty"`
}
