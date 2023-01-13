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

func (c *LinnworksClient) GetCategories() {
	c.refreshToken()
	fmt.Println(c.auth.Token)
}

func (c *LinnworksClient) refreshToken() {

    url := "https://api.linnworks.net/api/Auth/AuthorizeByApplication"
    method := "POST"

    body := fmt.Sprintf("applicationId=%s&applicationSecret=%s&token=%s", c.id, c.secret, c.token)
    payload := strings.NewReader(body)

    client := &http.Client {}
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
