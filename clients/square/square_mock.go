package square

import (
	"time"

	"github.com/kampanosg/go-lsi/types"
	"github.com/stretchr/testify/mock"
)

type SquareMockClient struct {
	mock.Mock
}

	func (c *SquareMockClient) GetItemVersion(squareId string) (int64, error) {
        args := c.Called(squareId)
        return args.Get(0).(int64), args.Error(1)
    }

	func (c *SquareMockClient) UpsertCategories(categories []types.Category) (SquareUpsertResponse, error) {
        args := c.Called(categories)
        return args.Get(0).(SquareUpsertResponse), args.Error(1)
    }

	func (c *SquareMockClient) UpsertProducts(products []types.Product) (SquareUpsertResponse, error) {
        args := c.Called(products)
        return args.Get(0).(SquareUpsertResponse), args.Error(1)
    }

	func (c *SquareMockClient) BatchDeleteItems(itemIds []string) error {
        args := c.Called(itemIds)
        return args.Error(1)
    }

	func (c *SquareMockClient) SearchOrders(start time.Time, end time.Time) ([]SquareOrder, error) {
        args := c.Called(start, end)
        return args.Get(0).([]SquareOrder), args.Error(1)
    }
