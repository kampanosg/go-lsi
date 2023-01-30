package main

import (
	"os"

	"github.com/kampanosg/go-lsi/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	args := os.Args[1:]

	username := args[1]
	password, err := hashPassword(args[2])
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open(args[0]), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	user := models.User{Username: username, Password: password}
	db.Create(&user)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
