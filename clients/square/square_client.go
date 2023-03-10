package square

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

const (
	TypeCategory  = "CATEGORY"
	TypeItem      = "ITEM"
	TypeVariation = "ITEM_VARIATION"
	TypeService   = "APPOINTMENTS_SERVICE"

	Visibility  = "PRIVATE"
	ProductType = "REGULAR"

	PenceMultiplier = 100
	BatchSize       = 700
	Currency        = "GBP"

	VariationOrdinal = 1
	VariationName    = "Regular"
	VariationPricing = "FIXED_PRICING"

	ServiceSkuSuffix = "GTR-"

	OrderLimit = 50

	DefaultSleepDuration = 90 * time.Second
)

type SQ interface {
	GetItemVersion(squareId string) (int64, int64, error)
	UpsertCategories(categories []types.Category) (SquareUpsertResponse, error)
	UpsertProducts(products []types.Product) (SquareUpsertResponse, error)
	BatchDeleteItems(itemIds []string) error
	SearchOrders(start time.Time, end time.Time) ([]SquareOrder, error)
}

type SquareClient struct {
	AccessToken   string
	Host          string
	ApiVersion    string
	LocationId    string
	TeamMemberIds []string
	logger        *zap.SugaredLogger
}

func NewSquareClient(accessToken, host, version, location string, teamMembers []string, logger *zap.SugaredLogger) *SquareClient {
	return &SquareClient{
		AccessToken:   accessToken,
		Host:          host,
		ApiVersion:    version,
		LocationId:    location,
		TeamMemberIds: teamMembers,
		logger:        logger,
	}
}

func (c *SquareClient) GetItemVersion(squareId string) (int64, int64, error) {
	headers := make(map[string]string)
	headers["Square-Version"] = c.ApiVersion
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)

	url := fmt.Sprintf("%s/catalog/object/%s", c.Host, squareId)
	resp, err := c.makeRequest("GET", url, headers, []byte{})
	if err != nil {
		return 0, 0, err
	}

	var r SquareCatalogItemResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		c.logger.Errorw("failed to parse object", "error", err.Error())
		return 0, 0, err
	}

	itemVersion := r.Object.Version
	var itemVarVersion int64

	if len(r.Object.ItemData.Variations) > 0 {
		itemVarVersion = r.Object.ItemData.Variations[0].Version
	}

	return itemVersion, itemVarVersion, nil
}

func (c *SquareClient) UpsertCategories(categories []types.Category) (SquareUpsertResponse, error) {

	objects := []SquareUpsertCategoryObject{}

	for _, category := range categories {
		object := SquareUpsertCategoryObject{
			Id:        category.SquareID,
			Type:      TypeCategory,
			IsDeleted: false,
			Version:   category.Version,
			CategoryData: SquareCategoryData{
				Name: category.Name,
			},
		}
		objects = append(objects, object)
	}

	squareBatch := SquareCategoryBatch{
		Objects: objects,
	}

	batchRequest := SquareBatchUpsertCatalogItemRequest{
		IdempotencyKey: uuid.New().String(),
		Batches:        []SquareCategoryBatch{squareBatch},
	}

	var squareResp SquareUpsertResponse

	jsonReq, err := json.Marshal(batchRequest)
	if err != nil {
		c.logger.Errorw("unable to marshall request", "error", err)
		return squareResp, err
	}

	url := fmt.Sprintf("%s/catalog/batch-upsert", c.Host)
	headers := make(map[string]string)
	headers["Square-Version"] = c.ApiVersion
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)

	c.logger.Debugw("attempting to call square", "url", url, "req", string(jsonReq))

	resp, err := c.makeRequest(POST, url, headers, jsonReq)
	if err != nil {
		return squareResp, err
	}

	err = json.Unmarshal(resp, &squareResp)

	return squareResp, err
}

func (c *SquareClient) UpsertProducts(products []types.Product) (SquareUpsertResponse, error) {

	objects := make([]SquareProductObject, 0)
	batches := make([]SquareProductBatch, 0)
	currentBatch := SquareProductBatch{}

	for _, product := range products {

		if len(objects) >= BatchSize {
			currentBatch.Objects = objects
			batches = append(batches, currentBatch)

			currentBatch = SquareProductBatch{}
			objects = make([]SquareProductObject, 0)
		}

		serviceDuration := c.getServiceDuration(product)
		availableForBooking := false
		teamMemberIds := make([]string, 0)

		if strings.HasPrefix(product.SKU, ServiceSkuSuffix) {
			availableForBooking = true
			teamMemberIds = c.TeamMemberIds
		}

		itemMoney := SquarePriceMoney{
			Amount:   int(product.Price * PenceMultiplier),
			Currency: Currency,
		}

		variationData := SquareProductVariationData{
			ItemID:              product.SquareID,
			Sku:                 product.SKU,
			Upc:                 product.Barcode,
			Name:                VariationName,
			PricingType:         VariationPricing,
			Ordinal:             VariationOrdinal,
			PriceMoney:          itemMoney,
			ServiceDuration:     serviceDuration,
			AvailableForBooking: availableForBooking,
			TeamMemberIds:       teamMemberIds,
		}

		itemVariations := []SquareProductVariation{
			{
				Type:                  TypeVariation,
				ID:                    product.SquareVarID,
				IsDeleted:             false,
				PresentAtAllLocations: true,
				Version:               product.VariationVersion,
				ItemVariationData:     variationData,
			},
		}

		productType := ProductType
		if strings.HasPrefix(product.SKU, "GTR-") {
			productType = "APPOINTMENTS_SERVICE"
		}

		itemData := SquareProductData{
			Name:               product.Title,
			CategoryID:         product.SquareCategoryID,
			Visibility:         Visibility,
			ProductType:        productType,
			SkipModifierScreen: false,
			IsTaxable:          true,
			Variations:         itemVariations,
		}

		object := SquareProductObject{
			Type:                  TypeItem,
			ID:                    product.SquareID,
			IsDeleted:             false,
			PresentAtAllLocations: true,
			Version:               product.Version,
			ItemData:              itemData,
		}

		objects = append(objects, object)
	}

	if len(objects) > 0 {
		currentBatch.Objects = objects
		batches = append(batches, currentBatch)
	}

	batchRequest := SquareBatchUpsertProductRequest{
		IdempotencyKey: uuid.New().String(),
		Batches:        batches,
	}

	var squareResp SquareUpsertResponse

	jsonReq, err := json.Marshal(batchRequest)
	if err != nil {
		c.logger.Errorw("unable to marshall request", "error", err)
		return squareResp, err
	}

	url := fmt.Sprintf("%s/catalog/batch-upsert", c.Host)
	headers := make(map[string]string, 0)
	headers["Square-Version"] = c.ApiVersion
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)

	c.logger.Debugw("attempting to call square", "url", url, "req", string(jsonReq))

	resp, err := c.makeRequest(POST, url, headers, jsonReq)
	if err != nil {
		return squareResp, err
	}

	if err := json.Unmarshal(resp, &squareResp); err != nil {
		c.logger.Errorw("unable to unmarshall request", "error", err)
		return squareResp, err
	}

	return squareResp, nil
}

