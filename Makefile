DATABASE_URL ?= postgres://airfa:airfa@localhost:5432/airfa_dev?sslmode=disable
export DATABASE_URL

.PHONY: dev migrate migrate-down migrate-create sqlc seed test

dev:
	@trap 'kill 0' EXIT; \
	(cd apps/api && go run ./cmd/server) & \
	(npm run dev --workspace apps/web) & \
	wait

migrate:
	migrate -path apps/api/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path apps/api/migrations -database "$(DATABASE_URL)" down 1

# usage: make migrate-create name=add_something
migrate-create:
	migrate create -ext sql -dir apps/api/migrations -seq $(name)

sqlc:
	cd apps/api && sqlc generate

seed:
	cd apps/api && go run ./cmd/seed

test:
	cd apps/api && go test ./...
	npm test --workspaces --if-present