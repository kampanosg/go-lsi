package types

import "time"

type ErrorResp struct {
	Message   string
	Timestamp time.Time
}

type AuthResponse struct {
	Message   string
	Token     string
	Timestamp time.Time
}

type HttpResponse struct {
	Total int `json:"total"`
	Items any `json:"items"`
}
