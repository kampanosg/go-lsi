package sync

import (
	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/clients/linnworks"
	"github.com/kampanosg/go-lsi/clients/square"
	"go.uber.org/zap"
)

const (
	reasonKey = "reason"
	errKey    = "error"
	msgDbErr  = "db client error"
	msgLwErr  = "linnworks client error"
	msgSqErr  = "square client error"
)

type SyncTool struct {
	LinnworksClient *linnworks.LinnworksClient
	SquareClient    *square.SquareClient
	Db              db.DB
	logger          *zap.SugaredLogger
}

func NewSyncTool(lwClient *linnworks.LinnworksClient, sqClient *square.SquareClient, db db.DB, logger *zap.SugaredLogger) *SyncTool {
	return &SyncTool{
		LinnworksClient: lwClient,
		SquareClient:    sqClient,
		Db:              db,
		logger:          logger,
	}
}
