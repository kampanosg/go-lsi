package gormsqlite

import (
	"github.com/kampanosg/go-lsi/models"
	"github.com/kampanosg/go-lsi/types"
)

func (db SqliteDb) GetCategories() ([]types.Category, error) {
	categories := make([]models.Category, 0)
	result := db.Connection.Find(&categories)
	if result.Error != nil {
		return []types.Category{}, result.Error
	}
	return fromCategoryModelsToTypes(categories), nil
}

func (db SqliteDb) InsertCategories(categories []types.Category) error {
	categoryModels := fromCategoryTypeToModels(categories)
	return db.Connection.Create(categoryModels).Error
}

func (db SqliteDb) ClearCategories() error {
	return db.Connection.Delete(&models.Category{}).Error
}

func (db SqliteDb) GetProducts() ([]types.Product, error) {
	products := make([]models.Product, 0)
	result := db.Connection.Find(&products)
	if result.Error != nil {
		return []types.Product{}, result.Error
	}
	return fromProductModelsToTypes(products), nil
}

func (db SqliteDb) GetProductByVarId(varId string) (types.Product, error) {
	var result models.Product
	db.Connection.Where(&models.Product{SquareVarID: varId}).First(&result)
	if result.ID == 0 {
		return types.Product{}, errRecordNotFound
	}
	return fromProductModelToType(result), nil
}

func (db SqliteDb) InsertProducts(products []types.Product) error {
	productModels := fromProductTypesToModels(products)
	return db.Connection.Create(productModels).Error
}

func (db SqliteDb) ClearProducts() error {
	return db.Connection.Delete(&models.Product{}).Error
}

func fromCategoryModelsToTypes(categoryModels []models.Category) []types.Category {
	categories := make([]types.Category, len(categoryModels))
	for index, categoryModel := range categoryModels {
		category := types.Category{
			ID:       categoryModel.ID,
			SquareID: categoryModel.SquareID,
			Name:     categoryModel.Name,
			Version:  categoryModel.Version,
		}
		categories[index] = category
	}
	return categories
}

func fromCategoryTypeToModels(categories []types.Category) []models.Category {
	categoryModels := make([]models.Category, len(categories))
	for index, category := range categories {
		categoryModel := models.Category{
			SquareID: category.SquareID,
			Name:     category.Name,
			Version:  category.Version,
		}
		categoryModels[index] = categoryModel
	}
	return categoryModels
}

func fromProductModelsToTypes(productModels []models.Product) []types.Product {
	products := make([]types.Product, len(productModels))
	for index, productModel := range productModels {
		products[index] = fromProductModelToType(productModel)
	}
	return products
}

func fromProductModelToType(productModel models.Product) types.Product {
	return types.Product{
		ID:               productModel.ID,
		LinnworksID:      productModel.LinnworksID,
		SquareID:         productModel.SquareID,
		SquareVarID:      productModel.SquareVarID,
		SquareCategoryID: productModel.SquareCategoryID,
		Title:            productModel.Title,
		Price:            productModel.Price,
		Barcode:          productModel.Barcode,
		SKU:              productModel.SKU,
		Version:          productModel.Version,
	}
}

func fromProductTypesToModels(products []types.Product) []models.Product {
	modelProducts := make([]models.Product, len(products))
	for index, product := range products {
		modelProduct := models.Product{
			LinnworksID:      product.LinnworksID,
			SquareID:         product.SquareID,
			SquareVarID:      product.SquareVarID,
			SquareCategoryID: product.SquareCategoryID,
			Title:            product.Title,
			Price:            product.Price,
			Barcode:          product.Barcode,
			SKU:              product.SKU,
			Version:          product.Version,
		}
		modelProducts[index] = modelProduct
	}
	return modelProducts
}
