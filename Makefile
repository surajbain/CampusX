.PHONY: help dev build run test test-verbose lint tidy fmt docker-up docker-down docker-logs clean migrate migrate-fresh psql db-reset test-integration

BINARY_NAME=campusx-api

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

tidy: ## Run go mod tidy
	cd backend && go mod tidy

build: ## Build binary
	cd backend && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/$(BINARY_NAME) ./cmd/api

run: ## Run API locally (needs Postgres + Redis)
	cd backend && go run ./cmd/api

dev: ## Run API with live logs
	cd backend && go run ./cmd/api

test: ## Run all tests
	cd backend && go test ./... -count=1 -race

test-verbose: ## Run tests verbose
	cd backend && go test ./... -count=1 -race -v

lint: ## Run go vet + staticcheck (if installed)
	cd backend && go vet ./...

fmt: ## gofmt all files
	cd backend && gofmt -s -w .

docker-up: ## Start postgres + redis + api
	docker compose up --build -d

docker-down: ## Stop all containers
	docker compose down

docker-logs: ## Tail api logs
	docker compose logs -f api

clean: ## Remove build artifacts
	rm -rf backend/bin/ backend/tmp/ backend/coverage.out

migrate: ## Run migrations against dev DB
	cd backend && RUN_MIGRATIONS=true go run ./cmd/api

migrate-fresh: ## Drop all tables + re-run migrations (DEV ONLY)
	docker compose exec -T postgres psql -U campusx -d campusx -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO campusx;"
	@echo "Schema reset. Restart API with RUN_MIGRATIONS=true."

psql: ## Open psql shell in running Postgres
	docker compose exec postgres psql -U campusx -d campusx

db-reset: ## Full reset (containers + volumes)
	docker compose down -v
	docker compose up --build -d

test-integration: ## Run tests against running Postgres
	cd backend && TEST_DSN="host=localhost port=5432 user=campusx password=campusx_dev_password dbname=campusx sslmode=disable" go test ./tests/... -v -count=1