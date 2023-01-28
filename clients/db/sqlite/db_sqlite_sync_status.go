package sqlite

import (
	"github.com/kampanosg/go-lsi/transformers"
	"github.com/kampanosg/go-lsi/types"
)

func (db SqliteDb) InsertSyncStatus(ts int64) error {
	args := make([][]any, 1)
	args[0] = []any{ts}
	return db.commitTx(query_INSERT_SYNC_STATUS, args)
}

func (db SqliteDb) GetLastSyncStatus() (types.SyncStatus, error) {
	row := db.Connection.QueryRow(query_GET_LAST_SYNC_STATUS)
	var ts int64
	if err := row.Scan(&ts); err != nil {
		return types.SyncStatus{}, err
	}
	syncStatus := transformers.FromSyncStatusDbRowToDomain(ts)
	return syncStatus, nil
}
