package gormsqlite

import (
	"testing"

	"github.com/kampanosg/go-lsi/models"
	"github.com/kampanosg/go-lsi/types"
)

func TestDbAuth(t *testing.T) {

	tests := []struct {
		name     string
		username string
		user     types.User
		hasError bool
	}{
		{"return error when user doesn't exist", "test", types.User{}, true},
		{"return user when exists", "darth-vader", types.User{Username: "darth-vader"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.User{Username: "darth-vader", Password: "empire-rocks-123"})

			user, err := db.GetUserByUsername(tt.username)
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}

			if user.Username != tt.user.Username {
				t.Errorf("got %s, want %s", user.Username, tt.user.Username)
			}
		})
	}
}

func TestDbAuth_ChangePassword(t *testing.T) {

	tests := []struct {
		name     string
		userId   uint
		password string
		hasError bool
	}{
		{"return no error when user exists", 1, "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.User{Username: "darth-vader", Password: "empire-rocks-123"})

			err = db.UpdateUserPassword(tt.userId, tt.password)
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}
		})
	}
}
