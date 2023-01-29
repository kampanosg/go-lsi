package gormsqlite

import (
	"github.com/kampanosg/go-lsi/models"
	"github.com/kampanosg/go-lsi/types"
)

func (db SqliteDb) GetOrders() ([]types.Order, error) {
	orders := make([]models.Order, 0)
	result := db.Connection.Find(&orders)
	if result.Error != nil {
		return []types.Order{}, result.Error
	}
	return fromOrderModelsToType(orders), nil
}

func (db SqliteDb) InsertOrders(orders []types.Order) error {
	orderModels := fromOrderTypesToModels(orders)
	result := db.Connection.Create(orderModels)
	return result.Error
}

func fromOrderModelsToType(orderModels []models.Order) []types.Order {
	orders := make([]types.Order, len(orderModels))
	for index, orderModel := range orderModels {
		order := types.Order{
			ID:         orderModel.ID,
			LocationID: orderModel.LocationID,
			State:      orderModel.State,
			Version:    orderModel.Version,
			TotalMoney: orderModel.TotalMoney,
			CreatedAt:  orderModel.CreatedAtSquare,
		}
		orders[index] = order
	}
	return orders
}

func fromOrderTypesToModels(orders []types.Order) []models.Order {
	orderModels := make([]models.Order, len(orders))
	for index, order := range orders {
		orderModel := models.Order{
			LocationID:      order.LocationID,
			State:           order.State,
			Version:         order.Version,
			TotalMoney:      order.TotalMoney,
			CreatedAtSquare: order.CreatedAt,
		}
		orderModels[index] = orderModel
	}
	return orderModels
}
