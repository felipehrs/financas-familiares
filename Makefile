.PHONY: up down dev devcontainer-build test test-unit test-integration test-coverage build run lint

up:
	docker compose up -d

down:
	docker compose down

# Builda a imagem do devcontainer manualmente (opcional — VS Code faz isso automaticamente)
devcontainer-build:
	docker compose -f docker-compose.yml -f .devcontainer/docker-compose.yml build devcontainer

# Sobe o ambiente completo de desenvolvimento:
# 1. Garante que o PostgreSQL está rodando
# 2. Roda as migrations
# 3. Sobe o backend (Go) em background
# 4. Sobe o frontend (Vite dev server)
dev:
	@bash $(CURDIR)/dev.sh

# Backend targets
test:
	cd backend && go test ./...

test-unit:
	cd backend && go test -run Unit ./...

test-integration:
	cd backend && go test -run Integration ./...

test-coverage:
	cd backend && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

build:
	cd backend && go build -o bin/server ./cmd/server

run:
	cd backend && go run ./cmd/server

lint-backend:
	cd backend && golangci-lint run

# Frontend targets
test-frontend:
	cd frontend && pnpm test:run

test-frontend-coverage:
	cd frontend && pnpm test:coverage

lint-frontend:
	cd frontend && pnpm lint

# Combined
lint: lint-backend lint-frontend
