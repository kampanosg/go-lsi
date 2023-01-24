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

type LinnworksCreateOrdersRequest struct {
	Orders []LinnworksOrder `json:"orders"`
}
type LinnworksOrderItem struct {
	TaxCostInclusive bool          `json:"TaxCostInclusive"`
	UseChannelTax    bool          `json:"UseChannelTax"`
	PricePerUnit     float64       `json:"PricePerUnit"`
	Qty              int           `json:"Qty"`
	TaxRate          float64       `json:"TaxRate"`
	LineDiscount     float64       `json:"LineDiscount"`
	ItemNumber       string        `json:"ItemNumber"`
	ChannelSKU       string        `json:"ChannelSKU"`
	IsService        bool          `json:"IsService"`
	ItemTitle        string        `json:"ItemTitle"`
	Options          []interface{} `json:"Options"`
}
type LinnworksBillingAddress struct {
	MatchCountryCode string `json:"MatchCountryCode"`
	MatchCountryName string `json:"MatchCountryName"`
	FullName         string `json:"FullName"`
	Company          string `json:"Company"`
	Address1         string `json:"Address1"`
	Address2         string `json:"Address2"`
	Address3         string `json:"Address3"`
	Town             string `json:"Town"`
	Region           string `json:"Region"`
	PostCode         string `json:"PostCode"`
	Country          string `json:"Country"`
	PhoneNumber      string `json:"PhoneNumber"`
	EmailAddress     string `json:"EmailAddress"`
	IsEmpty          bool   `json:"isEmpty"`
}
type LinnworksDeliveryAddress struct {
	MatchCountryCode string `json:"MatchCountryCode"`
	MatchCountryName string `json:"MatchCountryName"`
	FullName         string `json:"FullName"`
	Company          string `json:"Company"`
	Address1         string `json:"Address1"`
	Address2         string `json:"Address2"`
	Address3         string `json:"Address3"`
	Town             string `json:"Town"`
	Region           string `json:"Region"`
	PostCode         string `json:"PostCode"`
	Country          string `json:"Country"`
	PhoneNumber      string `json:"PhoneNumber"`
	EmailAddress     string `json:"EmailAddress"`
	IsEmpty          bool   `json:"isEmpty"`
}
type LinnworksOrder struct {
	UseChannelTax               bool                     `json:"UseChannelTax"`
	PkOrderID                   string                   `json:"pkOrderId"`
	AutomaticallyLinkBySKU      bool                     `json:"AutomaticallyLinkBySKU"`
	Site                        string                   `json:"Site"`
	MatchPostalServiceTag       string                   `json:"MatchPostalServiceTag"`
	PostalServiceName           string                   `json:"PostalServiceName"`
	SavePostalServiceIfNotExist bool                     `json:"SavePostalServiceIfNotExist"`
	MatchPaymentMethodTag       string                   `json:"MatchPaymentMethodTag"`
	PaymentMethodName           string                   `json:"PaymentMethodName"`
	SavePaymentMethodIfNotExist bool                     `json:"SavePaymentMethodIfNotExist"`
	MappingSource               string                   `json:"MappingSource"`
	OrderState                  string                   `json:"OrderState"`
	PaymentStatus               string                   `json:"PaymentStatus"`
	OrderItems                  []LinnworksOrderItem     `json:"OrderItems"`
	ExtendedProperties          []interface{}            `json:"ExtendedProperties"`
	Notes                       []interface{}            `json:"Notes"`
	Source                      string                   `json:"Source"`
	SubSource                   string                   `json:"SubSource"`
	ChannelBuyerName            string                   `json:"ChannelBuyerName"`
	ReferenceNumber             string                   `json:"ReferenceNumber"`
	ExternalReference           string                   `json:"ExternalReference"`
	SecondaryReferenceNumber    string                   `json:"SecondaryReferenceNumber"`
	Currency                    string                   `json:"Currency"`
	ConversionRate              float64                  `json:"ConversionRate"`
	ReceivedDate                time.Time                `json:"ReceivedDate"`
	DispatchBy                  time.Time                `json:"DispatchBy"`
	PaidOn                      time.Time                `json:"PaidOn"`
	PostalServiceCost           float64                  `json:"PostalServiceCost"`
	PostalServiceTaxRate        float64                  `json:"PostalServiceTaxRate"`
	PostalServiceDiscount       float64                  `json:"PostalServiceDiscount"`
	Discount                    float64                  `json:"Discount"`
	DiscountType                string                   `json:"DiscountType"`
	DiscountTaxType             string                   `json:"DiscountTaxType"`
	BillingAddress              LinnworksBillingAddress  `json:"BillingAddress"`
	DeliveryAddress             LinnworksDeliveryAddress `json:"DeliveryAddress"`
	DeliveryStartDate           time.Time                `json:"DeliveryStartDate"`
	DeliveryEndDate             time.Time                `json:"DeliveryEndDate"`
	OrderIdentifierTags         []string                 `json:"OrderIdentifierTags"`
	ForceReSaveFulfilledOrder   bool                     `json:"ForceReSaveFulfilledOrder"`
}

type LinnworksCreateOrdersResponse []string
