package gormsqlite

import "github.com/kampanosg/go-lsi/types"

func (db SqliteDb) GetOrders() ([]types.Order, error)                                        {}
func (db SqliteDb) GetOrderProductsForCategory(categoryId int) ([]types.OrderProduct, error) {}
func (db SqliteDb) InsertOrders([]types.Order) error                                         {}
func (db SqliteDb) InsertOrderProducts([]types.OrderProduct) error {
}
