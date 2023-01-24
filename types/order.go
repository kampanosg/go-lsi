package types

import "time"

type Order struct {
	Id                 int
	LinnworksId        string
	SquareId           string
	LocationId         string
	State              string
	Version            int
	TotalMoney         float64
	TotalDiscount      float64
	TotalTip           float64
	TotalServiceCharge float64
	CreatedAt          time.Time
	Products           []OrderProduct
}

type OrderProduct struct {
	Id            int
	OrderId       int
	SquareOrderId string
	SquareVarId   string
	Quantity      string
	ItemNumber    string
	SKU           string
	Title         string
	PricePerUnit  float64
}
