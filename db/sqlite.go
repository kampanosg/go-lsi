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
	return sqliteDb{Connection: conn}
}

func (db *sqliteDb) GetCategories() (categories []domain.Category, err error) {
	q := `SELECT lw_id, square_id, name, version FROM category;`
	rows, err := db.Connection.Query(q)
	if err != nil {
		log.Printf("unable to get categories from db, reason=%v\n", err)
		return categories, err
	}

	defer rows.Close()

	for rows.Next() {
		var id, squareId, name string
		var version int64
		rows.Scan(&id, &squareId, &name, &version)
		categories = append(categories, domain.Category{
			Id:       id,
			SquareId: squareId,
			Name:     name,
			Version:  version,
		})
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
	q := `INSERT INTO category (lw_id, square_id, name, version) VALUES (?, ?, ?, ?);`
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

		res, err := tx.Stmt(dl).Exec(category.Id, category.SquareId, category.Name, category.Version)

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

func (db *sqliteDb) GetProducts() (products []domain.Product, err error) {
	q := `SELECT lw_id, square_id, square_var_id, category_id, square_category_id, title, price, barcode, sku, version FROM product;`
	rows, err := db.Connection.Query(q)
	if err != nil {
		log.Printf("unable to get products from db, reason=%v\n", err)
		return products, err
	}

	defer rows.Close()

	for rows.Next() {
		var id, squareId, squareVarId, categoryId, squareCategoryId, title, barcode, sku string
		var price float64
		var version int64
		rows.Scan(&id, &squareId, &squareVarId, &categoryId, &squareCategoryId, &title, &price, &barcode, &sku, &version)
		products = append(products, domain.Product{
			Id:          id,
			SquareId:    squareId,
			SquareVarId: squareVarId,
			CategoryId:  categoryId,
            SquareCategoryId: squareCategoryId,
			Title:       title,
			Price:       price,
			Barcode:     barcode,
			SKU:         sku,
			Version:     version,
		})

	}
	return products, nil
}

func (db *sqliteDb) ClearProducts() error {
	q := `DELETE FROM product;`
	dl, err := db.Connection.Prepare(q)
	if err != nil {
		log.Printf("db: failed to clear products. reason=%v", err.Error())
		return err
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		log.Printf("db: failed to clear products. reason=%v", err.Error())
		return err
	}

	res, err := tx.Stmt(dl).Exec()

	if err != nil {
		tx.Rollback()
		log.Printf("db: failed to clear products. reason=%v", err.Error())
		return err
	}

	tx.Commit()
	log.Printf("db: deleted products. res=%v", res)
	return nil
}

func (db *sqliteDb) InsertProducts(products []domain.Product) error {
	q := `INSERT INTO product
            (lw_id, square_id, square_var_id, category_id, square_category_id, title, price, barcode, sku, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	dl, err := db.Connection.Prepare(q)
	if err != nil {
		log.Printf("db: failed to create, reason=%v", err.Error())
		return err
	}

	tx, err := db.Connection.Begin()
	for _, product := range products {
		if err != nil {
			log.Printf("db: failed to create category=%s, reason=%v", product.Id, err.Error())
			return err
		}

		res, err := tx.Stmt(dl).Exec(product.Id, product.SquareId, product.SquareVarId, product.CategoryId, 
            product.SquareCategoryId, product.Title, product.Price, product.Barcode, product.SKU, product.Version)

		if err != nil {
			tx.Rollback()
			log.Printf("db: failed to create product=%s, reason=%v", product.Id, err.Error())
			return err
		}

		log.Printf("db: created product, id=%s, name=%s, res=%v", product.Id, product.Title, res)
	}
	tx.Commit()
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("DB Exception: %s\n", err.Error())
	}
}
