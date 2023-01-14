package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kev/transform"
	"kev/types/domain"
	"kev/types/response"
	"log"
	"net/http"
	"strings"
)

type LinnworksClient struct {
	id     string
	secret string
	token  string
	auth   domain.Auth
}

func NewLinnworksClient(id, secret, token string) *LinnworksClient {
	return &LinnworksClient{
		id:     id,
		secret: secret,
		token:  token,
	}
}

func (c *LinnworksClient) GetCategories() ([]domain.Category, error) {
	c.refreshToken()

	categories := []domain.Category{}

	url := "https://eu-ext.linnworks.net/api/Inventory/GetCategories"
	method := "POST"

	client := &http.Client{}
	payload := strings.NewReader("=")
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Printf("setup failed, reason=%v\n", err)
		return categories, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", c.auth.Token)

	res, err := client.Do(req)
	if err != nil {
		log.Printf("request failed with reason: %v\n", err)
		return categories, err
	}

	defer res.Body.Close()

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("unable to process response, reason: %v\n", err)
		return categories, err
	}

	var authResp []response.CategoryResponse
	json.Unmarshal(responseData, &authResp)

	return transform.FromCategoriesRespToDomain(authResp), nil
}

func (c *LinnworksClient) GetProducts() ([]domain.Product, error) {
	c.refreshToken()
	products := []domain.Product{}
	entriesPerPage := 200 // TODO: move this to config
	pageNumber := 1

	url := "https://eu-ext.linnworks.net/api/Stock/GetStockItemsFull"
	method := "POST"

	var builder strings.Builder
	builder.WriteString("loadCompositeParents=True")
	builder.WriteString("&loadVariationParents=False")
	builder.WriteString("&dataRequirements=%5B1%2C8%5D&searchTypes=%5B0%2C1%2C2%5D")
	builder.WriteString(fmt.Sprintf("&entriesPerPage=%d", entriesPerPage))

	for {
		pld := fmt.Sprintf("%s&pageNumber=%d", builder.String(), pageNumber)
		payload := strings.NewReader(pld)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			log.Printf("setup failed, reason=%v", err)
			return products, err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", c.auth.Token)

		res, err := client.Do(req)
		if err != nil {
			log.Printf("request failed, reason=%v", err)
			return products, err
		}
		defer res.Body.Close()

		responseData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("unable to process response, reason: %v\n", err)
			return products, err
		}

		var authResp []response.ProductResponse
		json.Unmarshal(responseData, &authResp)
		for _, p := range transform.FromProductsRespToDomain(authResp) {
			products = append(products, p)
		}

		pageNumber += 1

		if len(responseData) < entriesPerPage {
			break
		}
	}
	return products, nil

}

func (c *LinnworksClient) refreshToken() {

	// TODO: efficiency - check if current token has expired
	url := "https://api.linnworks.net/api/Auth/AuthorizeByApplication"
	method := "POST"

	body := fmt.Sprintf("applicationId=%s&applicationSecret=%s&token=%s", c.id, c.secret, c.token)
	payload := strings.NewReader(body)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var authResp response.Auth
	json.Unmarshal(responseData, &authResp)

	c.auth = transform.FromAuthResponseToDomain(authResp)
}
