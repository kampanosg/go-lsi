package types

import "time"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SyncRequest struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}
