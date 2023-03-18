-- name: CreateShop :exec
INSERT INTO shop ("seller_id", "name", "avatar") VALUES ($1, $2, $3);

-- name: GetShopByID :one
SELECT * FROM shop WHERE "seller_id" = $1;

-- name: UpdateShopName :exec
UPDATE "shop"
SET "name" = $1
WHERE "seller_id" = $2;
