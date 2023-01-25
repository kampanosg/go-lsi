package controllers

import (
	"net/http"

	"github.com/kampanosg/go-lsi/sync"
)

type SyncController struct {
	tool *sync.SyncTool
}

func NewSyncController(syncTool *sync.SyncTool) *SyncController {
	return &SyncController{tool: syncTool}
}

func (c *SyncController) HandleSyncRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("PONG"))
}
