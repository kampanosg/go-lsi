package linnworks

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	LW_SERVER_1 = "https://api.linnworks.net/api/"
	LW_SERVER_2 = "https://eu-ext.linnworks.net/api/"
)

type LinnworksClient struct {
	Id     string
	Secret string
	Token  string
	auth   linnworksAuth
}

func NewLinnworksClient(id, secret, token string) *LinnworksClient {
	return &LinnworksClient{
		Id:     id,
		Secret: secret,
		Token:  token,
	}
}

func (c *LinnworksClient) GetCategories() ([]LinnworksCategoryResponse, error) {
	c.refreshToken()

	url := fmt.Sprintf("%s/Inventory/GetCategories", LW_SERVER_2)
	payload := strings.NewReader("=")
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = c.auth.Token

	response, err := makeRequest(POST, url, payload, headers)
	if err != nil {
		return []LinnworksCategoryResponse{}, err
	}

	var categoriesResps []LinnworksCategoryResponse
	json.Unmarshal(response, &categoriesResps)

	return categoriesResps, nil
}

func (c *LinnworksClient) GetProducts() ([]LinnworksProductResponse, error) {
	c.refreshToken()

	entriesPerPage := 200
	pageNumber := 1

	url := fmt.Sprintf("%s/Stock/GetStockItemsFull", LW_SERVER_2)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = c.auth.Token

	var builder strings.Builder
	builder.WriteString("loadCompositeParents=True")
	builder.WriteString("&loadVariationParents=False")
	builder.WriteString("&dataRequirements=%5B1%2C8%5D&searchTypes=%5B0%2C1%2C2%5D")
	builder.WriteString(fmt.Sprintf("&entriesPerPage=%d", entriesPerPage))

	products := make([]LinnworksProductResponse, 0)

	for {
		pld := fmt.Sprintf("%s&pageNumber=%d", builder.String(), pageNumber)
		payload := strings.NewReader(pld)

		resp, err := makeRequest(POST, url, payload, headers)
		if err != nil {
			return products, err
		}

		var productResps []LinnworksProductResponse
		json.Unmarshal(resp, &productResps)

		for _, product := range productResps {
			if !product.IsBatchedStockType {
				if product.ItemTitle != "" && product.RetailPrice > 0 {
					products = append(products, product)
				}
			}
		}

		pageNumber += 1

		if len(productResps) < entriesPerPage {
			break
		}
	}
	return products, nil
}

func (c *LinnworksClient) refreshToken() {
	if c.auth.ExpirationDate.After(time.Now()) {
		log.Printf("lw: token has not expired, no need to refresh\n")
		return
	}

	url := fmt.Sprintf("%s/Auth/AuthorizeByApplication", LW_SERVER_1)
	body := fmt.Sprintf("applicationId=%s&applicationSecret=%s&token=%s", c.Id, c.Secret, c.Token)
	payload := strings.NewReader(body)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	response, _ := makeRequest(POST, url, payload, headers)

	var authResp linnworksAuth
	json.Unmarshal(response, &authResp)

	c.auth = authResp
}
