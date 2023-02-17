package db

import (
	"time"

	"github.com/kampanosg/go-lsi/types"
)

type DB interface {
	GetCategories() ([]types.Category, error)
	InsertCategories([]types.Category) error
	UpsertCategory(category types.Category) error
	DeleteCategoriesBySquareIds(squareIds []string) error

	GetProducts() ([]types.Product, error)
	GetProductByBarcode(barcode string) (types.Product, error)
	GetProductBySku(sku string) (types.Product, error)
	GetProductByVarId(varId string) (types.Product, error)
	GetProductByTitle(title string) (types.Product, error)
	UpsertProduct(product types.Product) error
	DeleteProductsBySquareIds(squareIds []string) error

	GetUserByUsername(username string) (types.User, error)
    UpdateUserPassword(userId uint, password string) error

	GetOrdersWithinRange(start, end time.Time) ([]types.Order, error)
	GetOrderBySquareId(squareId string) (types.Order, error)
	InsertOrders(order []types.Order) error

	GetLastSyncStatus() (types.SyncStatus, error)
	InsertSyncStatus(ts int64) error
}
