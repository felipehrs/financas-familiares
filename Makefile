.PHONY: up down test test-unit test-integration test-coverage build run lint

up:
	docker compose up -d

down:
	docker compose down

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
