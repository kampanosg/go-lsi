package transformers

import (
	"github.com/kampanosg/go-lsi/clients/linnworks"
	"github.com/kampanosg/go-lsi/types"
)

func FromCategoryDbRowToDomain(id, squareId, name string, version int64) types.Category {
	return types.Category{
		Id:       id,
		SquareId: squareId,
		Name:     name,
		Version:  version,
	}
}

func FromProductDbRowToDomain(id, squareId, squareVarId, categoryId, squareCategoryId, title, barcode, sku string, price float64, version int64) types.Product {
	return types.Product{
		Id:               id,
		SquareId:         squareId,
		SquareVarId:      squareVarId,
		CategoryId:       categoryId,
		SquareCategoryId: squareCategoryId,
		Title:            title,
		Barcode:          barcode,
		SKU:              sku,
		Price:            price,
		Version:          version,
	}
}

func FromCategoryLinnworksResponsesToDomain(lwCategories []linnworks.LinnworksCategoryResponse) (categories []types.Category) {
	for _, lwCategory := range lwCategories {
		categories = append(categories, FromCategoryLinnworksResponseToDomain(lwCategory))
	}
	return categories
}

func FromCategoryLinnworksResponseToDomain(lwCategory linnworks.LinnworksCategoryResponse) types.Category {
	return types.Category{
		Id:   lwCategory.Id,
		Name: lwCategory.Name,
	}
}

func FromProductLinnworksResponsesToDomain(lwProducts []linnworks.LinnworksProductResponse) (products []types.Product) {
	for _, p := range lwProducts {
		products = append(products, FromProductLinnworksResponseToDomain(p))
	}
	return products
}

func FromProductLinnworksResponseToDomain(lwProduct linnworks.LinnworksProductResponse) types.Product {
	return types.Product{
		Id:         lwProduct.StockItemID,
		CategoryId: lwProduct.CategoryID,
		Title:      lwProduct.ItemTitle,
		Barcode:    lwProduct.BarcodeNumber,
		Price:      lwProduct.RetailPrice,
		SKU:        lwProduct.SKU,
		Images:     FromProductLinnworksImageResponsesToDomain(lwProduct.Images),
	}
}

func FromProductLinnworksImageResponseToDomain(lwImage linnworks.LinnworksProductImageResponse) types.ProductImage {
	return types.ProductImage{
		Id:        lwImage.PkRowID,
		Source:    lwImage.FullSource,
		Thumbnail: lwImage.Source,
		Order:     lwImage.SortOrder,
	}
}

func FromProductLinnworksImageResponsesToDomain(lwImages []linnworks.LinnworksProductImageResponse) (images []types.ProductImage) {
	for _, i := range lwImages {
		images = append(images, FromProductLinnworksImageResponseToDomain(i))
	}
	return images
}
