FROM golang:1.25-bookworm

# System deps + PostgreSQL client (pg_isready)
RUN apt-get update && apt-get install -y \
    curl \
    git \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Node.js 22
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# pnpm
RUN npm install -g pnpm

# go-migrate com driver postgres
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# golangci-lint
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

WORKDIR /workspace
