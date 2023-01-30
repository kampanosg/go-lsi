package controllers

import (
	"errors"
	"net/http"

	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

type InventoryController struct {
	db     db.DB
	logger *zap.SugaredLogger
}

var (
	errParamsNotProvided = errors.New("required params not provided")
)

func NewInventoryController(db db.DB, logger *zap.SugaredLogger) InventoryController {
	return InventoryController{db: db, logger: logger}
}

func (c *InventoryController) HandleInventoryRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		c.logger.Errorw("request failed", "reason", "method not supported", "method", r.Method, "uri", r.RequestURI)
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	foundSku := r.URL.Query().Has("sku")
	foundBarcode := r.URL.Query().Has("barcode")

	if !foundSku && !foundBarcode {
		failed(w, errParamsNotProvided, http.StatusBadRequest)
		return
	}

	var product types.Product
	var err error

	if foundBarcode {
		product, err = c.db.GetProductByBarcode(r.URL.Query().Get("barcode"))
	} else {
		product, err = c.db.GetProductBySku(r.URL.Query().Get("sku"))
	}

	if err != nil {
		failed(w, err, http.StatusNotFound)
		return
	}

	resp := types.OkResp{Items: []types.Product{product}, Total: 1}

	ok(w, resp)
}
