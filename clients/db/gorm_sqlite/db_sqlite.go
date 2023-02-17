package gormsqlite

import (
	"errors"

	"github.com/kampanosg/go-lsi/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	errRecordNotFound = errors.New("query returned no results")
)

type SqliteDb struct {
	Connection *gorm.DB
}

func NewSqliteDb(dbPath string) (SqliteDb, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		CreateBatchSize:   5000,
		AllowGlobalUpdate: true,
	})

	if err != nil {
		return SqliteDb{}, err
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.SyncStatus{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Order{})


    db.Unscoped().Delete(&models.User{})
    db.Save(&models.User{ Username: "geroge", Password: "reset-password" })
    db.Save(&models.User{ Username: "kev", Password: "reset-password" })

	return SqliteDb{Connection: db}, nil
}