func (c *SquareClient) BatchDeleteItems(itemIds []string) error {

	batchRequest := BatchDeleteItemsRequest{}
	objectIds := []string{}

	batchRequests := make([]BatchDeleteItemsRequest, 0)

	for _, itemId := range itemIds {
		objectIds = append(objectIds, itemId)
		if len(objectIds) == 200 {
			batchRequest.ObjectIds = objectIds
			batchRequests = append(batchRequests, batchRequest)

			batchRequest = BatchDeleteItemsRequest{}
			objectIds = make([]string, 0)
		}
	}

	batchRequest.ObjectIds = objectIds
	batchRequests = append(batchRequests, batchRequest)

	url := fmt.Sprintf("%s/catalog/batch-delete", c.Host)
	headers := make(map[string]string, 0)
	headers["Square-Version"] = c.ApiVersion
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)

	for index, br := range batchRequests {
		jsonReq, err := json.Marshal(br)
		if err != nil {
			c.logger.Errorw("unable to marshall request", "error", err)
			return err
		}

		c.logger.Debugw("attempting to call square", "url", url, "req", string(jsonReq))

		if _, err := c.makeRequest(POST, url, headers, jsonReq); err != nil {
			return err
		}

		if index >= len(batchRequests)-1 {
			break
		}

		c.logger.Infow("will sleep to avoid rate limiting", "duration", DefaultSleepDuration)
		time.Sleep(DefaultSleepDuration)
	}
	return nil
}

func (c *SquareClient) SearchOrders(start time.Time, end time.Time) ([]SquareOrder, error) {

	url := fmt.Sprintf("%s/orders/search", c.Host)
	cursor := ""

	headers := make(map[string]string)
	headers["Square-Version"] = c.ApiVersion
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)

	orders := make([]SquareOrder, 0)

	for {
		searchRequest := SquareSearchOrdersRequest{
			ReturnEntries: false,
			Limit:         OrderLimit,
			Query: SquareQuery{
				Filter: SquareFilter{
					DateTimeFilter: SquareDateTimeFilter{
						CreatedAt: SquareDateRange{
							StartAt: start,
							EndAt:   end,
						},
					},
				},
			},
			LocationIds: []string{c.LocationId},
			Cursor:      cursor,
		}

		var squareResp SquareOrderSearchResponse

		jsonReq, err := json.Marshal(searchRequest)
		if err != nil {
			c.logger.Errorw("unable to marshall request", "error", err)
			return make([]SquareOrder, 0), err
		}

		c.logger.Debugw("attempting to call square", "url", url, "req", string(jsonReq))

		resp, err := c.makeRequest(POST, url, headers, jsonReq)
		if err != nil {
			return orders, err
		}

		if err := json.Unmarshal(resp, &squareResp); err != nil {
			c.logger.Errorw("cannot unmarshall resp", "error", err.Error())
			return orders, err
		}

		if len(squareResp.Orders) == 0 {
			break
		}

		orders = append(orders, squareResp.Orders...)

		cursor = squareResp.Cursor
		if cursor == "" {
			break
		}
	}

	return orders, nil
}

func (c *SquareClient) getServiceDuration(product types.Product) int {
	var duration int
	switch product.SKU {
	case "GTR-001":
		duration = 40
	case "GTR-002":
		duration = 40
	case "GTR-003":
		duration = 30
	case "GTR-004":
		duration = 30
	case "GTR-005":
		duration = 40
	case "GTR-006":
		duration = 40
	case "GTR-007":
		duration = 30
	case "GTR-008":
		duration = 40
	case "GTR-009":
		duration = 40
	case "GTR-010":
		duration = 40
	default:
		return 0
	}
	durationMillis := duration * 60 * 1000
	c.logger.Debugw("will retrieve service duration", "sku", product.SKU, "durationMins", duration, "durationMillis", durationMillis)
	return durationMillis
}
