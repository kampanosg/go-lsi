package linnworks

import (
	"github.com/kampanosg/go-lsi/types"
	"github.com/stretchr/testify/mock"
)

type LinnworksMockClient struct {
	mock.Mock
}

func (c *LinnworksMockClient) GetCategories() ([]LinnworksCategoryResponse, error) {
	args := c.Called()
	return args.Get(0).([]LinnworksCategoryResponse), args.Error(1)
}

func (c *LinnworksMockClient) GetProducts() ([]LinnworksProductResponse, error) {
	args := c.Called()
	return args.Get(0).([]LinnworksProductResponse), args.Error(1)
}

func (c *LinnworksMockClient) CreateOrders(orders []types.Order) (LinnworksCreateOrdersResponse, error) {
	args := c.Called(orders)
	return args.Get(0).(LinnworksCreateOrdersResponse), args.Error(1)
}
