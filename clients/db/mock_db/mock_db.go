package mock

import (
	"github.com/kampanosg/go-lsi/types"
	"github.com/stretchr/testify/mock"
)

type MockDb struct {
	mock.Mock
}

func (db *MockDb) GetCategories() ([]types.Category, error) {
	args := db.Called()
	return args.Get(0).([]types.Category), args.Error(1)
}

func (db *MockDb) InsertCategories(categories []types.Category) error {
	args := db.Called(categories)
	return args.Error(0)
}

func (db *MockDb) UpsertCategory(category types.Category) error {
	args := db.Called(category)
	return args.Error(0)
}

func (db *MockDb) DeleteCategoriesBySquareIds(squareIds []string) error {
	args := db.Called(squareIds)
	return args.Error(0)
}

func (db *MockDb) GetProducts() ([]types.Product, error) {
	args := db.Called()
	return args.Get(0).([]types.Product), args.Error(1)
}

func (db *MockDb) GetProductByBarcode(barcode string) (types.Product, error) {
	args := db.Called(barcode)
	return args.Get(0).(types.Product), args.Error(1)
}

func (db *MockDb) GetProductBySku(sku string) (types.Product, error) {
	args := db.Called(sku)
	return args.Get(0).(types.Product), args.Error(1)
}

func (db *MockDb) GetProductByVarId(varId string) (types.Product, error) {
	args := db.Called(varId)
	return args.Get(0).(types.Product), args.Error(1)
}

func (db *MockDb) GetProductByTitle(title string) (types.Product, error) {
	args := db.Called(title)
	return args.Get(0).(types.Product), args.Error(1)
}

func (db *MockDb) UpsertProduct(product types.Product) error {
	args := db.Called(product)
	return args.Error(0)
}

func (db *MockDb) DeleteProductsBySquareIds(squareIds []string) error {
	args := db.Called(squareIds)
	return args.Error(0)
}

func (db *MockDb) GetUserByUsername(username string) (types.User, error) {
	args := db.Called(username)
	return args.Get(0).(types.User), args.Error(1)
}

func (db *MockDb) GetOrders() ([]types.Order, error) {
	args := db.Called()
	return args.Get(0).([]types.Order), args.Error(1)
}

func (db *MockDb) GetOrderBySquareId(squareId string) (types.Order, error) {
	args := db.Called(squareId)
	return args.Get(0).(types.Order), args.Error(1)
}

func (db *MockDb) InsertOrders(orders []types.Order) error {
	args := db.Called(orders)
	return args.Error(0)
}

func (db *MockDb) GetLastSyncStatus() (types.SyncStatus, error) {
	args := db.Called()
	return args.Get(0).(types.SyncStatus), args.Error(1)
}

func (db *MockDb) InsertSyncStatus(ts int64) error {
	args := db.Called(ts)
	return args.Error(0)
}
