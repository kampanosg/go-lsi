package types

import "time"

type ErrorResp struct {
	Message   string
	Timestamp time.Time
}

type AuthResp struct {
	Message   string
	Token     string
	Timestamp time.Time
}

type OkResp struct {
	Total int `json:"total"`
	Items any `json:"items"`
}
