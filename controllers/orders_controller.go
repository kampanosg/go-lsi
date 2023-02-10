package controllers

import (
	"net/http"
	"time"

	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

type OrdersController struct {
	db     db.DB
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

	foundStartDate := r.URL.Query().Has("start")
	foundEndDate := r.URL.Query().Has("end")

	if !foundStartDate || !foundEndDate {
		failed(w, errParamsNotProvided, http.StatusBadRequest)
		return
	}

	start, err := parseDateTime(r.URL.Query().Get("start"))
	if err != nil {
		failed(w, err, http.StatusBadRequest)
		return
	}

	end, err := parseDateTime(r.URL.Query().Get("end"))
	if err != nil {
		failed(w, err, http.StatusBadRequest)
		return
	}

	orders, err := c.db.GetOrdersWithinRange(start, end)
	if err != nil {
		c.logger.Errorw("request failed", "error retrieving orders from db", "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	resp := types.OkResp{
		Total: len(orders),
		Items: orders,
	}

	ok(w, resp)
}

func parseDateTime(s string) (time.Time, error) {
	dateFormat := "2006-01-02T15:04:05Z0700"
	date, err := time.Parse(dateFormat, s)
	if err != nil {
		return time.Now(), err
	}
	return date, nil
}