package sync

import (
	"time"

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

type ST interface {
	Sync(from time.Time, to time.Time) error
	SyncCategories() error
	SyncProducts() error
	SyncOrders(start time.Time, end time.Time) error
}

type SyncTool struct {
	LinnworksClient linnworks.LW
	SquareClient    square.SQ
	Db              db.DB
	logger          *zap.SugaredLogger
}

func NewSyncTool(lwClient linnworks.LW, sqClient square.SQ, db db.DB, logger *zap.SugaredLogger) *SyncTool {
	return &SyncTool{
		LinnworksClient: lwClient,
		SquareClient:    sqClient,
		Db:              db,
		logger:          logger,
	}
}

func (s *SyncTool) Sync(from time.Time, to time.Time) error {
	s.logger.Infow("start syncing process", "from", from, "to", to)

	startTime := time.Now()
	if err := s.SyncCategories(); err != nil {
		return err
	}

	if err := s.SyncProducts(); err != nil {
		return err
	}

	if err := s.SyncOrders(from, to); err != nil {
		return err
	}

	s.logger.Infow("finished syncing process", "from", from, "to", to, "elapsed", time.Since(startTime))

	if err := s.Db.InsertSyncStatus(startTime.UnixMilli()); err != nil {
		return err
	}

	return nil
}
