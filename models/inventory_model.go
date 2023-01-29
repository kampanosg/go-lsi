package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	LinnworksID string
	SquareID    string
	Name        string
	Version     int64
	Products    []Product
}

type Product struct {
	gorm.Model
	LinnworksID      string `gorm:"index"`
	SquareId         string
	SquareVarID      string `gorm:"index"`
	CategoryID       uint
	SquareCategoryID string
	Title            string
	Barcode          string `gorm:"index"`
	Price            float64
	SKU              string `gorm:"index"`
	Version          int64
}
