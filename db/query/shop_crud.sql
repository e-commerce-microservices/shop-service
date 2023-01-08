-- name: CreateShop :exec
INSERT INTO shop ("seller_id", "name", "avatar") VALUES ($1, $2, $3);