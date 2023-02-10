package gormsqlite

import (
	"errors"
	"time"

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

	db.Save(&models.Order{SquareID: "sample-id-1", LocationID: "location-1", TotalMoney: 100, CreatedAtSquare: time.Now()})
	db.Save(&models.Order{SquareID: "sample-id-2", LocationID: "location-1", TotalMoney: 99, CreatedAtSquare: time.Now()})
	db.Save(&models.Order{SquareID: "sample-id-3", LocationID: "location-1", TotalMoney: 98, CreatedAtSquare: time.Now()})

	return SqliteDb{Connection: db}, nil
}
