package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDb struct {
	Connection *sql.DB
}

func NewSqliteDB(dbPath string) SqliteDb {
	conn, err := sql.Open("sqlite3", dbPath)
	checkErr(err)
	return SqliteDb{Connection: conn}
}

func (db SqliteDb) commitTx(query string, args [][]any) error {
	stmt, err := db.Connection.Prepare(query)
	if err != nil {
		log.Printf("db: failed to commit tx. reason=%v\n", err.Error())
		return err
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		log.Printf("db: failed to commit tx. reason=%v\n", err.Error())
		return err
	}

	err = nil

	if len(args) > 0 {
		for _, arg := range args {
			if _, err := tx.Stmt(stmt).Exec(arg...); err != nil {
				return err
			}
		}
	} else {
		_, err = tx.Stmt(stmt).Exec()
	}

	if err != nil {
		tx.Rollback()
		log.Printf("db: failed to exec tx. reason=%v\n", err.Error())
		return err
	}

	return tx.Commit()
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("db: exception, err=%s\n", err.Error())
	}
}
