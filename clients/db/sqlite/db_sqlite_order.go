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
		orders = append(orders, order)
	}

	return orders, nil
}

func (db SqliteDb) InsertOrders(orders []types.Order) error {
	args := make([][]any, len(orders))
	for index, order := range orders {
		args[index] = []any{order.SquareID, order.LocationID, order.State, order.Version, order.TotalMoney, order.TotalTax, order.TotalDiscount, order.TotalTip, order.TotalServiceCharge, order.CreatedAt.Unix()}
	}
	return db.commitTx(query_INSERT_ORDERS, args)
}
