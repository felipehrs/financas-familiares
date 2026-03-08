.PHONY: up down dev test test-unit test-integration test-coverage build run lint

up:
	docker compose up -d

down:
	docker compose down

# Sobe o ambiente completo de desenvolvimento:
# 1. Garante que o PostgreSQL está rodando
# 2. Roda as migrations
# 3. Sobe o backend (Go) em background
# 4. Sobe o frontend (Vite dev server)
dev:
	@set -e; \
	ROOT=$(CURDIR); \
	echo "→ Verificando PostgreSQL..."; \
	if ! docker compose ps postgres 2>/dev/null | grep -q "running\|Up"; then \
		echo "→ Subindo PostgreSQL..."; \
		docker compose up -d postgres; \
		echo "→ Aguardando PostgreSQL ficar pronto..."; \
		until docker compose exec postgres pg_isready -U postgres -q 2>/dev/null; do sleep 1; done; \
	else \
		echo "→ PostgreSQL já está rodando."; \
	fi; \
	echo "→ Rodando migrations..."; \
	migrate -path $$ROOT/backend/migrations -database "postgres://postgres:postgres@localhost:5432/financas_familiares?sslmode=disable" up 2>&1 | grep -v "no change" || true; \
	echo "→ Subindo backend em background (logs em /tmp/backend.log)..."; \
	pkill -f "go run ./cmd/server" 2>/dev/null || true; \
	sleep 2; \
	(cd $$ROOT/backend && go run ./cmd/server > /tmp/backend.log 2>&1) & \
	echo "→ Aguardando backend iniciar..."; \
	for i in $$(seq 1 20); do \
		if curl -sf http://localhost:8080/health > /dev/null 2>&1; then \
			echo "→ Backend rodando em http://localhost:8080"; break; \
		fi; \
		if [ $$i -eq 20 ]; then \
			echo "✗ Backend falhou ao iniciar. Logs:"; cat /tmp/backend.log; exit 1; \
		fi; \
		sleep 1; \
	done; \
	echo "→ Subindo frontend..."; \
	cd $$ROOT/frontend && npm run dev

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
