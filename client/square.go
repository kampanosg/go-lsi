package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kampanosg/go-lsi/types/domain"
	"github.com/kampanosg/go-lsi/types/request"
	"github.com/kampanosg/go-lsi/types/response"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type SquareClient struct {
	AccessToken string `json:"accessToken"`
	Host        string `json:"host"`
}

func NewSquareClient(accessToken, host string) SquareClient {
	return SquareClient{AccessToken: accessToken, Host: host}
}

func (c *SquareClient) UpsertCategories(categories []domain.Category) response.SquareUpsertCategoryResponse {
	objects := []request.SquareUpsertCategoryRequest{}

	for _, category := range categories {
		categoryRequest := request.SquareUpsertCategoryRequest{
			Id:        category.SquareId,
			Type:      "CATEGORY",
			IsDeleted: false,
			Version:   category.Version,
			CategoryData: request.CategoryData{
				Name: category.Name,
			},
		}

		objects = append(objects, categoryRequest)
	}

	squareBatch := request.SquareBatch{
		Objects: objects,
	}

	batchRequest := request.SquareBatchUpsertCategoryRequest{
		IdempotencyKey: uuid.New().String(),
		Batches:        []request.SquareBatch{squareBatch},
	}

	jsonReq, _ := json.Marshal(batchRequest)

	// log.Printf("upsert req: %s", string(jsonReq))

	// panic("fhfhfhfhfh")

	url := fmt.Sprintf("%s/catalog/batch-upsert", c.Host)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	req.Header.Set("Square-Version", "2022-12-14")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// log.Println("response Body:", string(body))

	var squareResp response.SquareUpsertCategoryResponse
	json.Unmarshal(body, &squareResp)

	return squareResp

}

func (c *SquareClient) DeleteCategories(categories []domain.Category) response.SquareUpsertCategoryResponse {

	batchRequest := request.BatchDeleteCategoriesRequest{}
	objectIds := []string{}

	for _, category := range categories {
		objectIds = append(objectIds, category.SquareId)
	}

	batchRequest.ObjectIds = objectIds

	jsonReq, _ := json.Marshal(batchRequest)

	url := fmt.Sprintf("%s/catalog/batch-delete", c.Host)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	req.Header.Set("Square-Version", "2022-12-14")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	// log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// log.Println("response Body:", string(body))

	var squareResp response.SquareUpsertCategoryResponse
	json.Unmarshal(body, &squareResp)

	return squareResp

}

func (c *SquareClient) UpsertProducts(products []domain.Product) response.SquareUpsertItemResponse {

	objects := make([]request.ItemObject, 0)
	batches := make([]request.ItemBatch, 0)
	currentBatch := request.ItemBatch{}

	for _, product := range products {

		if len(objects) >= 500 { // TODO: Move len to config
			currentBatch.Objects = objects
			batches = append(batches, currentBatch)

			currentBatch = request.ItemBatch{}
			objects = make([]request.ItemObject, 0)
		}

		itemMoney := request.PriceMoney{
			Amount:   int(product.Price * 100),
			Currency: "GBP",
		}

		variationData := request.ItemVariationData{
			ItemID:      product.SquareId,
			Sku:         product.SKU,
			Name:        "Regular",
			PricingType: "FIXED_PRICING",
			Ordinal:     1,
			PriceMoney:  itemMoney,
		}

		itemVariations := []request.ItemVariation{
			{
				Type:                  "ITEM_VARIATION",
				ID:                    product.SquareVarId,
				IsDeleted:             false,
				PresentAtAllLocations: true,
				Version:               product.Version,
				ItemVariationData:     variationData,
			},
		}

		itemData := request.ItemData{
			Name:               product.Title,
			Visibility:         "PRIVATE",
			CategoryID:         product.SquareCategoryId,
			ProductType:        "REGULAR",
			SkipModifierScreen: false,
			IsTaxable:          true,
			Variations:         itemVariations,
		}

		object := request.ItemObject{
			Type:                  "ITEM",
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

	batchRequest := request.SquareBatchUpsertItemRequest{
		IdempotencyKey: uuid.New().String(),
		Batches:        batches,
	}

	jsonReq, _ := json.Marshal(batchRequest)

	// log.Printf("upsert req: %s", string(jsonReq))
	// panic("dadasadadsads")

	url := fmt.Sprintf("%s/catalog/batch-upsert", c.Host)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	req.Header.Set("Square-Version", "2022-12-14")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	// log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// log.Println("response Body:", string(body))

	var squareResp response.SquareUpsertItemResponse
	json.Unmarshal(body, &squareResp)

	return squareResp
}

func (c *SquareClient) DeleteProducts(products []domain.Product) {

	batchRequest := request.BatchDeleteCategoriesRequest{}
	objectIds := []string{}

	batchRequests := make([]request.BatchDeleteCategoriesRequest, 0)

	for _, product := range products {
		objectIds = append(objectIds, product.SquareId)
		if len(objectIds) == 200 {
			batchRequest.ObjectIds = objectIds
			batchRequests = append(batchRequests, batchRequest)

			batchRequest = request.BatchDeleteCategoriesRequest{}
			objectIds = make([]string, 0)
		}
	}

	batchRequest.ObjectIds = objectIds
	batchRequests = append(batchRequests, batchRequest)

	for _, br := range batchRequests {

		jsonReq, _ := json.Marshal(br)

		url := fmt.Sprintf("%s/catalog/batch-delete", c.Host)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
		req.Header.Set("Square-Version", "2022-12-14")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response Body:", string(body))

	}

}
