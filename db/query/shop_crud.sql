-- name: CreateShop :exec
INSERT INTO shop ("seller_id", "name", "avatar") VALUES ($1, $2, $3);

-- name: GetShopByID :one
SELECT * FROM shop WHERE "id" = $1;