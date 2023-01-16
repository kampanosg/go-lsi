package transform

import (
	"github.com/kampanosg/go-lsi/types/domain"
	"github.com/kampanosg/go-lsi/types/response"
)

func FromProductRespToDomain(resp response.ProductResponse) domain.Product {
	var price float64
	var title string

	if resp.RetailPrice < 0 {
		price = 0
	} else {
		price = resp.RetailPrice
	}

	if resp.ItemTitle == "" {
		title = "---"
	} else {
		title = resp.ItemTitle
	}

	return domain.Product{
		Id:         resp.StockItemID,
		CategoryId: resp.CategoryID,
		Title:      title,
		Barcode:    resp.BarcodeNumber,
		Price:      price,
		SKU:        resp.SKU,
		Images:     FromProductImagesResponseToDomain(resp.Images),
	}
}

func FromProductsRespToDomain(resps []response.ProductResponse) (products []domain.Product) {
	for _, p := range resps {
		products = append(products, FromProductRespToDomain(p))
	}
	return products
}

func FromProductImageResponseToDomain(resp response.ProductImageResponse) domain.ProductImage {
	return domain.ProductImage{
		Id:        resp.PkRowID,
		Source:    resp.FullSource,
		Thumbnail: resp.Source,
		Order:     resp.SortOrder,
	}
}

func FromProductImagesResponseToDomain(resps []response.ProductImageResponse) (images []domain.ProductImage) {
	for _, i := range resps {
		images = append(images, FromProductImageResponseToDomain(i))
	}
	return images
}
