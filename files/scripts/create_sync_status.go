package main

import (
	"os"
	"strconv"
	"time"

	"github.com/kampanosg/go-lsi/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	args := os.Args[1:]

	db, err := gorm.Open(sqlite.Open(args[0]), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	tsMillis, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		panic(err)
	}

	user := models.SyncStatus{LastRun: time.UnixMilli(tsMillis), Success: true}
	db.Create(&user)
}
