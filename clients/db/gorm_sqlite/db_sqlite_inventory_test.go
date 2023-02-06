package gormsqlite

import (
	"testing"

	"github.com/kampanosg/go-lsi/models"
	"github.com/kampanosg/go-lsi/types"
)

func TestDbInventory_GetCategories(t *testing.T) {
	tests := []struct {
		name       string
		categories []types.Category
		hasError   bool
	}{
		{"return all categories", []types.Category{{Name: "test"}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Category{Name: "test category 1"})

			categories, err := db.GetCategories()
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}

			if len(categories) != 1 {
				t.Errorf("got %d, want %d", len(categories), len(tt.categories))
			}
		})
	}
}

func TestDbInventory_InsertCategory(t *testing.T) {
	tests := []struct {
		name        string
		categories  []types.Category
		hasError    bool
		expectedLen int
	}{
		{"inserts category", []types.Category{{Name: "test", LinnworksID: "test-id-1", SquareID: "square-id-1"}}, false, 1},
		{"returns error for duplicate keys", []types.Category{
			{Name: "test 2", LinnworksID: "test-id-2", SquareID: "square-id-2"},
			{Name: "test 2", LinnworksID: "test-id-2", SquareID: "square-id-2"},
		}, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			err = db.InsertCategories(tt.categories)
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}

			categories, err := db.GetCategories()
			if err != nil {
				t.Errorf("unexpected error, got %v", err.Error())
			}

			if len(categories) != tt.expectedLen {
				t.Errorf("got %d, want %d", len(categories), tt.expectedLen)
			}
		})
	}

}

func TestDbInventory_UpsertCategory(t *testing.T) {
	tests := []struct {
		name     string
		category types.Category
		hasError bool
	}{
		{"inserts new category", types.Category{Name: "test", LinnworksID: "test-id-1", SquareID: "square-id-1"}, false},
		{"updates existing category", types.Category{Name: "test 2", LinnworksID: "test-id-2", SquareID: "square-id-2"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Category{Name: "Existing Category", LinnworksID: "test-id-2", SquareID: "square-id-2"})

			err = db.UpsertCategory(tt.category)
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}

			var res models.Category
			if err := db.Connection.Where("square_id = ?", tt.category.SquareID).First(&res).Error; err != nil {
				t.Errorf("threw unexpected error, got %s", err.Error())
			}

			if res.Name != tt.category.Name {
				t.Errorf("category has not been upserted, got %s, want %s", res.Name, tt.category.Name)
			}
		})
	}
}