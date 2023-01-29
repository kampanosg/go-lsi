package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kampanosg/go-lsi/sync"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

var (
	errSyncCategories = errors.New("failed to sync categories")
	errSyncProducts   = errors.New("failed to sync products")
	errSyncOrders     = errors.New("failed to sync orders")
	errSyncStatus     = errors.New("failed to add sync status")
)

type SyncController struct {
	tool   *sync.SyncTool
	logger *zap.SugaredLogger
}

func NewSyncController(syncTool *sync.SyncTool, logger *zap.SugaredLogger) *SyncController {
	return &SyncController{tool: syncTool, logger: logger}
}

func (c *SyncController) HandleSyncRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.logger.Errorw("request failed", "reason", "method not supported", "method", r.Method, "uri", r.RequestURI)
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	var req types.SyncRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Errorw("request failed", "reason", "unable to decode body", "uri", r.RequestURI, "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	syncResp := types.SyncStatus{
		Timestamp: time.Now(),
	}

	if err := c.sync(req.From, req.To); err != nil {
		failed(w, err, http.StatusBadRequest)
		return
	}

	ok(w, syncResp)
}

func (c *SyncController) HandleSyncStatusRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		c.logger.Errorw("request failed", "reason", "method not supported", "method", r.Method, "uri", r.RequestURI)
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	status, err := c.tool.Db.GetLastSyncStatus()
	if err != nil {
		failed(w, err, http.StatusBadRequest)
		return
	}

	ok(w, status)
}

func (c *SyncController) HandleSyncRecentRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.logger.Errorw("request failed", "reason", "method not supported", "method", r.Method, "uri", r.RequestURI)
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	from := time.Now().Add(-time.Minute * 30)
	to := time.Now()

	syncResp := types.SyncStatus{
		Timestamp: time.Now(),
	}

	if err := c.sync(from, to); err != nil {
		c.logger.Errorw("syncing failed", "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	ok(w, syncResp)
}

func (c *SyncController) sync(from time.Time, to time.Time) error {
	c.logger.Infow("start syncing process", "from", from, "to", to)

	startTime := time.Now()

	if false {
		if err := c.tool.SyncCategories(); err != nil {
			return errSyncCategories
		}

		if err := c.tool.SyncProducts(); err != nil {
			return errSyncProducts
		}
	}

	if err := c.tool.SyncOrders(from, to); err != nil {
		return errSyncOrders
	}

	c.logger.Infow("finished syncing process", "from", from, "to", to, "elapsed", time.Since(startTime))

	if err := c.tool.Db.InsertSyncStatus(startTime.UnixMilli()); err != nil {
		return errSyncStatus
	}

	return nil
}
