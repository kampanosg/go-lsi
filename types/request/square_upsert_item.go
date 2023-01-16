package request

type SquareBatchUpsertItemRequest struct {
	IdempotencyKey string      `json:"idempotency_key"`
	Batches        []ItemBatch `json:"batches"`
}

type PriceMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type ItemVariationData struct {
	ItemID      string     `json:"item_id"`
	Name        string     `json:"name"`
	Sku         string     `json:"sku"`
	Ordinal     int        `json:"ordinal"`
	PricingType string     `json:"pricing_type"`
	PriceMoney  PriceMoney `json:"price_money"`
}

type ItemVariation struct {
	Type                  string            `json:"type"`
	ID                    string            `json:"id"`
	IsDeleted             bool              `json:"is_deleted"`
	PresentAtAllLocations bool              `json:"present_at_all_locations"`
	Version               int64             `json:"version"`
	ItemVariationData     ItemVariationData `json:"item_variation_data"`
}

type ItemData struct {
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	IsTaxable          bool            `json:"is_taxable"`
	Visibility         string          `json:"visibility"`
	CategoryID         string          `json:"category_id"`
	Variations         []ItemVariation `json:"variations"`
	ProductType        string          `json:"product_type"`
	SkipModifierScreen bool            `json:"skip_modifier_screen"`
}

type ItemObject struct {
	Type                  string   `json:"type"`
	ID                    string   `json:"id"`
	IsDeleted             bool     `json:"is_deleted"`
	PresentAtAllLocations bool     `json:"present_at_all_locations"`
	Version               int64    `json:"version"`
	ItemData              ItemData `json:"item_data"`
}

type ItemBatch struct {
	Objects []ItemObject `json:"objects"`
}
