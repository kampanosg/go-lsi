package sqlite

const (
	query_GET_CATEGORIES   = `SELECT lw_id, square_id, name, version FROM category;`
	query_INSERT_CATEGORY  = `INSERT INTO category (lw_id, square_id, name, version) VALUES (?, ?, ?, ?);`
	query_CLEAR_CATEGORIES = `DELETE FROM category;`

	query_GET_PRODUCTS   = `SELECT lw_id, square_id, square_var_id, category_id, square_category_id, title, price, barcode, sku, version FROM product;`
	query_INSERT_PRODUCT = `INSERT INTO product (lw_id, square_id, square_var_id, category_id, square_category_id, title, price, barcode, sku, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	query_CLEAR_PRODUCTS = `DELETE FROM product;`

	query_GET_USER_BY_USERNAME = `SELECT id, username, password, friendly_name FROM user u WHERE u.username = ?`
	query_INSERT_SESSION       = `INSERT`
)
