include .env
export

GOPATH ?= $(shell go env GOPATH)
MIGRATE = $(GOPATH)/bin/migrate

DB_URL := ${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?multiStatements=true

.PHONY: help migrate-create migrate-up migrate-down migrate-force gorm db_sql sqlx

migrate-create:
ifndef NAME
	$(error NAME is not set. Usage: make migrate-create NAME=<migration_name>)
endif
	@echo "Creating migration files for: $(NAME)"
	migrate create -ext sql -dir migration -seq $(NAME)

migrate-up:
	@echo "Running up migrations..."
	$(MIGRATE) -database "$(DB_URL)" -path migration up

migrate-down:
	@echo "Running down migrations..."
	$(MIGRATE) -database "$(DB_URL)" -path migration down

migrate-force:
	@echo "Forcing migration to version $(VERSION)..."
	$(MIGRATE) -database "$(DB_URL)" -path migration force $(VERSION)


gorm:
	go run gorm/cmd/main.go

db_sql:
	go run db_sql/cmd/main.go

sqlx:
	go run cmd/sqlx/main.go