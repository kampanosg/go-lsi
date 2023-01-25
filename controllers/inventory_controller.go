package controllers

import (
	"net/http"

	"github.com/kampanosg/go-lsi/clients/db"
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

	resp := types.HttpResponse{
		Total: len(items),
		Items: items,
	}

	ok(w, resp)
}
