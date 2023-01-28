package types

import "time"

type SyncStatus struct {
	Timestamp time.Time `json:"ts"`
}
