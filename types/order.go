package types

import "time"

type Order struct {
	ID         uint      `json:"id"`
	SquareID   string    `json:"squareId"`
	LocationID string    `json:"locationId"`
	State      string    `json:"state"`
	Version    int64     `json:"version"`
	TotalMoney float64   `json:"totalMoney"`
	CreatedAt  time.Time `json:"createdAt"`
}
