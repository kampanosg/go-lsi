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
		var id, createdAt int
		var squareId, locationId, state string
		var totalMoney float64
		if rows.Scan(&id, &squareId, &locationId, &state, &totalMoney, &createdAt); err != nil {
			return orders, err
		}
		order := transformers.FromOrderDbRowToDomain(id, squareId, locationId, state, totalMoney, createdAt)
		orders = append(orders, order)
	}

	return orders, nil
}

func (db SqliteDb) InsertOrders([]types.Order) error {
	return nil
}

func (db SqliteDb) InsertOrderProducts([]types.OrderProduct) error {
	return nil
}
