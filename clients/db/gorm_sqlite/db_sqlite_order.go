package gormsqlite

import (
	"time"

	"github.com/kampanosg/go-lsi/models"
	"github.com/kampanosg/go-lsi/types"
)

func (db SqliteDb) GetOrdersWithinRange(start, end time.Time) ([]types.Order, error) {
	orders := make([]models.Order, 0)
	err := db.Connection.
		Order("created_at_square desc").
		Where("created_at_square BETWEEN ? AND ?", start, end).Find(&orders).Error
	return fromOrderModelsToTypes(orders), err
}

func (db SqliteDb) InsertOrders(orders []types.Order) error {
	orderModels := fromOrderTypesToModels(orders)
	return db.Connection.Create(&orderModels).Error
}

func (db SqliteDb) GetOrderBySquareId(squareId string) (types.Order, error) {
	var result models.Order
	db.Connection.Where(&models.Order{SquareID: squareId}).Limit(1).Find(&result)
	if result.ID == 0 {
		return types.Order{}, errRecordNotFound
	}
	return fromOrderModelToType(result), nil
}

func fromOrderModelsToTypes(orderModels []models.Order) []types.Order {
	orders := make([]types.Order, len(orderModels))
	for index, orderModel := range orderModels {
		orders[index] = fromOrderModelToType(orderModel)
	}
	return orders
}

func fromOrderModelToType(orderModel models.Order) types.Order {
	return types.Order{
		ID:         orderModel.ID,
		SquareID:   orderModel.SquareID,
		LocationID: orderModel.LocationID,
		State:      orderModel.State,
		Version:    orderModel.Version,
		TotalMoney: orderModel.TotalMoney,
		CreatedAt:  orderModel.CreatedAtSquare,
	}
}

func fromOrderTypesToModels(orders []types.Order) []models.Order {
	orderModels := make([]models.Order, len(orders))
	for index, order := range orders {
		orderModel := models.Order{
			SquareID:        order.SquareID,
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
