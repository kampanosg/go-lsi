package controllers

import (
	"net/http"

	"github.com/kampanosg/go-lsi/clients/db"
)

type InventoryController struct {
	db db.DB
}

func NewInventoryController(db db.DB) InventoryController {
	return InventoryController{db: db}
}

func (c *InventoryController) HandleInventoryRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	items, err := c.db.GetInventory()
	if err != nil {
		failed(w, err, http.StatusBadRequest)
		return
	}

    
}
