package models

import (
	"time"

	"gorm.io/gorm"
)

type SyncStatus struct {
    gorm.Model
    LastRun time.Time
    Success bool
}
