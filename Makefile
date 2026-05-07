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
	-go run ./cmd/api/main.go

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
	psql "$(DB_DSN)" -f $(MIGRATIONS_DIR)/999_seed_sample_data.sql

setup: migrate-up seed

reset-db: migrate-down migrate-up seed

fresh: reset-db

# =========================
# DEMO COMMANDS
# =========================

.PHONY: frontend health branches loans fines cors uncompressed compressed compressed-view ratelimit metrics
frontend:
	python3 -m http.server 5500 --directory frontend

health:
	curl -s http://localhost:8080/health | jq

branches:
	curl -s http://localhost:8080/branches | jq

loans:
	curl -s http://localhost:8080/loans | jq

fines:
	curl -s http://localhost:8080/fines | jq

cors:
	curl -i -X OPTIONS http://localhost:8080/branches \
	-H "Origin: http://localhost:5500" \
	-H "Access-Control-Request-Method: GET"

uncompressed:
	curl -i http://localhost:8080/branches

compressed:
	curl -i --compressed -H "Accept-Encoding: gzip" http://localhost:8080/branches

compressed-view:
	curl -s --compressed -H "Accept-Encoding: gzip" http://localhost:8080/branches | jq

ratelimit:
	for i in {1..10}; do curl -i -s http://localhost:8080/branches | head -n 5; done

metrics:
	curl -s http://localhost:8080/metrics | jq

# =========================
# DATABASE COMMANDS
# =========================

db:
	psql "postgres://postgres:test1Access@localhost:5432/library-test1"

tables:
	psql "postgres://postgres:test1Access@localhost:5432/library-test1" -c "\dt"

branches-db:
	psql "postgres://postgres:test1Access@localhost:5432/library-test1" -c "SELECT * FROM branches;"

books-db:
	psql "postgres://postgres:test1Access@localhost:5432/library-test1" -c "SELECT * FROM books;"

members-db:
	psql "postgres://postgres:test1Access@localhost:5432/library-test1" -c "SELECT * FROM members;"

copies-db:
	psql "postgres://postgres:test1Access@localhost:5432/library-test1" -c "SELECT * FROM book_copies;"

loans-db:
	psql "postgres://postgres:test1Access@localhost:5432/library-test1" -c "SELECT l.id AS loan_id, m.first_name || ' ' || m.last_name AS member, bc.id AS copy_id, b.title AS book, br.name AS branch, l.borrowed_at, l.due_at, l.returned_at, l.status FROM loans l JOIN members m ON l.member_id = m.id JOIN book_copies bc ON l.book_copy_id = bc.id JOIN books b ON bc.book_id = b.id JOIN branches br ON bc.branch_id = br.id ORDER BY l.id;"

fines-db:
	psql "postgres://postgres:test1Access@localhost:5432/library-test1" -c "SELECT f.id AS fine_id, l.id AS loan_id, m.first_name || ' ' || m.last_name AS member, b.title AS book, br.name AS branch, f.amount, f.reason, f.paid FROM fines f JOIN loans l ON f.loan_id = l.id JOIN members m ON f.member_id = m.id JOIN book_copies bc ON l.book_copy_id = bc.id JOIN books b ON bc.book_id = b.id JOIN branches br ON bc.branch_id = br.id ORDER BY f.id;"

login:
	curl -s -X POST http://localhost:8080/auth/login \
	-H "Content-Type: application/json" \
	-d '{"email":"admin@library.local","password":"secret123"}' | jq

register-staff:
	curl -s -X POST http://localhost:8080/auth/register-staff \
	-H "Content-Type: application/json" \
	-d '{"email":"2009210050@ub.edu.bz","password":"secret123456","role":"admin","first_name":"New","last_name":"Staff","position":"Librarian","branch_id":1}' | jq

delete-member-auth:
	@TOKEN=$$(curl -s -X POST http://localhost:8080/auth/login \
	-H "Content-Type: application/json" \
	-d '{"email":"admin@library.local","password":"secret123"}' | jq -r '.token'); \
	curl -s -X DELETE http://localhost:8080/members/3 \
	-H "Authorization: Bearer $$TOKEN" | jq

delete-member-no-auth:
	curl -s -X DELETE http://localhost:8080/members/1 | jq
