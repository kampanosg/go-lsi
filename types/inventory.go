package types

type Category struct {
	Id       string `json:"id"`
	SquareId string `json:"squareId"`
	Name     string `json:"name"`
	Version  int64  `json:"version"`
}

type Product struct {
	Id               string         `json:"id"`
	SquareId         string         `json:"squareId"`
	SquareVarId      string         `json:"squareVariationId"`
	CategoryId       string         `json:"categoryId"`
	SquareCategoryId string         `json:"squareCategoryId"`
	Title            string         `json:"title"`
	Price            float64        `json:"price"`
	Barcode          string         `json:"barcode"`
	SKU              string         `json:"sku"`
	Images           []ProductImage `json:"images"`
	Version          int64          `json:"version"`
}

type ProductImage struct {
	Id        string `json:"id"`
	Source    string `json:"source"`
	Thumbnail string `json:"thumbnail"`
	Order     int    `json:"order"`
}

type InventoryItem struct {
	Id           string
	SquareId     string
	Title        string
	CategoryName string
	Barcode      string
	SKU          string
	Price        float64
}
