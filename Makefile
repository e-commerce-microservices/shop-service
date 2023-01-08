DB_DSN := postgres://admin:admin@192.168.49.2:30013/shop?sslmode=disable

migratecreate:
	migrate create -ext sql -dir db/migration -seq ${f}

migrateup:
	migrate -path db/migration -database "${DB_DSN}" -verbose up ${v}

migratedown:
	migrate -path db/migration -database "${DB_DSN}" -verbose down ${v}

migrateforce:
	migrate -path db/migration -database "${DB_DSN}" -verbose force ${v}

protogen:
	protoc --proto_path=proto proto/shop_service.proto proto/general.proto proto/auth_service.proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative

sqlcgen:
	sqlc generate

.PHONY: dockerbp
dockerbp:
	docker build -t ngoctd/ecommerce-shop:latest . && \
	docker push ngoctd/ecommerce-shop

.PHONY: redeploy
redeploy:
	kubectl rollout restart deployment depl-shop

.PHONY: migratecreate migrateup migratedown migrateforce protogen_auth protogen_product sqlcgen