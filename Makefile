.PHONY: help build up down logs shell db-shell clean dev prod

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development commands
dev: ## Start development environment with hot reloading
	docker-compose up --build

dev-detached: ## Start development environment in detached mode
	docker-compose up --build -d

stop: ## Stop all containers
	docker-compose down

restart: ## Restart all containers
	docker-compose restart

logs: ## Show logs from all containers
	docker-compose logs -f

logs-app: ## Show logs from app container only
	docker-compose logs -f app

logs-db: ## Show logs from database container only
	docker-compose logs -f postgres

# Database commands
db-shell: ## Connect to PostgreSQL database
	docker-compose exec postgres psql -U trendy_user -d trendy_repos

db-reset: ## Reset database (WARNING: This will delete all data)
	docker-compose down -v
	docker-compose up postgres -d
	sleep 10
	docker-compose up app -d

# Container management
shell: ## Get shell access to app container
	docker-compose exec app sh

build: ## Build Docker images
	docker-compose build

clean: ## Clean up Docker resources
	docker-compose down -v --remove-orphans
	docker system prune -f

# Production commands
prod: ## Start production environment
	docker-compose -f docker-compose.yml up --build -d --target production

prod-build: ## Build production image
	docker build --target production -t trendy-repos:latest .

# Development tools
pgadmin: ## Start with PgAdmin (database management tool)
	docker-compose --profile dev up -d

test: ## Run tests in container
	docker-compose exec app go test ./...

fmt: ## Format Go code
	docker-compose exec app go fmt ./...

vet: ## Run go vet
	docker-compose exec app go vet ./...

# Git helpers
git-add: ## Add all changes to git
	git add .

git-status: ## Show git status
	git status

# Quick development setup
setup: ## Initial setup - copy env file and start containers
	cp app.env.example app.env
	$(MAKE) dev 