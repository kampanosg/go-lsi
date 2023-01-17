package linnworks

import "time"

type linnworksAuth struct {
	CustomerID            int         `json:"CustomerId"`
	FullName              string      `json:"FullName"`
	Company               string      `json:"Company"`
	ProductName           string      `json:"ProductName"`
	ExpirationDate        time.Time   `json:"ExpirationDate"`
	IsAccountHolder       bool        `json:"IsAccountHolder"`
	SessionUserID         int         `json:"SessionUserId"`
	ID                    string      `json:"Id"`
	EntityID              string      `json:"EntityId"`
	DatabaseName          string      `json:"DatabaseName"`
	DatabaseServer        interface{} `json:"DatabaseServer"`
	PrivateDatabaseServer interface{} `json:"PrivateDatabaseServer"`
	DatabaseUser          interface{} `json:"DatabaseUser"`
	DatabasePassword      interface{} `json:"DatabasePassword"`
	AppName               interface{} `json:"AppName"`
	SidRegistration       string      `json:"sid_registration"`
	UserName              string      `json:"UserName"`
	Md5Hash               string      `json:"Md5Hash"`
	Locality              string      `json:"Locality"`
	SuperAdmin            bool        `json:"SuperAdmin"`
	TTL                   int         `json:"TTL"`
	Token                 string      `json:"Token"`
	AccessToken           interface{} `json:"AccessToken"`
	GroupName             string      `json:"GroupName"`
	Device                string      `json:"Device"`
	DeviceType            string      `json:"DeviceType"`
	UserType              string      `json:"UserType"`
	Status                struct {
		State      string   `json:"State"`
		Reason     string   `json:"Reason"`
		Parameters struct{} `json:"Parameters"`
	} `json:"Status"`
	UserID     string `json:"UserId"`
	Properties struct {
	} `json:"Properties"`
	Email      string `json:"Email"`
	Server     string `json:"Server"`
	PushServer string `json:"PushServer"`
}

type LinnworksCategoryResponse struct {
	Id                string `json:"CategoryId"`
	Name              string `json:"CategoryName"`
	StructureId       int    `json:"StructureCategoryId"`
	ProductCategoryId int    `json:"ProductCategoryId"`
}

type LinnworksProductResponse struct {
	Suppliers                []interface{}                   `json:"Suppliers"`
	StockLevels              []interface{}                   `json:"StockLevels"`
	ItemChannelDescriptions  []interface{}                   `json:"ItemChannelDescriptions"`
	ItemExtendedProperties   []interface{}                   `json:"ItemExtendedProperties"`
	ItemChannelTitles        []interface{}                   `json:"ItemChannelTitles"`
	ItemChannelPrices        []interface{}                   `json:"ItemChannelPrices"`
	Images                   []LinnworksProductImageResponse `json:"Images"`
	SKU                      string                          `json:"ItemNumber"`
	ItemTitle                string                          `json:"ItemTitle"`
	BarcodeNumber            string                          `json:"BarcodeNumber"`
	MetaData                 string                          `json:"MetaData"`
	IsVariationParent        bool                            `json:"IsVariationParent"` // TODO: only want products with this set to false
	IsBatchedStockType       bool                            `json:"isBatchedStockType"`
	PurchasePrice            float64                         `json:"PurchasePrice"`
	RetailPrice              float64                         `json:"RetailPrice"`
	TaxRate                  float64                         `json:"TaxRate"`
	PostalServiceID          string                          `json:"PostalServiceId"`
	CategoryID               string                          `json:"CategoryId"`
	CategoryName             string                          `json:"CategoryName"`
	PackageGroupID           string                          `json:"PackageGroupId"`
	Height                   float64                         `json:"Height"`
	Width                    float64                         `json:"Width"`
	Depth                    float64                         `json:"Depth"`
	Weight                   float64                         `json:"Weight"`
	CreationDate             time.Time                       `json:"CreationDate"`
	InventoryTrackingType    int                             `json:"InventoryTrackingType"`
	BatchNumberScanRequired  bool                            `json:"BatchNumberScanRequired"`
	SerialNumberScanRequired bool                            `json:"SerialNumberScanRequired"`
	StockItemID              string                          `json:"StockItemId"`
	StockItemIntID           int                             `json:"StockItemIntId"`
}

type LinnworksProductImageResponse struct {
	Source         string `json:"Source"`
	FullSource     string `json:"FullSource"`
	PkRowID        string `json:"pkRowId"`
	IsMain         bool   `json:"IsMain"`
	SortOrder      int    `json:"SortOrder"`
	StockItemID    string `json:"StockItemId"`
	StockItemIntID int    `json:"StockItemIntId"`
}
