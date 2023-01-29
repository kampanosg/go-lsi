package gormsqlite

import (
	"time"

	"github.com/kampanosg/go-lsi/models"
	"github.com/kampanosg/go-lsi/types"
)

func (db SqliteDb) GetLastSyncStatus() (types.SyncStatus, error) {
	var result models.SyncStatus
	db.Connection.Last(&result)
	if result.ID == 0 {
		return types.SyncStatus{}, errRecordNotFound
	}
	return fromSyncStatusModelToType(result), nil
}

func (db SqliteDb) InsertSyncStatus(ts int64) error {
	result := db.Connection.Create(models.SyncStatus{LastRun: time.UnixMilli(ts), Success: true})
	return result.Error
}

func fromSyncStatusModelToType(s models.SyncStatus) types.SyncStatus {
	return types.SyncStatus{Timestamp: s.LastRun}
}
