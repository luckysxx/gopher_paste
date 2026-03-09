SHELL := /bin/bash

BACKEND_DIR := backend
FRONTEND_DIR := frontend
SWAG_CMD := go run github.com/swaggo/swag/cmd/swag@latest

.PHONY: help swagger swagger-paste db-init db-migrate-paste dev dev-infra stop-infra lint lint-backend lint-frontend

help:
	@echo "Available targets:"
	@echo "  make swagger        	# Generate Swagger docs (paste service)"
	@echo "  make swagger-paste  	# Generate Swagger docs for paste service"
	@echo "  make db-init        	# Initialize users/pastes tables in postgres"
	@echo "  make db-migrate-paste 	# Run paste schema migration script"
	@echo "  make dev            	# Start infra + run paste/user/frontend locally"
	@echo "  make dev-infra      	# Start local postgres and redis via compose"
	@echo "  make stop-infra     	# Stop local postgres and redis"
	@echo "  make lint           	# Run backend test checks and frontend lint"
	@echo "  make lint-backend   	# Run go test ./..."
	@echo "  make lint-frontend  	# Run pnpm lint"

swagger: swagger-paste

swagger-paste:
	@cd $(BACKEND_DIR) && $(SWAG_CMD) init -g main.go -d services/paste -o services/paste/docs
	@echo "Swagger docs generated at backend/services/paste/docs"

db-init:
	@docker exec -i gopher_db psql -U luckys -d gopher_paste < $(BACKEND_DIR)/services/user/db/schema.sql
	@docker exec -i gopher_db psql -U luckys -d gopher_paste < $(BACKEND_DIR)/services/paste/db/schema.sql
	@echo "Database initialized: users, pastes"

db-migrate-paste:
	@docker exec -i gopher_db psql -U luckys -d gopher_paste < $(BACKEND_DIR)/services/paste/db/migrations/001_upgrade_pastes_to_snippets.sql
	@docker exec -i gopher_db psql -U luckys -d gopher_paste -c "SELECT column_name, data_type FROM information_schema.columns WHERE table_name='pastes' ORDER BY ordinal_position;"
	@echo "Paste migration applied and verified"

dev-infra:
	@docker network inspect gopher-net >/dev/null 2>&1 || docker network create gopher-net
	@docker compose -f docker-compose-infra.yaml up -d postgres redis
	@echo "Infra is running: postgres(5432), redis(6379)"

stop-infra:
	@docker compose -f docker-compose-infra.yaml stop postgres redis
	@echo "Infra stopped"

dev: dev-infra
	@set -euo pipefail; \
	trap 'kill 0' INT TERM EXIT; \
	(cd $(BACKEND_DIR)/services/paste && go run main.go) & \
	(cd $(BACKEND_DIR)/services/user && SERVER_PORT=8081 go run main.go) & \
	(cd $(FRONTEND_DIR) && pnpm dev) & \
	wait

lint: lint-backend lint-frontend

lint-backend:
	@cd $(BACKEND_DIR) && go test ./...

lint-frontend:
	@cd $(FRONTEND_DIR) && pnpm lint
