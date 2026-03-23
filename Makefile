include .env
export

APP_NAME=test1
MAIN_PATH=./cmd/api/main.go
MIGRATIONS_DIR=./migrations

.PHONY: help run build tidy setup migrate-up migrate-down seed reset-db fresh

help:
	@echo ""
	@echo "Commands:"
	@echo " make setup      -> run migrations + seed data"
	@echo " make run        -> run the API"
	@echo " make build      -> build binary"
	@echo " make tidy       -> go mod tidy"
	@echo " make migrate-up -> run migrations"
	@echo " make migrate-down -> rollback migrations"
	@echo " make seed       -> insert sample data"
	@echo " make reset-db   -> rebuild database"
	@echo ""

run:
	go run $(MAIN_PATH)

build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

tidy:
	go mod tidy

migrate-up:
	@for file in $$(ls $(MIGRATIONS_DIR)/*.up.sql | sort); do \
		echo "Applying $$file"; \
		psql "$(DB_DSN)" -f $$file || exit 1; \
	done

migrate-down:
	@for file in $$(ls $(MIGRATIONS_DIR)/*.down.sql | sort -r); do \
		echo "Rolling back $$file"; \
		psql "$(DB_DSN)" -f $$file || exit 1; \
	done

seed:
	psql "$(DB_DSN)" -f $(MIGRATIONS_DIR)/005_seed_sample_data.sql

setup: migrate-up seed

reset-db: migrate-down migrate-up seed

fresh: reset-db