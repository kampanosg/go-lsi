package sqlite

import "github.com/kampanosg/go-lsi/types"

func (db SqliteDb) GetUserByUsername(username string) (types.User, error) {
	row := db.Connection.QueryRow(query_GET_USER_BY_USERNAME, username)

	var id int
	var uname, password, friendlyName string

	err := row.Scan(&id, &uname, &password, &friendlyName)
	if err != nil {
		return types.User{}, err
	}

	return types.User{
		Id:           id,
		Username:     uname,
		Password:     password,
		FriendlyName: friendlyName,
	}, nil
}