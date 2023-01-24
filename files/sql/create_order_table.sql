CREATE TABLE "square_order" (
    "id"    INTEGER NOT NULL UNIQUE,
    "square_id" TEXT NOT NULL,
    "location_id"   TEXT NOT NULL,
    "state" TEXT NOT NULL,
    "version"   INTEGER NOT NULL,
    "total_money"   REAL NOT NULL,
    "total_tax" REAL,
    "total_discount"    REAL,
    "total_tip" REAL,
    "total_service_charge"  REAL,
    "created_at"    REAL NOT NULL,

    PRIMARY KEY("id" AUTOINCREMENT)
)
