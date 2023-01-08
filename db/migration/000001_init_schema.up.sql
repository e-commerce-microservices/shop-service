
CREATE TABLE shop (
    "id" serial8 PRIMARY KEY,
    "seller_id" serial8 NOT NULL,
    "name" varchar(64) NOT NULL,
    "avatar" varchar(256),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);