.PHONY: help dev-up dev-down dev-logs prod-up prod-down build test migrate

# Default target
help:
	@echo "CloudManager - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev-up       - Start development infrastructure (docker-compose)"
	@echo "  make dev-down     - Stop development infrastructure"
	@echo "  make dev-logs     - View development logs"
	@echo "  make migrate      - Run database migrations"
	@echo ""
	@echo "Production:"
	@echo "  make prod-up      - Start production deployment"
	@echo "  make prod-down    - Stop production deployment"
	@echo "  make build        - Build all services"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run all tests"
	@echo "  make test-edge    - Run edge-agent tests"
	@echo "  make test-cloud   - Run cloud-backend tests"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean        - Clean build artifacts and volumes"
	@echo "  make proto        - Generate protobuf code"

# Development
dev-up:
	docker-compose up -d

dev-down:
	docker-compose down

dev-logs:
	docker-compose logs -f

dev-ps:
	docker-compose ps

# Database migrations
migrate:
	@echo "Running database migrations..."
	psql $(DATABASE_URL) -f migrations/001_initial.sql

# Production
prod-up:
	docker-compose -f docker-compose.prod.yml up -d

prod-down:
	docker-compose -f docker-compose.prod.yml down

prod-build:
	docker-compose -f docker-compose.prod.yml build

# Build
build:
	@echo "Building edge-agent..."
	cd edge-agent && go build -o bin/edge-agent ./cmd/agent
	@echo "Building cloud-backend services..."
	cd cloud-backend && go build -o bin/api-gateway ./api-gateway/cmd
	cd cloud-backend && go build -o bin/asset-service ./asset-service/cmd
	cd cloud-backend && go build -o bin/telemetry-service ./telemetry-service/cmd
	@echo "Building web-console..."
	cd web-console && npm run build

# Testing
test:
	@echo "Running edge-agent tests..."
	cd edge-agent && go test ./...
	@echo "Running cloud-backend tests..."
	cd cloud-backend && go test ./...
	@echo "Running web-console tests..."
	cd web-console && npm test

test-edge:
	cd edge-agent && go test ./... -v

test-cloud:
	cd cloud-backend && go test ./... -v

# Protocol Buffers
proto:
	@echo "Generating protobuf code..."
	cd edge-agent && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/edge.proto
	cd cloud-backend && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		shared/proto/edge.proto

# Cleaning
clean:
	@echo "Cleaning build artifacts..."
	rm -rf edge-agent/bin
	rm -rf cloud-backend/bin
	rm -rf web-console/dist
	@echo "Cleaning Docker volumes..."
	docker-compose down -v
	docker-compose -f docker-compose.prod.yml down -v

# Frontend development
web-dev:
	cd web-console && npm run dev

web-build:
	cd web-console && npm run build

web-lint:
	cd web-console && npm run lint

# Edge Agent
edge-build:
	cd edge-agent && go build -o bin/edge-agent ./cmd/agent

edge-run:
	cd edge-agent && go run ./cmd/agent -config config.yaml
