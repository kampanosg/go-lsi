package sqlite

const (
	query_GET_INVENTORY    = `SELECT p.lw_id, p.square_id, p.title, c.name, p.price, p.barcode, p.sku FROM product AS p JOIN category AS c ON c.lw_id = p.category_id`
	query_GET_CATEGORIES   = `SELECT lw_id, square_id, name, version FROM category;`
	query_INSERT_CATEGORY  = `INSERT INTO category (lw_id, square_id, name, version) VALUES (?, ?, ?, ?);`
	query_CLEAR_CATEGORIES = `DELETE FROM category;`

	query_GET_PRODUCTS          = `SELECT lw_id, square_id, square_var_id, category_id, square_category_id, title, price, barcode, sku, version FROM product;`
	query_GET_PRODUCT_BY_VAR_ID = `SELECT lw_id, square_id, square_var_id, category_id, square_category_id, title, price, barcode, sku, version FROM product WHERE square_var_id=?;`
	query_INSERT_PRODUCT        = `INSERT INTO product (lw_id, square_id, square_var_id, category_id, square_category_id, title, price, barcode, sku, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	query_CLEAR_PRODUCTS        = `DELETE FROM product;`

	query_GET_USER_BY_USERNAME = `SELECT id, username, password, friendly_name FROM user u WHERE u.username = ?;`

	query_GET_ORDERS         = `SELECT id, square_id, location_id, state, total_money, created_at FROM square_order;`
	query_GET_ORDER_PRODUCTS = `SELECT p.id, p.square_order_id, p.square_var_id, p.quantity FROM square_order_product AS p WHERE p.order_id = ?;`
	query_INSERT_ORDERS      = `
        INSERT INTO square_order (square_id, location_id, state, version, total_money, total_tax, total_discount, total_tip, total_service_charge, created_at) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	query_INSERT_ORDER_PRODUCTS = `
        INSERT INTO square_order_product (order_id, square_order_id, square_var_id, quantity) VALUES (?, ?, ?, ?);
    `

	query_INSERT_SYNC_STATUS   = `INSERT INTO sync_status (last_run) VALUES (?);`
	query_GET_LAST_SYNC_STATUS = `SELECT s.last_run FROM sync_status s ORDER BY s.last_run DESC LIMIT 1;`
)
