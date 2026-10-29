.PHONY: help build run test docker-build docker-up docker-down docker-logs docker-tools docker-restart docker-clean

help: ## Muestra la ayuda de los comandos disponibles
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ==========================================
# Desarrollo Local en Go
# ==========================================
build: ## Compila el binario localmente
	go build -v -o bin/clienthub ./cmd/main.go

run: ## Ejecuta la aplicación localmente
	go run ./cmd/main.go

test: ## Ejecuta las pruebas unitarias
	go test -v ./...

# ==========================================
# Docker & Docker Compose
# ==========================================
docker-build: ## Construye la imagen Docker de la API
	docker compose build

docker-up: ## Inicia los contenedores esenciales (API y PostgreSQL) en segundo plano
	docker compose up -d --build

docker-down: ## Detiene y remueve los contenedores activos
	docker compose down

docker-logs: ## Muestra los logs de todos los contenedores en tiempo real
	docker compose logs -f

docker-logs-api: ## Muestra los logs de la API Go
	docker compose logs -f api

docker-logs-db: ## Muestra los logs de PostgreSQL
	docker compose logs -f postgres

docker-tools: ## Inicia todos los servicios incluyendo pgAdmin (puerto 5050)
	docker compose --profile tools up -d

docker-restart: ## Reinicia los contenedores de la aplicación
	docker compose restart

docker-clean: ## Detiene los contenedores y elimina los volúmenes de datos
	docker compose down -v

