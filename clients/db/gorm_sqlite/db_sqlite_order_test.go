package gormsqlite

import (
	"testing"
	"time"

	"github.com/kampanosg/go-lsi/models"
	"github.com/kampanosg/go-lsi/types"
)

var (
	date1     = time.Date(2023, 2, 10, 21, 0, 0, 0, time.Local)
	date2     = time.Date(2023, 2, 10, 21, 45, 0, 0, time.Local)
	date3     = time.Date(2023, 2, 10, 12, 0, 0, 0, time.Local)
	date4     = time.Date(2023, 2, 10, 15, 45, 0, 0, time.Local)
	orderDate = time.Date(2023, 2, 10, 13, 30, 0, 0, time.Local)
)

func TestDbOrders_WithinRange(t *testing.T) {
	tests := []struct {
		name        string
		start       time.Time
		end         time.Time
		expectedLen int
	}{
		{"returns empty array when no orders in range", date1, date2, 0},
		{"returns orders within range", date3, date4, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Order{CreatedAtSquare: orderDate})

			orders, err := db.GetOrdersWithinRange(tt.start, tt.end)
			if err != nil {
				t.Errorf("threw unexpected error, got %s", err.Error())
			}

			if len(orders) != tt.expectedLen {
				t.Errorf("got %d, want %d", len(orders), tt.expectedLen)
			}
		})
	}
}

func TestDbOrders_InsertOrders(t *testing.T) {
	tests := []struct {
		name   string
		orders []types.Order
		err    error
	}{
		{"succeeds when adding order", []types.Order{{SquareID: "test-order-1"}}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			if err := db.InsertOrders(tt.orders); err != tt.err {
				t.Errorf("threw unexpected error, got %s, want %s", err.Error(), tt.err.Error())
			}
		})
	}
}

func TestDbOrders_GetOrderBySquareId(t *testing.T) {
	tests := []struct {
		name     string
		squareId string
		hasErr   bool
	}{
		{"fails for id that doesnt exist", "square-id-2", true},
		{"succeeds for id that exists", "square-id-1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Order{SquareID: "square-id-1"})

			_, err = db.GetOrderBySquareId(tt.squareId)
			if tt.hasErr && err == nil {
				t.Errorf("was expecting error, got %v", err)
			}

		})
	}
}
