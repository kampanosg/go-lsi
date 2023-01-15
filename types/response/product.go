package response

import "time"

type ProductResponse struct {
	Suppliers                []interface{}          `json:"Suppliers"`
	StockLevels              []interface{}          `json:"StockLevels"`
	ItemChannelDescriptions  []interface{}          `json:"ItemChannelDescriptions"`
	ItemExtendedProperties   []interface{}          `json:"ItemExtendedProperties"`
	ItemChannelTitles        []interface{}          `json:"ItemChannelTitles"`
	ItemChannelPrices        []interface{}          `json:"ItemChannelPrices"`
	Images                   []ProductImageResponse `json:"Images"`
	SKU                      string                 `json:"ItemNumber"`
	ItemTitle                string                 `json:"ItemTitle"`
	BarcodeNumber            string                 `json:"BarcodeNumber"`
	MetaData                 string                 `json:"MetaData"`
	IsVariationParent        bool                   `json:"IsVariationParent"` // TODO: only want products with this set to false
	IsBatchedStockType       bool                   `json:"isBatchedStockType"`
	PurchasePrice            float64                `json:"PurchasePrice"`
	RetailPrice              float64                `json:"RetailPrice"`
	TaxRate                  float64                `json:"TaxRate"`
	PostalServiceID          string                 `json:"PostalServiceId"`
	CategoryID               string                 `json:"CategoryId"`
	CategoryName             string                 `json:"CategoryName"`
	PackageGroupID           string                 `json:"PackageGroupId"`
	Height                   float64                `json:"Height"`
	Width                    float64                `json:"Width"`
	Depth                    float64                `json:"Depth"`
	Weight                   float64                `json:"Weight"`
	CreationDate             time.Time              `json:"CreationDate"`
	InventoryTrackingType    int                    `json:"InventoryTrackingType"`
	BatchNumberScanRequired  bool                   `json:"BatchNumberScanRequired"`
	SerialNumberScanRequired bool                   `json:"SerialNumberScanRequired"`
	StockItemID              string                 `json:"StockItemId"`
	StockItemIntID           int                    `json:"StockItemIntId"`
}

type ProductImageResponse struct {
	Source         string `json:"Source"`
	FullSource     string `json:"FullSource"`
	PkRowID        string `json:"pkRowId"`
	IsMain         bool   `json:"IsMain"`
	SortOrder      int    `json:"SortOrder"`
	StockItemID    string `json:"StockItemId"`
	StockItemIntID int    `json:"StockItemIntId"`
}
