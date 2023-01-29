package gormsqlite

import (
	"github.com/kampanosg/go-lsi/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteDb struct {
	Connection *gorm.DB
}

func NewSqliteDb(dbPath string) (SqliteDb, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return SqliteDb{}, err
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.SyncStatus{})
    db.AutoMigrate(&models.Category{})
    db.AutoMigrate(&models.Product{})
    db.AutoMigrate(&models.Order{})

	return SqliteDb{Connection: db}, nil
}
