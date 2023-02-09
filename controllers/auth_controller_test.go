package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	mock "github.com/kampanosg/go-lsi/clients/db/mock_db"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

const (
	SigningKey   = "test-signing-key"
	EasyPassword = "$2a$14$ZbEQyjYsmCGBC8zAviOobuWQW/2YxjO0eyUS8St1p.9r0D.wVPM9K"
)

var (
	errUserNotFound = errors.New("user with that username was not found")
)

func TestAuthController(t *testing.T) {

	tests := []struct {
		name   string
		method string
		req    types.AuthRequest
		dbRes  types.User
		dbErr  error
		status int
	}{
		{"fail when receive GET req", http.MethodGet, types.AuthRequest{}, types.User{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive CONNECT req", http.MethodConnect, types.AuthRequest{}, types.User{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive DELETE req", http.MethodDelete, types.AuthRequest{}, types.User{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive HEAD req", http.MethodHead, types.AuthRequest{}, types.User{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive OPTIONS req", http.MethodOptions, types.AuthRequest{}, types.User{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive PATCH req", http.MethodPatch, types.AuthRequest{}, types.User{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive PUT req", http.MethodPut, types.AuthRequest{}, types.User{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive TRACE req", http.MethodTrace, types.AuthRequest{}, types.User{}, nil, http.StatusMethodNotAllowed},
		{"fail when user is not found", http.MethodPost, types.AuthRequest{Username: "darth-vadre", Password: "="}, types.User{}, errUserNotFound, http.StatusUnauthorized},
		{"fail when password is wrong", http.MethodPost, types.AuthRequest{Username: "darth-vader", Password: "password"}, types.User{Username: "darth-vader", Password: EasyPassword}, nil, http.StatusUnauthorized},
		{"succeed when correct auth details", http.MethodPost, types.AuthRequest{Username: "darth-vader", Password: "empire-123"}, types.User{Username: "darth-vader", Password: EasyPassword}, nil, http.StatusOK},
	}

	logger, _ := zap.NewDevelopment()
	db := new(mock.MockDb)
	sk := []byte(SigningKey)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.On("GetUserByUsername", tt.req.Username).Return(tt.dbRes, tt.dbErr)

			ctrl := NewAuthController(db, sk, logger.Sugar())

			router := mux.NewRouter()
			router.HandleFunc("/", ctrl.HandleAuthRequest)

			b, err := json.Marshal(tt.req)
			if err != nil {
				t.Errorf("threw unexpected error, got %s", err.Error())
			}

			req, err := http.NewRequest(tt.method, "/", bytes.NewBuffer(b))
			if err != nil {
				t.Errorf("threw error when calling endpoint, got %s", err.Error())
			}

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.status)
			}
		})
	}
}
