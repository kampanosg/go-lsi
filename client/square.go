package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kev/types/domain"
	"kev/types/request"
	"kev/types/response"
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
	log.Println("response Body:", string(body))

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

	log.Printf("%v", objectIds)

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
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	var squareResp response.SquareUpsertCategoryResponse
	json.Unmarshal(body, &squareResp)

	return squareResp

}
