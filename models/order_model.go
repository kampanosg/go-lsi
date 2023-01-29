package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	SquareID        string
	LocationID      string
	State           string
	Version         int64
	TotalMoney      float64
	CreatedAtSquare time.Time
}
