package transformers

import (
	"time"

	"github.com/kampanosg/go-lsi/types"
)

func FromSyncStatusDbRowToDomain(ts int64) types.SyncStatus {
	return types.SyncStatus{
		Timestamp: time.UnixMilli(ts),
	}
}
