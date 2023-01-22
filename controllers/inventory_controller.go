package controllers

import (
	"net/http"

	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/transformers"
	"github.com/kampanosg/go-lsi/types"
)

type InventoryController struct {
	db db.DB
}

func NewInventoryController(db db.DB) InventoryController {
	return InventoryController{db: db}
}

func (c *InventoryController) HandleInventoryRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	items, err := c.db.GetInventory()
	if err != nil {
		failed(w, err, http.StatusBadRequest)
		return
	}

	resps := make([]types.InventoryItemResponse, len(items))
	for index, item := range items {
		resps[index] = transformers.FromInventoryDomainToResponse(item)
	}

	resp := types.InventoryResponse{
		Total: len(resps),
		Items: resps,
	}

	ok(w, resp)
}
