CREATE TABLE "product" (
	"id"	INTEGER NOT NULL UNIQUE,
	"lw_id"	TEXT NOT NULL,
	"category_id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"barcode"	TEXT,
	"price"	REAL NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("category_id") REFERENCES "category"("lw_id") ON DELETE CASCADE
);
