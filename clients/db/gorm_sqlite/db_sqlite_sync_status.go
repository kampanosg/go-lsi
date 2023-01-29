package gormsqlite

func (db SqliteSqliteDb) GetLastSyncStatus() (types.SyncStatus, error) {}
func (db SqliteDb) InsertSyncStatus(ts int64) error                    {}
