CREATE TABLE "product" (
	"id"	INTEGER NOT NULL UNIQUE,
	"lw_id"	TEXT NOT NULL,
	"square_id"	TEXT NOT NULL,
	"square_var_id"	TEXT NOT NULL,
	"category_id"	INTEGER NOT NULL,
	"square_category_id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"barcode"	TEXT,
	"price"	REAL NOT NULL,
	"sku"	TEXT,
	"version"	INTEGER DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("category_id") REFERENCES "category"("lw_id")
);

CREATE INDEX "product_barcode_idx" ON "product" (
	"barcode"
);

CREATE INDEX "product_lw_id_idx" ON "product" (
	"lw_id"
);

CREATE INDEX "product_price_idx" ON "product" (
	"price"	ASC
);

CREATE INDEX "product_sku_idx" ON "product" (
	"sku"
);

CREATE INDEX "product_square_category_id_idx" ON "product" (
	"square_category_id"
);

CREATE INDEX "product_square_id_idx" ON "product" (
	"square_id"
);

CREATE INDEX "product_square_var_id_idx" ON "product" (
	"square_var_id"
);
