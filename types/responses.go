package types

import "time"

type ErrorResp struct {
	Message   string
	Timestamp time.Time
}

type AuthResponse struct {
	Message   string
	Token     string
	Timestamp time.Time
}

type InventoryItemResponse struct {
	Id           string  `json:"linnworksId"`
	SquareId     string  `json:"squareId"`
	Title        string  `json:"title"`
	CategoryName string  `json:"categoryName"`
	Barcode      string  `json:"barcode"`
	SKU          string  `json:"sku"`
	Price        float64 `json:"price"`
}
type InventoryResponse struct {
	Total int                     `json:"total"`
	Items []InventoryItemResponse `json:"items"`
}
