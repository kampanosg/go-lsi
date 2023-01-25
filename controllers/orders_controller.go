package controllers

import (
	"net/http"

	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/types"
)

type OrdersController struct {
	db db.DB
}

func NewOrdersController(db db.DB) OrdersController {
	return OrdersController{db: db}
}

func (c *OrdersController) HandleOrdersRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	orders, err := c.db.GetOrders()
	if err != nil {
		failed(w, err, http.StatusBadRequest)
		return
	}

	resp := types.HttpResponse{
		Total: len(orders),
		Items: orders,
	}

	ok(w, resp)
}
