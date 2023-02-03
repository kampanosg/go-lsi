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
	ItemID              string           `json:"item_id"`
	Name                string           `json:"name"`
	Sku                 string           `json:"sku"`
	Ordinal             int              `json:"ordinal"`
	PricingType         string           `json:"pricing_type"`
	Upc                 string           `json:"upc"`
	ServiceDuration     int              `json:"service_duration"`
	PriceMoney          SquarePriceMoney `json:"price_money"`
	AvailableForBooking bool             `json:"available_for_booking"`
	TeamMemberIds       []string         `json:"team_member_ids"`
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

type SquareOrderSearchResponse struct {
	Orders []SquareOrder `json:"orders"`
	Cursor string        `json:"cursor"`
}
type SquareBasePriceMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareGrossSalesMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareTotalTaxMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareTotalDiscountMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareTotalMoney struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
type SquareVariationTotalPriceMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareLineItem struct {
	UID                      string                         `json:"uid"`
	CatalogObjectID          string                         `json:"catalog_object_id"`
	CatalogVersion           int64                          `json:"catalog_version"`
	Quantity                 string                         `json:"quantity"`
	Name                     string                         `json:"name"`
	VariationName            string                         `json:"variation_name"`
	BasePriceMoney           SquareBasePriceMoney           `json:"base_price_money"`
	Note                     string                         `json:"note"`
	GrossSalesMoney          SquareGrossSalesMoney          `json:"gross_sales_money"`
	TotalTaxMoney            SquareTotalTaxMoney            `json:"total_tax_money"`
	TotalDiscountMoney       SquareTotalDiscountMoney       `json:"total_discount_money"`
	TotalMoney               SquareTotalMoney               `json:"total_money"`
	VariationTotalPriceMoney SquareVariationTotalPriceMoney `json:"variation_total_price_money"`
	ItemType                 string                         `json:"item_type"`
}
type SquareTotalTipMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareTotalServiceChargeMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareTaxMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareDiscountMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareTipMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type ServiceChargeMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareNetAmounts struct {
	TotalMoney         SquareTotalMoney    `json:"total_money"`
	TaxMoney           SquareTaxMoney      `json:"tax_money"`
	DiscountMoney      SquareDiscountMoney `json:"discount_money"`
	TipMoney           SquareTipMoney      `json:"tip_money"`
	ServiceChargeMoney ServiceChargeMoney  `json:"service_charge_money"`
}
type SquareSource struct {
	Name string `json:"name"`
}
type SquareNetAmountDueMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type SquareOrder struct {
	ID                      string                        `json:"id"`
	LocationID              string                        `json:"location_id"`
	LineItems               []SquareLineItem              `json:"line_items"`
	CreatedAt               time.Time                     `json:"created_at"`
	UpdatedAt               time.Time                     `json:"updated_at"`
	State                   string                        `json:"state"`
	Version                 int64                         `json:"version"`
	TotalTaxMoney           SquareTotalTaxMoney           `json:"total_tax_money"`
	TotalDiscountMoney      SquareTotalDiscountMoney      `json:"total_discount_money"`
	TotalTipMoney           SquareTotalTipMoney           `json:"total_tip_money"`
	TotalMoney              SquareTotalMoney              `json:"total_money"`
	TotalServiceChargeMoney SquareTotalServiceChargeMoney `json:"total_service_charge_money"`
	NetAmounts              SquareNetAmounts              `json:"net_amounts"`
	Source                  SquareSource                  `json:"source"`
	NetAmountDueMoney       SquareNetAmountDueMoney       `json:"net_amount_due_money"`
}

type SquareSearchOrdersRequest struct {
	ReturnEntries bool        `json:"return_entries"`
	Limit         int         `json:"limit"`
	Query         SquareQuery `json:"query"`
	LocationIds   []string    `json:"location_ids"`
	Cursor        string      `json:"cursor"`
}
type SquareDateRange struct {
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}
type SquareDateTimeFilter struct {
	CreatedAt SquareDateRange `json:"created_at"`
}
type SquareFilter struct {
	DateTimeFilter SquareDateTimeFilter `json:"date_time_filter"`
}
type SquareQuery struct {
	Filter SquareFilter `json:"filter"`
}

type SquareCatalogItemResponse struct {
	Object struct {
		Type                  string    `json:"type"`
		ID                    string    `json:"id"`
		UpdatedAt             time.Time `json:"updated_at"`
		CreatedAt             time.Time `json:"created_at"`
		Version               int64     `json:"version"`
		IsDeleted             bool      `json:"is_deleted"`
		PresentAtAllLocations bool      `json:"present_at_all_locations"`
		ItemData              struct {
			Name       string `json:"name"`
			IsTaxable  bool   `json:"is_taxable"`
			Visibility string `json:"visibility"`
			CategoryID string `json:"category_id"`
			Variations []struct {
				Type                  string    `json:"type"`
				ID                    string    `json:"id"`
				UpdatedAt             time.Time `json:"updated_at"`
				CreatedAt             time.Time `json:"created_at"`
				Version               int64     `json:"version"`
				IsDeleted             bool      `json:"is_deleted"`
				PresentAtAllLocations bool      `json:"present_at_all_locations"`
				ItemVariationData     struct {
					ItemID      string `json:"item_id"`
					Name        string `json:"name"`
					Sku         string `json:"sku"`
					Upc         string `json:"upc"`
					Ordinal     int    `json:"ordinal"`
					PricingType string `json:"pricing_type"`
					PriceMoney  struct {
						Amount   int    `json:"amount"`
						Currency string `json:"currency"`
					} `json:"price_money"`
					ServiceDuration     int  `json:"service_duration"`
					AvailableForBooking bool `json:"available_for_booking"`
					Sellable            bool `json:"sellable"`
					Stockable           bool `json:"stockable"`
				} `json:"item_variation_data"`
			} `json:"variations"`
			ProductType        string `json:"product_type"`
			SkipModifierScreen bool   `json:"skip_modifier_screen"`
			EcomAvailable      bool   `json:"ecom_available"`
			EcomVisibility     string `json:"ecom_visibility"`
		} `json:"item_data"`
	} `json:"object"`
}
