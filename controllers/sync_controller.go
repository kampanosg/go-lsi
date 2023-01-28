package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kampanosg/go-lsi/sync"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

type SyncController struct {
	tool   *sync.SyncTool
	logger *zap.SugaredLogger
}

func NewSyncController(syncTool *sync.SyncTool, logger *zap.SugaredLogger) *SyncController {
	return &SyncController{tool: syncTool}
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

	c.sync(req.From, req.To)
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

	c.sync(from, to)
}

func (c *SyncController) sync(from time.Time, to time.Time) {
	c.logger.Infow("start syncing process", "from", from, "to", to)

	startTime := time.Now()

	if err := c.tool.SyncCategories(); err != nil {
		c.logger.Errorw("syncing failed", "reason", "syncing categories failed", "error", err.Error())
		return
	}

	if err := c.tool.SyncProducts(); err != nil {
		c.logger.Errorw("syncing failed", "reason", "syncing products failed", "error", err.Error())
		return
	}

	if err := c.tool.SyncOrders(from, to); err != nil {
		c.logger.Errorw("syncing failed", "reason", "syncing orders failed", "error", err.Error())
		return
	}

	c.logger.Infow("finished syncing process", "from", from, "to", to, "elapsed", time.Since(startTime))

	if err := c.tool.Db.InsertSyncStatus(startTime.UnixMilli()); err != nil {
		c.logger.Errorw("unable to save sync status to db", "error", err.Error())
	}
}
