package linnworks

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	Get  string = "GET"
	Post        = "POST"
)

func (c *LinnworksClient) makeRequest(method, url string, payload *strings.Reader, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

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

	return responseData, nil
}
