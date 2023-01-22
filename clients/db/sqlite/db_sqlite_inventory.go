package sqlite

import (
	"github.com/kampanosg/go-lsi/transformers"
	"github.com/kampanosg/go-lsi/types"
)

func (db SqliteDb) GetInventory() (inventory []types.InventoryItem, err error) {
	rows, err := db.Connection.Query(query_GET_INVENTORY)
	defer rows.Close()

	for rows.Next() {
		var id, squareId, categoryName, title, barcode, sku string
		var price float64
		if rows.Scan(&id, &squareId, &title, &categoryName, &price, &barcode, &sku); err != nil {
			return inventory, err
		}
		inventory_item := transformers.FromInventoryDbRowToDomain(id, squareId, title, categoryName, barcode, sku, price)
		inventory = append(inventory, inventory_item)
	}

	return inventory, nil
}

func (db SqliteDb) GetCategories() (categories []types.Category, err error) {
	rows, err := db.Connection.Query(query_GET_CATEGORIES)
	defer rows.Close()

	for rows.Next() {
		var id, squareId, name string
		var version int64
		if rows.Scan(&id, &squareId, &name, &version); err != nil {
			return categories, err
		}
		category := transformers.FromCategoryDbRowToDomain(id, squareId, name, version)
		categories = append(categories, category)
	}

	return categories, nil
}

func (db SqliteDb) InsertCategories(categories []types.Category) error {
	args := make([][]any, len(categories))
	for index, category := range categories {
		args[index] = []any{category.Id, category.SquareId, category.Name, category.Version}
	}
	return db.commitTx(query_INSERT_CATEGORY, args)
}

func (db SqliteDb) ClearCategories() error {
	return db.commitTx(query_CLEAR_CATEGORIES, make([][]any, 0))
}

func (db SqliteDb) GetProducts() (products []types.Product, err error) {
	rows, err := db.Connection.Query(query_GET_PRODUCTS)
	defer rows.Close()

	for rows.Next() {
		var id, squareId, squareVarId, categoryId, squareCategoryId, title, barcode, sku string
		var price float64
		var version int64
		if rows.Scan(&id, &squareId, &squareVarId, &categoryId, &squareCategoryId, &title, &price, &barcode, &sku, &version); err != nil {
			return products, err
		}
		product := transformers.FromProductDbRowToDomain(id, squareId, squareVarId, categoryId, squareCategoryId, title, barcode, sku, price, version)
		products = append(products, product)
	}

	return products, nil
}

func (db SqliteDb) InsertProducts(products []types.Product) error {
	args := make([][]any, len(products))
	for index, product := range products {
		args[index] = []any{product.Id, product.SquareId, product.SquareVarId,
			product.CategoryId, product.SquareCategoryId, product.Title, product.Price,
			product.Barcode, product.SKU, product.Version,
		}
	}
	return db.commitTx(query_INSERT_PRODUCT, args)
}

func (db SqliteDb) ClearProducts() error {
	return db.commitTx(query_CLEAR_PRODUCTS, make([][]any, 0))
}
