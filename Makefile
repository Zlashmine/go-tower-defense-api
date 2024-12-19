include .envrc

MIGRATIONS_PATH = ./cmd/migrations

.PHONY: go
run-app:
	@go run cmd/api/*~*_test.go*

.PHONY: db
db:
	@docker-compose up --build

.PHONY: migrate-create
migration:
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,lib && swag fmt
