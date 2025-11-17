.PHONY: help install dev build start stop clean docker-build docker-up docker-down docker-logs

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install all dependencies (Go and Node.js)
	@echo "📦 Installing dependencies..."
	@cd go-core && go mod download
	@cd web && npm install
	@echo "✅ Dependencies installed!"

dev: ## Start development servers (backend + frontend)
	@echo "🚀 Starting development servers..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@make -j2 dev-backend dev-frontend

dev-backend: ## Start backend development server
	@cd go-core && go run cmd/main.go

dev-frontend: ## Start frontend development server
	@cd web && npm run dev

build: ## Build both backend and frontend
	@echo "🔨 Building backend..."
	@cd go-core && go build -o ../bin/100xtrader-api ./cmd/main.go
	@echo "🔨 Building frontend..."
	@cd web && npm run build
	@echo "✅ Build complete!"

start: ## Start production servers
	@echo "🚀 Starting production servers..."
	@./bin/100xtrader-api &
	@cd web && npm start

stop: ## Stop all running servers
	@pkill -f 100xtrader-api || true
	@pkill -f "next" || true
	@echo "✅ Servers stopped"

clean: ## Clean build artifacts and dependencies
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -rf web/.next
	@rm -rf web/node_modules
	@rm -rf go-core/vendor
	@rm -f *.sqlite
	@rm -f go-core/db.sqlite go-core/temp_db.sqlite
	@echo "✅ Clean complete!"

docker-build: ## Build Docker images
	@echo "🐳 Building Docker images..."
	docker-compose build
	@echo "✅ Docker images built!"

docker-up: ## Start Docker containers
	@echo "🐳 Starting Docker containers..."
	docker-compose up -d
	@echo "✅ Containers started!"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"

docker-down: ## Stop Docker containers
	@echo "🐳 Stopping Docker containers..."
	docker-compose down
	@echo "✅ Containers stopped!"

docker-logs: ## View Docker container logs
	docker-compose logs -f

docker-clean: ## Remove Docker containers and volumes
	@echo "🐳 Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f
	@echo "✅ Docker cleaned!"

setup: install ## Complete setup (install dependencies)
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  make dev          - Start development servers"
	@echo "  make docker-up    - Start with Docker"

