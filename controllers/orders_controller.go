package controllers

import (
	"net/http"

	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

type OrdersController struct {
	db db.DB
logger *zap.SugaredLogger
}

func NewOrdersController(db db.DB, logger *zap.SugaredLogger) OrdersController {
    return OrdersController{db: db, logger: logger}
}

func (c *OrdersController) HandleOrdersRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
        c.logger.Errorw("request failed", "reason", "method not supported", "method", r.Method, "uri", r.RequestURI)
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	orders, err := c.db.GetOrders()
	if err != nil {
c.logger.Errorw("request failed", "error retrieving orders from db", "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	resp := types.HttpResponse{
		Total: len(orders),
		Items: orders,
	}

	ok(w, resp)
}
