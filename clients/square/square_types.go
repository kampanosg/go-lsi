package square

import "time"

type BatchDeleteItemsRequest struct {
	ObjectIds []string `json:"object_ids"`
}

type SquareBatchUpsertCatalogItemRequest struct {
	IdempotencyKey string                `json:"idempotency_key"`
	Batches        []SquareCategoryBatch `json:"batches"`
}

type SquareCategoryBatch struct {
	Objects []SquareUpsertCategoryObject `json:"objects"`
}

type SquareUpsertCategoryObject struct {
	Type         string             `json:"type"`
	Id           string             `json:"id"`
	Version      int64              `json:"version"`
	IsDeleted    bool               `json:"is_deleted"`
	CategoryData SquareCategoryData `json:"category_data"`
}

type SquareCategoryData struct {
	Name string `json:"name"`
}

type SquareBatchUpsertProductRequest struct {
	IdempotencyKey string               `json:"idempotency_key"`
	Batches        []SquareProductBatch `json:"batches"`
}

type SquareProductBatch struct {
	Objects []SquareProductObject `json:"objects"`
}

type SquareProductObject struct {
	Type                  string            `json:"type"`
	ID                    string            `json:"id"`
	IsDeleted             bool              `json:"is_deleted"`
	PresentAtAllLocations bool              `json:"present_at_all_locations"`
	Version               int64             `json:"version"`
	ItemData              SquareProductData `json:"item_data"`
}

type SquareProductData struct {
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	IsTaxable          bool                     `json:"is_taxable"`
	Visibility         string                   `json:"visibility"`
	CategoryID         string                   `json:"category_id"`
	Variations         []SquareProductVariation `json:"variations"`
	ProductType        string                   `json:"product_type"`
	SkipModifierScreen bool                     `json:"skip_modifier_screen"`
}

type SquareProductVariation struct {
	Type                  string                     `json:"type"`
	ID                    string                     `json:"id"`
	IsDeleted             bool                       `json:"is_deleted"`
	PresentAtAllLocations bool                       `json:"present_at_all_locations"`
	Version               int64                      `json:"version"`
	ItemVariationData     SquareProductVariationData `json:"item_variation_data"`
}

type SquareProductVariationData struct {
	ItemID      string           `json:"item_id"`
	Name        string           `json:"name"`
	Sku         string           `json:"sku"`
	Ordinal     int              `json:"ordinal"`
	PricingType string           `json:"pricing_type"`
	Upc         string           `json:"upc"`
	PriceMoney  SquarePriceMoney `json:"price_money"`
}

type SquarePriceMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type SquareUpsertResponse struct {
	Objects []struct {
		Type                  string    `json:"type"`
		ID                    string    `json:"id"`
		UpdatedAt             time.Time `json:"updated_at"`
		CreatedAt             time.Time `json:"created_at"`
		Version               int64     `json:"version"`
		IsDeleted             bool      `json:"is_deleted"`
		PresentAtAllLocations bool      `json:"present_at_all_locations"`
		CategoryData          struct {
			Name       string `json:"name"`
			IsTopLevel bool   `json:"is_top_level"`
		} `json:"category_data"`
	} `json:"objects"`
	IDMappings []SquareIdMapping `json:"id_mappings"`
}

type SquareIdMapping struct {
	ClientObjectID string `json:"client_object_id"`
	ObjectID       string `json:"object_id"`
}

