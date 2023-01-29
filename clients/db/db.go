package db

import (
	"github.com/kampanosg/go-lsi/types"
)

type DB interface {
	GetInventory() ([]types.InventoryItem, error)

	GetCategories() ([]types.Category, error)
	InsertCategories(categories []types.Category) error
	ClearCategories() error

	GetProducts() ([]types.Product, error)
	GetProductByVarId(varId string) (types.Product, error)
	InsertProducts(products []types.Product) error
	ClearProducts() error

	GetUserByUsername(username string) (types.User, error)

	GetOrders() ([]types.Order, error)
	InsertOrders([]types.Order) error

	GetLastSyncStatus() (types.SyncStatus, error)
	InsertSyncStatus(ts int64) error
}
