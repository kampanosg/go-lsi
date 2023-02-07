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

func TestDbInventory_DeleteCategoriesBySquareIds(t *testing.T) {
	tests := []struct {
		name        string
		squareIds   []string
		expectedLen int
	}{
		{"doesn't delete - squareIds is empty", []string{}, 3},
		{"doesn't delete - squareIds are not valid", []string{"bad-square-id-1", "bad-square-id-2"}, 3},
		{"deletes only the correct categories", []string{"square-id-1", "square-id-2"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Category{Name: "Existing Category 1", LinnworksID: "test-id-1", SquareID: "square-id-1"})
			db.Connection.Save(&models.Category{Name: "Existing Category 2", LinnworksID: "test-id-2", SquareID: "square-id-2"})
			db.Connection.Save(&models.Category{Name: "Existing Category 3", LinnworksID: "test-id-3", SquareID: "square-id-3"})

			db.DeleteCategoriesBySquareIds(tt.squareIds)

			res, err := db.GetCategories()
			if err != nil {
				t.Errorf("threw, unexpected error, got %s", err.Error())
			}

			if len(res) != tt.expectedLen {
				t.Errorf("returned wrong res, got %v, want %d", len(res), tt.expectedLen)
			}

		})
	}
}

func TestDbInventory_GetProducts(t *testing.T) {
	tests := []struct {
		name        string
		expextedLen int
		hasError    bool
	}{
		{"return all products", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Product{Title: "Test Product"})

			products, err := db.GetProducts()
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}

			if len(products) != tt.expextedLen {
				t.Errorf("got %d, want %d", len(products), tt.expextedLen)
			}
		})
	}
}

func TestDbInventory_GetProductBySku(t *testing.T) {
	tests := []struct {
		name            string
		sku             string
		hasError        bool
		expectedProduct types.Product
	}{
		{"returns error for invalid sku", "bad-sku", true, types.Product{}},
		{"returns correct product for ok sku", "sku-001", false, types.Product{Title: "Test Product 1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Product{Title: "Test Product 1", SKU: "sku-001"})

			product, err := db.GetProductBySku(tt.sku)
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}

			if tt.expectedProduct.Title != product.Title {
				t.Errorf("returned wrong product, got %s, want %s", product.Title, tt.expectedProduct.Title)
			}
		})
	}
}

func TestDbInventory_GetProductByBarcode(t *testing.T) {
	tests := []struct {
		name            string
		barcode         string
		hasError        bool
		expectedProduct types.Product
	}{
		{"returns error for invalid sku", "bad-barcode", true, types.Product{}},
		{"returns correct product for ok barcode", "barcode-001", false, types.Product{Title: "Test Product 1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Product{Title: "Test Product 1", Barcode: "barcode-001"})

			product, err := db.GetProductByBarcode(tt.barcode)
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}

			if tt.expectedProduct.Title != product.Title {
				t.Errorf("returned wrong product, got %s, want %s", product.Title, tt.expectedProduct.Title)
			}
		})
	}
}

func TestDbInventory_GetProductByVariationId(t *testing.T) {
	tests := []struct {
		name            string
		varId           string
		hasError        bool
		expectedProduct types.Product
	}{
		{"returns error for invalid sku", "bad-variation", true, types.Product{}},
		{"returns correct product for ok variation", "variation-001", false, types.Product{Title: "Test Product 1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Product{Title: "Test Product 1", SquareVarID: "variation-001"})

			product, err := db.GetProductByVarId(tt.varId)
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}

			if tt.expectedProduct.Title != product.Title {
				t.Errorf("returned wrong product, got %s, want %s", product.Title, tt.expectedProduct.Title)
			}
		})
	}
}

func TestDbInventory_GetProductByTitle(t *testing.T) {
	tests := []struct {
		name            string
		title           string
		hasError        bool
		expectedProduct types.Product
	}{
		{"returns error for invalid sku", "Elaborate Title which doesn't Work", true, types.Product{}},
		{"returns correct product for ok variation", "Test Product 1", false, types.Product{Title: "Test Product 1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Product{Title: "Test Product 1"})

			product, err := db.GetProductByTitle(tt.title)
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}

			if tt.expectedProduct.Title != product.Title {
				t.Errorf("returned wrong product, got %s, want %s", product.Title, tt.expectedProduct.Title)
			}
		})
	}
}

func TestDbInventory_UpsertProduct(t *testing.T) {
	tests := []struct {
		name     string
		product  types.Product
		hasError bool
	}{
		{"inserts new product", types.Product{Title: "test", LinnworksID: "test-id-1", SquareID: "square-id-1", SquareVarID: "square-var-id-1", Price: 4.20}, false},
		{"updates existing product", types.Product{Title: "test", LinnworksID: "test-id-2", SquareID: "square-id-2", SquareVarID: "square-var-id-2", Price: 6.9}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Product{Title: "test", LinnworksID: "test-id-2", SquareID: "square-id-2", SquareVarID: "square-var-id-2", Price: 69.99})

			err = db.UpsertProduct(tt.product)
			if tt.hasError && err == nil {
				t.Errorf("expecting to throw error")
			}

			var res models.Product
			if err := db.Connection.Where("square_id = ?", tt.product.SquareID).First(&res).Error; err != nil {
				t.Errorf("threw unexpected error, got %s", err.Error())
			}

			if res.Title != tt.product.Title || res.Price != tt.product.Price {
				t.Errorf("category has not been upserted, title = { got %s, want %s }, price = { got %v, want %v }", res.Title, tt.product.Title, res.Price, tt.product.Price)
			}
		})
	}
}

func TestDbInventory_DeleteProductsBySquareIds(t *testing.T) {
	tests := []struct {
		name        string
		squareIds   []string
		expectedLen int
	}{
		{"doesn't delete - squareIds is empty", []string{}, 3},
		{"doesn't delete - squareIds are not valid", []string{"bad-square-id-1", "bad-square-id-2"}, 3},
		{"deletes only the correct categories", []string{"square-id-1", "square-id-2"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			db, err := NewSqliteDb(tmpDb)
			if err != nil {
				t.Errorf("failed to open db, err=%s", err.Error())
			}

			db.Connection.Save(&models.Product{Title: "Existing Product 1", LinnworksID: "test-id-1", SquareID: "square-id-1", SquareVarID: "square-var-id-1"})
			db.Connection.Save(&models.Product{Title: "Existing Product 2", LinnworksID: "test-id-2", SquareID: "square-id-2", SquareVarID: "square-var-id-2"})
			db.Connection.Save(&models.Product{Title: "Existing Product 3", LinnworksID: "test-id-3", SquareID: "square-id-3", SquareVarID: "square-var-id-3"})

			db.DeleteProductsBySquareIds(tt.squareIds)

			res, err := db.GetProducts()
			if err != nil {
				t.Errorf("threw, unexpected error, got %s", err.Error())
			}

			if len(res) != tt.expectedLen {
				t.Errorf("returned wrong res, got %v, want %d", len(res), tt.expectedLen)
			}

		})
	}
}
