package linnworks

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kampanosg/go-lsi/types"
)

const (
	LinnworksServer1 = "https://api.linnworks.net/api/"
	LinnworksServer2 = "https://eu-ext.linnworks.net/api/"
	DefaultDryRun    = true
)

type LinnworksClient struct {
	Id     string
	Secret string
	Token  string
	DryRun bool
	auth   linnworksAuth
}

func NewLinnworksClient(id, secret, token string) *LinnworksClient {
	return &LinnworksClient{
		Id:     id,
		Secret: secret,
		Token:  token,
		DryRun: DefaultDryRun,
	}
}

func (c *LinnworksClient) GetCategories() ([]LinnworksCategoryResponse, error) {
	c.refreshToken()

	url := fmt.Sprintf("%s/Inventory/GetCategories", LinnworksServer2)
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

	url := fmt.Sprintf("%s/Stock/GetStockItemsFull", LinnworksServer2)
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

func (c *LinnworksClient) CreateOrders(orders []types.Order) (LinnworksCreateOrdersResponse, error) {
	c.refreshToken()

	url := fmt.Sprintf("%s/Orders/CreateOrders", LinnworksServer2)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = c.auth.Token

	for _, order := range orders {

		pld := fmt.Sprintf(orderTemplate, "irieiireii")
        fmt.Println(pld)
        fmt.Println(order)
        
		payload := strings.NewReader(pld)

		if c.DryRun {
		} else {
			resp, err := makeRequest(POST, url, payload, headers)
			if err != nil {
				return LinnworksCreateOrdersResponse{}, err
			}
            var productResps []LinnworksProductResponse
            json.Unmarshal(resp, &productResps)
		}
	}

	return LinnworksCreateOrdersResponse{}, nil
}

func (c *LinnworksClient) refreshToken() {
	if c.auth.ExpirationDate.After(time.Now()) {
		log.Printf("lw: token has not expired, no need to refresh\n")
		return
	}

	url := fmt.Sprintf("%s/Auth/AuthorizeByApplication", LinnworksServer1)
	body := fmt.Sprintf("applicationId=%s&applicationSecret=%s&token=%s", c.Id, c.Secret, c.Token)
	payload := strings.NewReader(body)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	response, _ := makeRequest(POST, url, payload, headers)

	var authResp linnworksAuth
	json.Unmarshal(response, &authResp)

	c.auth = authResp
}
