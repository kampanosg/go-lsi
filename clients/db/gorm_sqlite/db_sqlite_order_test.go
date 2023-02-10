package gormsqlite

import (
	"testing"
	"time"

	"github.com/kampanosg/go-lsi/models"
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
