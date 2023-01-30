package types

import "time"

type Category struct {
	ID          uint   `json:"id"`
	LinnworksID string `json:"linnworksId"`
	SquareID    string `json:"squareId"`
	Name        string `json:"name"`
	Version     int64  `json:"version"`
}

type Product struct {
	ID                  uint      `json:"id"`
	LinnworksID         string    `json:"linnworksId"`
	SquareID            string    `json:"squareId"`
	SquareVarID         string    `json:"squareVariationId"`
	CategoryID          uint      `json:"categoryId"`
	LinnworksCategoryID string    `json:"linnworksCategoryId"`
	SquareCategoryID    string    `json:"squareCategoryId"`
	Title               string    `json:"title"`
	Price               float64   `json:"price"`
	Barcode             string    `json:"barcode"`
	SKU                 string    `json:"sku"`
	Version             int64     `json:"version"`
	UpdatedAt           time.Time `json:"updatedAt"`
}
