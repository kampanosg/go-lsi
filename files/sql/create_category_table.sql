CREATE TABLE "category" (
	"id"	INTEGER NOT NULL UNIQUE,
	"lw_id"	TEXT NOT NULL,
	"square_id"	INTEGER NOT NULL,
	"name"	TEXT NOT NULL,
	"version"	INTEGER DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE INDEX "category_lw_id_idx" ON "category" (
	"lw_id"
);

CREATE INDEX "category_square_id_idx" ON "category" (
	"square_id"
);
