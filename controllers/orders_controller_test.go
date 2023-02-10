package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	mock "github.com/kampanosg/go-lsi/clients/db/mock_db"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

var (
	d1                = time.Date(2023, 2, 10, 12, 0, 0, 0, time.UTC)
	d2                = time.Date(2023, 2, 10, 21, 0, 0, 0, time.UTC)
	errOrdersNotFound = errors.New("failed to retrieve orders")
)

func TestOrdersControllers(t *testing.T) {
	tests := []struct {
		name   string
		method string
		start  string
		end    string
		dbRes  []types.Order
		dbErr  error
		status int
	}{
		{"fail when receive POST req", http.MethodPost, "", "", []types.Order{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive CONNECT req", http.MethodConnect, "", "", []types.Order{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive DELETE req", http.MethodDelete, "", "", []types.Order{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive HEAD req", http.MethodHead, "", "", []types.Order{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive OPTIONS req", http.MethodOptions, "", "", []types.Order{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive PATCH req", http.MethodPatch, "", "", []types.Order{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive PUT req", http.MethodPut, "", "", []types.Order{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive TRACE req", http.MethodTrace, "", "", []types.Order{}, nil, http.StatusMethodNotAllowed},
		{"fail when no url params are not provided", http.MethodGet, "", "", []types.Order{}, nil, http.StatusBadRequest},
		{"fail when start url param not provided", http.MethodGet, "", "2023-02-10T16:00:00Z", []types.Order{}, nil, http.StatusBadRequest},
		{"fail when end url param not provided", http.MethodGet, "2023-02-10T12:00:00Z", "", []types.Order{}, nil, http.StatusBadRequest},
		{"fail when url params are malformed", http.MethodGet, "boop", "beep", []types.Order{}, nil, http.StatusBadRequest},
		{"fail when start url param is malformed", http.MethodGet, "boop", "2023-02-10T16:00:00Z", []types.Order{}, nil, http.StatusBadRequest},
		{"fail when end url param is malformed", http.MethodGet, "2023-02-10T16:00:00Z", "beep", []types.Order{}, nil, http.StatusBadRequest},
		{"fail db returns error", http.MethodGet, "2023-02-10T12:00:00Z", "2023-02-10T21:00:00Z", []types.Order{}, errOrdersNotFound, http.StatusBadRequest},
		{"ok when db returns orders", http.MethodGet, "2023-02-10T12:00:00Z", "2023-02-10T21:00:00Z", []types.Order{{SquareID: "square-id"}}, nil, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logger, _ := zap.NewDevelopment()
			db := new(mock.MockDb)

			db.On("GetOrdersWithinRange", d1, d2).Return(tt.dbRes, tt.dbErr)

			ctrl := NewOrdersController(db, logger.Sugar())

			router := mux.NewRouter()
			router.HandleFunc("/orders", ctrl.HandleOrdersRequest)

			var url strings.Builder
			url.WriteString("/orders")

			url.WriteString("?start=")
			url.WriteString(tt.start)

			url.WriteString("&end=")
			url.WriteString(tt.end)

			req, err := http.NewRequest(tt.method, url.String(), nil)
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
