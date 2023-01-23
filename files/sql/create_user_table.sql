CREATE TABLE "user" (
    "id"        INTEGER NOT NULL UNIQUE,
    "username"  TEXT NOT NULL,
    "password"  TEXT NOT NULL,
    "friendly_name" TEXT,
    PRIMARY KEY("id" AUTOINCREMENT)
)
