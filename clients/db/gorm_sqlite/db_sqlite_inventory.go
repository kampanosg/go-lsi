package gormsqlite

import (
	"github.com/kampanosg/go-lsi/models"
	"github.com/kampanosg/go-lsi/types"
)

func (db SqliteDb) GetCategories() ([]types.Category, error) {
	categories := make([]models.Category, 0)
	err := db.Connection.Find(&categories).Error
	return fromCategoryModelsToTypes(categories), err
}

func (db SqliteDb) InsertCategories(categories []types.Category) error {
	categoryModels := fromCategoryTypeToModels(categories)
	return db.Connection.Create(&categoryModels).Error
}

func (db SqliteDb) UpsertCategory(category types.Category) error {
	var existingCategory models.Category
	if err := db.Connection.Where("square_id = ?", category.SquareID).First(&existingCategory).Error; err == nil {
		db.Connection.Unscoped().Delete(&existingCategory)
	}
	upsertCategory := fromCategoryTypeToModel(category)
	return db.Connection.Save(&upsertCategory).Error
}

func (db SqliteDb) DeleteCategoriesBySquareIds(squareIds []string) error {
	return db.Connection.Where("square_id IN ?", squareIds).Unscoped().Delete(&models.Category{}).Error
}

func (db SqliteDb) GetProducts() ([]types.Product, error) {
	products := make([]models.Product, 0)
	result := db.Connection.Find(&products)
	if result.Error != nil {
		return []types.Product{}, result.Error
	}
	return fromProductModelsToTypes(products), nil
}

func (db SqliteDb) GetProductBySku(sku string) (types.Product, error) {
	var result models.Product
	err := db.Connection.Where("sku = ?", sku).First(&result).Error
	return transformResult(result), err
}

func (db SqliteDb) GetProductByBarcode(barcode string) (types.Product, error) {
	var result models.Product
	err := db.Connection.Where("barcode = ?", barcode).First(&result).Error
	return transformResult(result), err
}

func (db SqliteDb) GetProductByVarId(varId string) (types.Product, error) {
	var result models.Product
	err := db.Connection.Where("square_var_id = ?", varId).First(&result).Error
	return transformResult(result), err
}

func (db SqliteDb) GetProductByTitle(title string) (types.Product, error) {
	var result models.Product
	err := db.Connection.Where("title = ?", title).First(&result).Error
	return transformResult(result), err
}

func (db SqliteDb) UpsertProduct(product types.Product) error {
	var existingProduct models.Product
	if err := db.Connection.Where("square_id = ?", product.SquareID).First(&existingProduct).Error; err == nil {
		db.Connection.Unscoped().Delete(&existingProduct)
	}
	upsertProduct := fromProductTypeToModel(product)
	return db.Connection.Save(&upsertProduct).Error
}

func (db SqliteDb) DeleteProductsBySquareIds(squareIds []string) error {
	return db.Connection.Where("square_id IN ?", squareIds).Unscoped().Delete(&models.Product{}).Error
}

func transformResult(result models.Product) types.Product {
	return fromProductModelToType(result)
}

func fromCategoryModelsToTypes(categoryModels []models.Category) []types.Category {
	categories := make([]types.Category, len(categoryModels))
	for index, categoryModel := range categoryModels {
		category := types.Category{
			ID:          categoryModel.ID,
			LinnworksID: categoryModel.LinnworksID,
			SquareID:    categoryModel.SquareID,
			Name:        categoryModel.Name,
			Version:     categoryModel.Version,
		}
		categories[index] = category
	}
	return categories
}

func fromCategoryTypeToModels(categories []types.Category) []models.Category {
	categoryModels := make([]models.Category, len(categories))
	for index, category := range categories {
		categoryModels[index] = fromCategoryTypeToModel(category)
	}
	return categoryModels
}

func fromCategoryTypeToModel(category types.Category) models.Category {
	return models.Category{
		LinnworksID: category.LinnworksID,
		SquareID:    category.SquareID,
		Name:        category.Name,
		Version:     category.Version,
	}
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
		ID:                  productModel.ID,
		LinnworksID:         productModel.LinnworksID,
		LinnworksCategoryID: productModel.LinnworksCategoryId,
		SquareID:            productModel.SquareID,
		SquareVarID:         productModel.SquareVarID,
		SquareCategoryID:    productModel.SquareCategoryID,
		CategoryID:          productModel.CategoryID,
		Title:               productModel.Title,
		Price:               productModel.Price,
		Barcode:             productModel.Barcode,
		SKU:                 productModel.SKU,
		Version:             productModel.Version,
		VariationVersion:    productModel.VariationVersion,
		UpdatedAt:           productModel.UpdatedAt,
	}
}

func fromProductTypeToModel(product types.Product) models.Product {
	return models.Product{
		LinnworksID:         product.LinnworksID,
		LinnworksCategoryId: product.LinnworksCategoryID,
		SquareID:            product.SquareID,
		SquareVarID:         product.SquareVarID,
		SquareCategoryID:    product.SquareCategoryID,
		CategoryID:          product.CategoryID,
		Title:               product.Title,
		Price:               product.Price,
		Barcode:             product.Barcode,
		SKU:                 product.SKU,
		Version:             product.Version,
		VariationVersion:    product.VariationVersion,
	}
}
