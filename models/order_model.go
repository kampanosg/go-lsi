package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	SquareID        string `gorm:"uniqueIndex"`
	LocationID      string
	State           string
	Version         int64
	TotalMoney      float64
	CreatedAtSquare time.Time
}
