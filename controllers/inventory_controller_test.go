package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	mock "github.com/kampanosg/go-lsi/clients/db/mock_db"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

var (
	errProductNotFound = errors.New("product not found")
)

func TestInventoryController(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		barcode string
		sku     string
		dbRes   types.Product
		dbErr   error
		status  int
	}{
		{"fail when receive POST req", http.MethodPost, "", "", types.Product{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive CONNECT req", http.MethodConnect, "", "", types.Product{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive DELETE req", http.MethodDelete, "", "", types.Product{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive HEAD req", http.MethodHead, "", "", types.Product{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive OPTIONS req", http.MethodOptions, "", "", types.Product{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive PATCH req", http.MethodPatch, "", "", types.Product{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive PUT req", http.MethodPut, "", "", types.Product{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive TRACE req", http.MethodTrace, "", "", types.Product{}, nil, http.StatusMethodNotAllowed},
		{"fail when no url params are not provided", http.MethodGet, "", "", types.Product{}, errProductNotFound, http.StatusBadRequest},
		{"return 404 when db returns err", http.MethodGet, "test", "", types.Product{}, errProductNotFound, http.StatusNotFound},
		{"return ok when db finds product", http.MethodGet, "test", "", types.Product{Title: "Test Product"}, nil, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logger, _ := zap.NewDevelopment()
			db := new(mock.MockDb)

			db.On("GetProductByBarcode", tt.barcode).Return(tt.dbRes, tt.dbErr)
			db.On("GetProductBySku", tt.sku).Return(tt.dbRes, tt.dbErr)

			ctrl := NewInventoryController(db, logger.Sugar())

			router := mux.NewRouter()
			router.HandleFunc("/inventory", ctrl.HandleInventoryRequest)

			var url strings.Builder
			url.WriteString("/inventory")

			if tt.barcode != "" {
				url.WriteString("?barcode=")
				url.WriteString(tt.barcode)
			}

			if tt.sku != "" {
				url.WriteString("?sku=")
				url.WriteString(tt.sku)
			}
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
