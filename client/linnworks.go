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
	method := "GET"

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

	return transform.FromArrCategoryRespToDomain(authResp), nil
}

func (c *LinnworksClient) refreshToken() {

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
