package gormsqlite

import (
	"github.com/kampanosg/go-lsi/models"
	"github.com/kampanosg/go-lsi/types"
)

func (db SqliteDb) GetUserByUsername(username string) (types.User, error) {
	var result models.User
	err := db.Connection.Where("username = ?", username).First(&result).Error
	return fromUserDbRowToType(result), err
}

func (db SqliteDb) UpdateUserPassword(userId uint, password string) error {
	return db.Connection.Model(&models.User{}).Where("id = ?", userId).Update("password", password).Error
}

func fromUserDbRowToType(user models.User) types.User {
	return types.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}
}
