package controllers

import (
	"net/http"

	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

type InventoryController struct {
	db     db.DB
	logger *zap.SugaredLogger
}

func NewInventoryController(db db.DB, logger *zap.SugaredLogger) InventoryController {
	return InventoryController{db: db, logger: logger}
}

func (c *InventoryController) HandleInventoryRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		c.logger.Errorw("request failed", "reason", "method not supported", "method", r.Method, "uri", r.RequestURI)
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	items, err := c.db.GetInventory()
	if err != nil {
		c.logger.Errorw("request failed", "error retrieving inventory from db", "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	resp := types.HttpResponse{
		Total: len(items),
		Items: items,
	}

	ok(w, resp)
}
