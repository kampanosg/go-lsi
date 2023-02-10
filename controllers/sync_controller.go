package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/sync"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

type SyncController struct {
	tool   sync.ST
	db     db.DB
	logger *zap.SugaredLogger
}

func NewSyncController(syncTool sync.ST, db db.DB, logger *zap.SugaredLogger) *SyncController {
	return &SyncController{tool: syncTool, db: db, logger: logger}
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

	c.handleSync(w, req.From, req.To)
}

func (c *SyncController) HandleSyncStatusRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		c.logger.Errorw("request failed", "reason", "method not supported", "method", r.Method, "uri", r.RequestURI)
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	status, err := c.db.GetLastSyncStatus()
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

	c.handleSync(w, from, to)
}

func (c *SyncController) handleSync(w http.ResponseWriter, from time.Time, to time.Time) {
	syncResp := types.SyncStatus{
		Timestamp: time.Now(),
	}

	if err := c.tool.Sync(from, to); err != nil {
		c.logger.Errorw("syncing failed", "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	ok(w, syncResp)
}
