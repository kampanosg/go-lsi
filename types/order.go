package types

import (
	"time"
)

type Order struct {
	ID         uint           `json:"id"`
	SquareID   string         `json:"squareId"`
	LocationID string         `json:"locationId"`
	State      string         `json:"state"`
	Version    int64          `json:"version"`
	TotalMoney float64        `json:"totalMoney"`
	CreatedAt  time.Time      `json:"createdAt"`
	Products   []OrderProduct `json:"products"`
}

type OrderProduct struct {
	Id            int     `json:"id"`
	OrderId       int     `json:"orderId"`
	SquareOrderId string  `json:"squareOrderId"`
	SquareVarId   string  `json:"squareVarId"`
	Quantity      string  `json:"qty"`
	ItemNumber    string  `json:"itemNumber"`
	SKU           string  `json:"sku"`
	Title         string  `json:"title"`
	PricePerUnit  float64 `json:"pricePerUnit"`
}
