package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/kampanosg/go-lsi/clients/square"
	"go.uber.org/zap"
)

var (
	excludedIds = map[string]bool{
		"VFCFO4YK5QCIGIPI3QGJVK4I": true,
		"AHGU7NEM4PODEE3UMRSKAAG2": true,
		"25O5FVWY76CRZRZAFIWF26MS": true,
		"6RDRSF6PLPI6GYAUOK7R6ASZ": true,
		"TLR4MWIMYVLFBWXLD7K3ZJD5": true,
		"RNHZ4FHU67S7AERXP4PB6WCW": true,
		"7WYBRBZIYEXUZL43UGHHDB7H": true,
		"D45BUSWUUII2MZO7XDHRKBMB": true,
		"DHDHZFYTAFHDBL7MG26TDAAB": true,
		"CDL27LZ5AVFOARQQHARB4W5F": true,
	}
)

// / usage: go run delete_square_inventory.go ACCESS_TOKEN LOCATION_ID
func main() {
	args := os.Args[1:]

	accessToken := args[0]
	host := "https://connect.squareup.com/v2"
	version := "2023-01-19"
	location := args[1]
	logger := zap.NewExample().Sugar()

	headers := make(map[string]string)
	headers["Square-Version"] = "2023-01-19"
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("Bearer %s", args[0])

	ids := make([]string, 0)
	titles := make(map[string]string, 0)
	cursor := ""

	for {
		url := fmt.Sprintf("%s/catalog/list?types=ITEM&cursor=%s", host, cursor)
		resp, err := makeRequest2("GET", url, headers, []byte{})
		if err != nil {
			panic(err)
		}

		var r squareResp2
		if err := json.Unmarshal(resp, &r); err != nil {
			panic(err)
		}

		for _, o := range r.Objects {
			if isExcludedID(o.ID) {
				continue
			}
			ids = append(ids, o.ID)
			titles[o.ID] = o.ItemData.Name
		}

		cursor = r.Cursor
		if cursor == "" {
			break
		}
	}

	// fmt.Printf("total items to delete: %d\n", len(ids))
	// for id, title := range titles {
	// fmt.Printf("%s - %s\n", id, title)
	// }

	// panic("stahp")

	client := square.NewSquareClient(accessToken, host, version, location, make([]string, 0), logger)
	client.BatchDeleteItems(ids)
}

func isExcludedID(id string) bool {
	return excludedIds[id]
}

func makeRequest2(method, url string, headers map[string]string, jsonReq []byte) ([]byte, error) {
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
		return []byte{}, err
	}
	defer res.Body.Close()

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return responseData, err
	}

	if res.StatusCode >= 400 && res.StatusCode <= 500 {
		return responseData, errors.New("bad request")
	}

	return responseData, nil
}

type squareResp2 struct {
	Cursor  string `json:"cursor"`
	Objects []struct {
		Type                  string    `json:"type"`
		ID                    string    `json:"id"`
		UpdatedAt             time.Time `json:"updated_at"`
		CreatedAt             time.Time `json:"created_at"`
		Version               int64     `json:"version"`
		IsDeleted             bool      `json:"is_deleted"`
		PresentAtAllLocations bool      `json:"present_at_all_locations"`
		ItemData              struct {
			Name       string `json:"name"`
			IsTaxable  bool   `json:"is_taxable"`
			Visibility string `json:"visibility"`
			CategoryID string `json:"category_id"`
			Variations []struct {
				Type                  string    `json:"type"`
				ID                    string    `json:"id"`
				UpdatedAt             time.Time `json:"updated_at"`
				CreatedAt             time.Time `json:"created_at"`
				Version               int64     `json:"version"`
				IsDeleted             bool      `json:"is_deleted"`
				PresentAtAllLocations bool      `json:"present_at_all_locations"`
				ItemVariationData     struct {
					ItemID      string `json:"item_id"`
					Name        string `json:"name"`
					Sku         string `json:"sku"`
					Upc         string `json:"upc"`
					Ordinal     int    `json:"ordinal"`
					PricingType string `json:"pricing_type"`
					PriceMoney  struct {
						Amount   int    `json:"amount"`
						Currency string `json:"currency"`
					} `json:"price_money"`
					Sellable  bool `json:"sellable"`
					Stockable bool `json:"stockable"`
				} `json:"item_variation_data"`
			} `json:"variations"`
			ProductType        string `json:"product_type"`
			SkipModifierScreen bool   `json:"skip_modifier_screen"`
		} `json:"item_data"`
	} `json:"objects"`
}
