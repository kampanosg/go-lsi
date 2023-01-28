package square

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	GET  string = "GET"
	POST        = "POST"
)

var (
	errBadRequest = errors.New("bad request")
)

func (c *SquareClient) makeRequest(method, url string, headers map[string]string, jsonReq []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonReq))

	if err != nil {
		return []byte{}, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		c.logger.Debugw("unable to make request", "error", err.Error())
		return []byte{}, err
	}
	defer res.Body.Close()

	responseData, err := ioutil.ReadAll(res.Body)
	c.logger.Debugw("http client responded", "body", string(responseData), "status", res.StatusCode)
	if err != nil {
		return responseData, err
	}

	if res.StatusCode >= 400 && res.StatusCode <= 500 {
		return responseData, errBadRequest
	}

	return responseData, nil
}
