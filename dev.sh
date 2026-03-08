#!/usr/bin/env bash
set -e

ROOT="$(cd "$(dirname "$0")" && pwd)"
BACKEND_PID_FILE="/tmp/financas-backend.pid"
BACKEND_LOG="/tmp/financas-backend.log"
DB_URL="postgres://postgres:postgres@localhost:5432/financas_familiares?sslmode=disable"

# ── 1. PostgreSQL ────────────────────────────────────────────────────────────
echo "→ Verificando PostgreSQL..."
if ! docker compose -f "$ROOT/docker-compose.yml" ps postgres 2>/dev/null | grep -q "running\|Up"; then
  echo "→ Subindo PostgreSQL..."
  docker compose -f "$ROOT/docker-compose.yml" up -d postgres
  echo "→ Aguardando PostgreSQL ficar pronto..."
  until docker compose -f "$ROOT/docker-compose.yml" exec postgres pg_isready -U postgres -q 2>/dev/null; do
    sleep 1
  done
else
  echo "→ PostgreSQL já está rodando."
fi

# ── 2. Migrations ────────────────────────────────────────────────────────────
echo "→ Rodando migrations..."
migrate -path "$ROOT/backend/migrations" -database "$DB_URL" up 2>&1 | grep -v "no change" || true

# ── 3. Backend ───────────────────────────────────────────────────────────────
echo "→ Parando backend anterior (se houver)..."
if [ -f "$BACKEND_PID_FILE" ]; then
  OLD_PID=$(cat "$BACKEND_PID_FILE")
  if kill -0 "$OLD_PID" 2>/dev/null; then
    kill "$OLD_PID" 2>/dev/null || true
    sleep 1
  fi
  rm -f "$BACKEND_PID_FILE"
fi

echo "→ Subindo backend (logs em $BACKEND_LOG)..."
(cd "$ROOT/backend" && go run ./cmd/server) > "$BACKEND_LOG" 2>&1 &
echo $! > "$BACKEND_PID_FILE"

echo "→ Aguardando backend iniciar..."
for i in $(seq 1 30); do
  if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
    echo "→ Backend rodando em http://localhost:8080 (PID $(cat $BACKEND_PID_FILE))"
    break
  fi
  if [ "$i" -eq 30 ]; then
    echo "✗ Backend falhou ao iniciar. Logs:"
    cat "$BACKEND_LOG"
    exit 1
  fi
  sleep 1
done

# ── 4. Frontend ──────────────────────────────────────────────────────────────
echo "→ Subindo frontend em http://localhost:5173"
echo "   (Ctrl+C para encerrar — o backend continuará rodando)"
echo "   Para parar o backend: kill \$(cat $BACKEND_PID_FILE)"
echo ""
cd "$ROOT/frontend" && npm run dev
