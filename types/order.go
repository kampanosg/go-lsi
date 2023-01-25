package types

import "time"

type Order struct {
	Id                 int            `json:"id"`
	LinnworksId        string         `json:"linnworksId"`
	SquareId           string         `json:"squareId"`
	LocationId         string         `json:"locationId"`
	State              string         `json:"state"`
	Version            int            `json:"version"`
	TotalMoney         float64        `json:"totalMoney"`
	TotalDiscount      float64        `json:"totalDiscount"`
	TotalTip           float64        `json:"totalTip"`
	TotalServiceCharge float64        `json:"totalServiceCharge"`
	CreatedAt          time.Time      `json:"createdAt"`
	Products           []OrderProduct `json:"products"`
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
