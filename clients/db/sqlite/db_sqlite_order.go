package sqlite

import (
	"github.com/kampanosg/go-lsi/transformers"
	"github.com/kampanosg/go-lsi/types"
)

func (db SqliteDb) GetOrders() ([]types.Order, error) {
	rows, err := db.Connection.Query(query_GET_ORDERS)
	orders := make([]types.Order, 0)
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var squareId, locationId, state string
		var totalMoney float64
		var createdAt int64
		if rows.Scan(&id, &squareId, &locationId, &state, &totalMoney, &createdAt); err != nil {
			return orders, err
		}
		order := transformers.FromOrderDbRowToDomain(id, squareId, locationId, state, totalMoney, createdAt)
		products, err := db.GetOrderProductsForCategory(order.Id)
		if err != nil {
			return orders, err
		}
		order.Products = products
		orders = append(orders, order)
	}

	return orders, nil
}

func (db SqliteDb) GetOrderProductsForCategory(categoryId int) ([]types.OrderProduct, error) {
	rows, err := db.Connection.Query(query_GET_ORDER_PRODUCTS, categoryId)
	products := make([]types.OrderProduct, 0)
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var squareOrderId, squareVarId, qty string
		if rows.Scan(&id, &squareOrderId, &squareVarId, &qty); err != nil {
			return products, err
		}
		order := transformers.FromOrderProductDbRowToDomain(id, squareOrderId, squareVarId, qty)

		products = append(products, order)
	}

	return products, nil
}

func (db SqliteDb) InsertOrders(orders []types.Order) error {
	args := make([][]any, len(orders))
	for index, order := range orders {
		args[index] = []any{order.SquareId, order.LocationId, order.State, order.Version, order.TotalMoney, order.CreatedAt.Unix()}
	}
	return db.commitTx(query_INSERT_ORDERS, args)
}

func (db SqliteDb) InsertOrderProducts(orderProducts []types.OrderProduct) error {
	args := make([][]any, len(orderProducts))
	for index, orderProduct := range orderProducts {
		args[index] = []any{orderProduct.OrderId, orderProduct.SquareOrderId, orderProduct.SquareVarId, orderProduct.Quantity}
	}
	return db.commitTx(query_INSERT_ORDER_PRODUCTS, args)
}
