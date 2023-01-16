package transform

import (
	"kev/types/domain"
	"kev/types/response"
)

func FromProductRespToDomain(resp response.ProductResponse) domain.Product {
	return domain.Product{
		Id:         resp.StockItemID,
		CategoryId: resp.CategoryID,
		Title:      resp.ItemTitle,
		Barcode:    resp.BarcodeNumber,
		Price:      resp.RetailPrice,
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
