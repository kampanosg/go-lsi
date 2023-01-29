package gormsqlite

import "github.com/kampanosg/go-lsi/types"

func (db SqliteDb) GetCategories() ([]types.Category, error)           {}
func (db SqliteDb) InsertCategories(categories []types.Category) error {}
func (db SqliteDb) ClearCategories() error                             {}

func (db SqliteDb) GetProducts() ([]types.Product, error)                 {}
func (db SqliteDb) GetProductByVarId(varId string) (types.Product, error) {}
func (db SqliteDb) InsertProducts(products []types.Product) error         {}
func (db SqliteDb) ClearProducts() error                                  {}
