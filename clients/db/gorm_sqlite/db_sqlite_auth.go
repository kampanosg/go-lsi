package gormsqlite

import ( 
    "github.com/kampanosg/go-lsi/types"
    "github.com/kampanosg/go-lsi/models"
)

func (db SqliteDb) GetUserByUsername(username string) (types.User, error) {
    var result models.User
    db.Connection.Where("username = ?", username).First(&result) 
    if result.ID == 0 {
        return types.User{}, errRecordNotFound
    }
    return fromUserDbRowToType(result), nil
}

func fromUserDbRowToType(user models.User) types.User {
    return types.User {
        ID: user.ID,
        Username: user.Username,
        Password: user.Password,
    }
}
