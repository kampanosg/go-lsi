package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	LinnworksID string `gorm:"uniqueIndex"`
	SquareID    string `gorm:"uniqueIndex"`
	Name        string
	Version     int64
	Products    []Product
}

type Product struct {
	gorm.Model
	LinnworksID         string `gorm:"uniqueIndex"`
	SquareID            string `gorm:"uniqueIndex"`
	SquareVarID         string `gorm:"uniqueIndex"`
	CategoryID          uint
	LinnworksCategoryId string
	SquareCategoryID    string
	Title               string
	Barcode             string `gorm:"index"`
	Price               float64
	SKU                 string `gorm:"index"`
	Version             int64
	VariationVersion    int64
}
