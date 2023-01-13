package db

import (
	"database/sql"
	"log"

	"kev/types/domain"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteDb struct {
	Connection *sql.DB
}

func NewSqliteDB(dbPath string) sqliteDb {
	conn, err := sql.Open("sqlite3", dbPath)
	checkErr(err)
	log.Printf("Connected to the database: %s\n", dbPath)
	return sqliteDb{Connection: conn}
}

func (db *sqliteDb) GetCategories() (categories []domain.Category, err error) {
	q := `SELECT lw_id, name FROM category;`
	rows, err := db.Connection.Query(q)
	if err != nil {
		log.Printf("unable to get categories from db, reason=%v\n", err)
		return categories, err
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		var title string
		rows.Scan(&id, &title)
		categories = append(categories, domain.Category{Id: id, Name: title})
	}
	return categories, nil
}

func (db *sqliteDb) ClearCategories() error {
	q := `DELETE FROM category;`
	dl, err := db.Connection.Prepare(q)
	if err != nil {
		log.Printf("db: failed to clear categories. reason=%v", err.Error())
		return err
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		log.Printf("db: failed to clear categories. reason=%v", err.Error())
		return err
	}

	res, err := tx.Stmt(dl).Exec()

	if err != nil {
		tx.Rollback()
		log.Printf("db: failed to clear categories. reason=%v", err.Error())
		return err
	}

	tx.Commit()
	log.Printf("db: deleted categories. res=%v", res)
	return nil
}

func (db *sqliteDb) InsertCategories(categories []domain.Category) error {
	q := `INSERT INTO category (lw_id, name) VALUES (?, ?);`
	dl, err := db.Connection.Prepare(q)
	if err != nil {
		log.Printf("db: failed to create, reason=%v", err.Error())
		return err
	}

	tx, err := db.Connection.Begin()
	for _, category := range categories {
		if err != nil {
			log.Printf("db: failed to create category=%s, reason=%v", category.Id, err.Error())
			return err
		}

		res, err := tx.Stmt(dl).Exec(category.Id, category.Name)

		if err != nil {
			tx.Rollback()
			log.Printf("db: failed to create category=%s, reason=%v", category.Id, err.Error())
			return err
		}

		log.Printf("db: created category, id=%s, name=%s, res=%v", category.Id, category.Name, res)
	}
	tx.Commit()
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("DB Exception: %s\n", err.Error())
	}
}
