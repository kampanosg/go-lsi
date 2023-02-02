package db

import (
	"github.com/kampanosg/go-lsi/types"
)

type DB interface {
	GetCategories() ([]types.Category, error)
	InsertCategories(categories []types.Category) error
	UpsertCategory(category types.Category) error
	DeleteCategoriesBySquareIds(squareIds []string) error

	GetProducts() ([]types.Product, error)
	GetProductByBarcode(barcode string) (types.Product, error)
	GetProductBySku(sku string) (types.Product, error)
	GetProductByVarId(varId string) (types.Product, error)
	InsertProducts(products []types.Product) error
	InsertProduct(product types.Product) error
	ClearProducts() error

	GetUserByUsername(username string) (types.User, error)

	GetOrders() ([]types.Order, error)
	GetOrderBySquareId(squareId string) (types.Order, error)
	InsertOrders([]types.Order) error

	GetLastSyncStatus() (types.SyncStatus, error)
	InsertSyncStatus(ts int64) error
}
