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

type InventoryResponse struct {
    Id           string `json:"linnworksId"`
    SquareId     string `json:"squareId"`
    Title        string `json:"title"`
    CategoryName string `json:"categoryName"`
    Barcode      string `json:"barcode"`
    SKU          string `json:"sku"`
    Price        string `json:"price"`
}
