package domain

type Product struct {
	Id         string         `json:"id"`
	CategoryId string         `json:"categoryId"`
	Title      string         `json:"title"`
	Price      float64        `json:"price"`
	Barcode    string         `json:"barcode"`
	Images     []ProductImage `json:"images"`
}

type ProductImage struct {
	Id        string `json:"id"`
	Source    string `json:"source"`
	Thumbnail string `json:"thumbnail"`
	Order     int    `json:"order"`
}
