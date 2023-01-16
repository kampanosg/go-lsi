package db

import (
	"github.com/kampanosg/go-lsi/types"
)

type DB interface {
	GetCategories() ([]types.Category, error)
	InsertCategories(categories []types.Category) error
	ClearCategories() error

	GetProducts() ([]types.Product, error)
	InsertProducts(products []types.Product) error
	ClearProducts() error
}
