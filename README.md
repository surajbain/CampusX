# CampusX API

Multi-college event, competition & verified entry platform.

## Stack
- Go 1.25 + Gin
- PostgreSQL 16 + GORM
- Redis 7
- JWT (access 15m + refresh 7d, rotation-ready)
- Argon2id (Phase 3)
- Razorpay + mock provider (Phase 8)
- zerolog structured logging

## Phase 1 — Foundation ✅
- Config with fail-fast validation
- Structured logging (zerolog)
- Postgres (GORM) with retry + pooling
- Redis with retry
- Middleware: RequestID, Logger, Recovery, CORS
- Health endpoints: `/health`, `/ready`
- Docker + docker-compose
- First passing tests

## Quick start

### 1. Local (with Docker for deps only)

```bash
cp .env.example .env
# edit .env if needed (defaults work for local docker)

# start postgres + redis
docker compose up -d postgres redis

# run API
make run

## Phase 2 — Database Schema + Migrations

- Embed-based migrations with `schema_migrations` tracking
- Idempotent, transactional, checksum-verified
- 16 core tables with UUID PKs, indexes, and FKs
- Anti-fraud UNIQUE constraints at DB level
- Seed data (IITB + 6 users + 2 events)
- GORM models for all tables

### Run migrations
```bash
# Auto (on API boot)
RUN_MIGRATIONS=true

# Manual
make migrate

# Verify
make psql
\dt              -- list tables
\d payments      -- inspect constraints
SELECT * FROM schema_migrations;