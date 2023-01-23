package square

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kampanosg/go-lsi/types"
)

const (
	TYPE_CATEGORY  = "CATEGORY"
	TYPE_ITEM      = "ITEM"
	TYPE_VARIATION = "ITEM_VARIATION"

	VISIBILITY   = "PRIVATE"
	PRODUCT_TYPE = "REGULAR"

	PENCE_MULTIPLIER = 100
	BATCH_SIZE       = 500
	CURRENCY         = "GBP"

	VARIATION_ORDINAL = 1
	VARIATION_NAME    = "Regular"
	VARIATION_PRICING = "FIXED_PRICING"
)

type SquareClient struct {
	AccessToken string `json:"accessToken"`
	Host        string `json:"host"`
}

func NewSquareClient(accessToken, host string) *SquareClient {
	return &SquareClient{AccessToken: accessToken, Host: host}
}

func (c *SquareClient) UpsertCategories(categories []types.Category) (SquareUpsertResponse, error) {

	objects := []SquareUpsertCategoryObject{}

	for _, category := range categories {
		object := SquareUpsertCategoryObject{
			Id:        category.SquareId,
			Type:      TYPE_CATEGORY,
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

	url := fmt.Sprintf("%s/catalog/batch-upsert", c.Host)
	jsonReq, _ := json.Marshal(batchRequest)
	headers := make(map[string]string)

	headers["Square-Version"] = "2022-12-14"
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)

	var squareResp SquareUpsertResponse

	resp, err := makeRequest(POST, url, headers, jsonReq)
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

		if len(objects) >= BATCH_SIZE {
			currentBatch.Objects = objects
			batches = append(batches, currentBatch)

			currentBatch = SquareProductBatch{}
			objects = make([]SquareProductObject, 0)
		}

		itemMoney := SquarePriceMoney{
			Amount:   int(product.Price * PENCE_MULTIPLIER),
			Currency: CURRENCY,
		}

		variationData := SquareProductVariationData{
			ItemID:      product.SquareId,
			Sku:         product.SKU,
			Upc:         product.Barcode,
			Name:        VARIATION_NAME,
			PricingType: VARIATION_PRICING,
			Ordinal:     VARIATION_ORDINAL,
			PriceMoney:  itemMoney,
		}

		itemVariations := []SquareProductVariation{
			{
				Type:                  TYPE_VARIATION,
				ID:                    product.SquareVarId,
				IsDeleted:             false,
				PresentAtAllLocations: true,
				Version:               product.Version,
				ItemVariationData:     variationData,
			},
		}

		itemData := SquareProductData{
			Name:               product.Title,
			CategoryID:         product.SquareCategoryId,
			Visibility:         VISIBILITY,
			ProductType:        PRODUCT_TYPE,
			SkipModifierScreen: false,
			IsTaxable:          true,
			Variations:         itemVariations,
		}

		object := SquareProductObject{
			Type:                  TYPE_ITEM,
			ID:                    product.SquareId,
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

	url := fmt.Sprintf("%s/catalog/batch-upsert", c.Host)
	jsonReq, _ := json.Marshal(batchRequest)
	headers := make(map[string]string, 0)

	headers["Square-Version"] = "2022-12-14"
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)

	var squareResp SquareUpsertResponse
	resp, err := makeRequest(POST, url, headers, jsonReq)
	if err != nil {
		return squareResp, err
	}

	err = json.Unmarshal(resp, &squareResp)

	return squareResp, err
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
	headers["Square-Version"] = "2022-12-14"
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)

	for _, br := range batchRequests {
		jsonReq, _ := json.Marshal(br)
		makeRequest(POST, url, headers, jsonReq)
		if _, err := makeRequest(POST, url, headers, jsonReq); err != nil {
			return err
		}

	}
	return nil
}

func (c *SquareClient) SearchOrders(start time.Time, end time.Time) (SquareOrderSearchResponse, error) {
	objects := []SquareUpsertCategoryObject{}

	for _, category := range categories {
		object := SquareUpsertCategoryObject{
			Id:        category.SquareId,
			Type:      TYPE_CATEGORY,
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

	url := fmt.Sprintf("%s/catalog/batch-upsert", c.Host)
	jsonReq, _ := json.Marshal(batchRequest)
	headers := make(map[string]string)

	headers["Square-Version"] = "2022-12-14"
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", c.AccessToken)

	var squareResp SquareUpsertResponse

	resp, err := makeRequest(POST, url, headers, jsonReq)
	if err != nil {
		return squareResp, err
	}

	err = json.Unmarshal(resp, &squareResp)

	return squareResp, err
}
