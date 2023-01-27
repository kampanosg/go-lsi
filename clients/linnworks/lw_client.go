package linnworks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kampanosg/go-lsi/types"
)

const (
	LinnworksServer1 = "https://api.linnworks.net/api/"
	LinnworksServer2 = "https://eu-ext.linnworks.net/api/"
	DefaultDryRun    = false
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

	response, err := makeRequest(Post, url, payload, headers)
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

		resp, err := makeRequest(Post, url, payload, headers)
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

	linnworksUrl := fmt.Sprintf("%s/Orders/CreateOrders", LinnworksServer2)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = c.auth.Token

	for _, order := range orders {

		var orderProducts bytes.Buffer
		orderProducts.WriteString("[")
		for index, product := range order.Products {
			p := fmt.Sprintf(orderItemTemplate,
				product.PricePerUnit,
				product.Quantity,
				product.ItemNumber,
				product.SKU,
				"Test",
			)
			orderProducts.WriteString(p)

			if index < len(order.Products)-1 {
				orderProducts.WriteString(",")
			}
		}
		orderProducts.WriteString("]")

		formattedTime := order.CreatedAt.Format("2006-01-02T15:04:05.000000+01:00")

		pld := fmt.Sprintf(orderTemplate,
			uuid.New().String(),
			orderProducts.String(),
			order.SquareId,
			order.SquareId,
			order.SquareId,
			formattedTime,
			formattedTime,
			formattedTime,
			order.SquareId,
			order.SquareId,
			order.SquareId,
		)

		encodedPld := url.QueryEscape(pld)
		f := fmt.Sprintf("orders=%s&location=Default", encodedPld)
		payload := strings.NewReader(f)

		log.Printf("\n")
		log.Printf("\n")
		log.Printf("payload = %v\n", payload)
		if c.DryRun {
		} else {
			resp, err := makeRequest(Post, linnworksUrl, payload, headers)
			if err != nil {
				return LinnworksCreateOrdersResponse{}, err
			}
			fmt.Printf("resp: %v\n", string(resp))
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

	response, _ := makeRequest(Post, url, payload, headers)

	var authResp linnworksAuth
	json.Unmarshal(response, &authResp)

	c.auth = authResp
}

func formatDatePart(part int) string {
	if part < 10 {
		return fmt.Sprintf("0%d", part)
	}
	return fmt.Sprintf("%d", part)
}
