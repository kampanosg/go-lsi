package sync

import (

	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/clients/linnworks"
	"github.com/kampanosg/go-lsi/clients/square"
)

type SyncTool struct {
	LinnworksClient *linnworks.LinnworksClient
	SquareClient    *square.SquareClient
	Db              db.DB
}

func NewSyncTool(lwClient *linnworks.LinnworksClient, sqClient *square.SquareClient, db db.DB) *SyncTool {
	return &SyncTool{
		LinnworksClient: lwClient,
		SquareClient:    sqClient,
		Db:              db,
	}
}

