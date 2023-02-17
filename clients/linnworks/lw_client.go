package linnworks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

const (
	LinnworksServer1 = "https://api.linnworks.net/api/"
	LinnworksServer2 = "https://eu-ext.linnworks.net/api/"
)

type LW interface {
	GetCategories() ([]LinnworksCategoryResponse, error)
	GetProducts() ([]LinnworksProductResponse, error)
	CreateOrders(orders []types.Order) (LinnworksCreateOrdersResponse, error)
}

type LinnworksClient struct {
	Id     string
	Secret string
	Token  string
	auth   linnworksAuth
	logger *zap.SugaredLogger
}

func NewLinnworksClient(id, secret, token string, logger *zap.SugaredLogger) *LinnworksClient {
	return &LinnworksClient{
		Id:     id,
		Secret: secret,
		Token:  token,
		logger: logger,
	}
}

func (c *LinnworksClient) GetCategories() ([]LinnworksCategoryResponse, error) {
	if err := c.refreshToken(); err != nil {
		c.logger.Errorw("unable to refresh linnworks auth token")
		return []LinnworksCategoryResponse{}, err
	}

	url := fmt.Sprintf("%s/Inventory/GetCategories", LinnworksServer2)
	payload := strings.NewReader("=")
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = c.auth.Token

	c.logger.Debugw("attempting to call linnworks", "url", url)

	response, err := c.makeRequest(Post, url, payload, headers)
	if err != nil {
		return []LinnworksCategoryResponse{}, err
	}

	var categoriesResps []LinnworksCategoryResponse
	json.Unmarshal(response, &categoriesResps)

	return categoriesResps, nil
}

func (c *LinnworksClient) GetProducts() ([]LinnworksProductResponse, error) {
	if err := c.refreshToken(); err != nil {
		c.logger.Errorw("unable to refresh linnworks auth token")
		return []LinnworksProductResponse{}, err
	}

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

		c.logger.Debugw("attempting to call linnworks", "url", url, "payload", pld)

		resp, err := c.makeRequest(Post, url, payload, headers)
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
	if err := c.refreshToken(); err != nil {
		c.logger.Errorw("unable to refresh linnworks auth token")
		return LinnworksCreateOrdersResponse{}, err
	}

	linnworksUrl := fmt.Sprintf("%s/Orders/CreateOrders", LinnworksServer2)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = c.auth.Token

	ordersResp := make(LinnworksCreateOrdersResponse, 0)

	for _, order := range orders {

		var orderProducts bytes.Buffer
		orderProducts.WriteString("[")
		for index, product := range order.Products {
			p := fmt.Sprintf(orderItemTemplate,
				product.PricePerUnit,
				product.Quantity,
				product.ItemNumber,
				product.SKU,
				"Prod",
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
			order.SquareID,
			order.SquareID,
			order.SquareID,
			formattedTime,
			formattedTime,
			formattedTime,
			order.SquareID,
			order.SquareID,
			order.SquareID,
		)

		encodedPld := url.QueryEscape(pld)
		encodedOrder := fmt.Sprintf("orders=%s&location=Default", encodedPld)
		payload := strings.NewReader(encodedOrder)

		c.logger.Debugw("attempting to call linnworks", "url", linnworksUrl, "payload", encodedOrder)

		resp, err := c.makeRequest(Post, linnworksUrl, payload, headers)
		if err != nil {
			// return ordersResp, err
			c.logger.Debugw("http client threw error", "error", err)
			continue
		}

		c.logger.Debugw("linnworks finished processing new order", "response", string(resp))
		var productResps LinnworksCreateOrdersResponse
		if err := json.Unmarshal(resp, &productResps); err != nil {
			// return ordersResp, err
			c.logger.Debugw("json parser could not parse resp", "error", err)
			continue
		}

		ordersResp = append(ordersResp, productResps...)
	}

	return ordersResp, nil
}

func (c *LinnworksClient) refreshToken() error {
	if c.auth.ExpirationDate.After(time.Now()) {
		return nil
	}

	c.logger.Debugw("refreshing linnworks auth token")

	url := fmt.Sprintf("%s/Auth/AuthorizeByApplication", LinnworksServer1)
	body := fmt.Sprintf("applicationId=%s&applicationSecret=%s&token=%s", c.Id, c.Secret, c.Token)
	payload := strings.NewReader(body)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	response, err := c.makeRequest(Post, url, payload, headers)
	if err != nil {
		return err
	}

	var authResp linnworksAuth
	if err := json.Unmarshal(response, &authResp); err != nil {
		return err
	}

	c.auth = authResp
	return nil
}
