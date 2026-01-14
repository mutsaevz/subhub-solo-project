# =========================
# Project configuration
# =========================
APP_NAME      := app
GO            := go
CMD_APP       := ./cmd/app/main.go
CMD_SEED      := ./cmd/seed/*.go
TMP_DIR       := tmp
BINARY        := $(TMP_DIR)/$(APP_NAME)

DOCKER_COMPOSE := docker compose

# =========================
# Phony targets
# =========================
.PHONY: help run dev build seed test cover fmt vet lint tidy clean \
        docker-up docker-down docker-rebuild docker-logs

# =========================
# Help
# =========================
help: ## Показать список команд
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# =========================
# Local Go commands
# =========================
run: ## Запуск HTTP-сервера локально
	$(GO) run $(CMD_APP)

dev: ## Запуск с hot reload (air)
	air -c .air.toml

build: ## Сборка бинарника
	mkdir -p $(TMP_DIR)
	$(GO) build -o $(BINARY) $(CMD_APP)

seed: ## Запуск сидов
	$(GO) run $(CMD_SEED)

# =========================
# Testing & quality
# =========================
test: ## Запуск всех тестов
	$(GO) test -v ./...

cover: ## Покрытие тестами
	$(GO) test -cover ./...

cover-html: ## HTML-отчёт покрытия
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

fmt: ## Форматирование кода
	$(GO) fmt ./...

vet: ## Статический анализ
	$(GO) vet ./...

lint: ## Линтинг (golangci-lint)
	golangci-lint run

tidy: ## Обновление go.mod / go.sum
	$(GO) mod tidy

# =========================
# Docker
# =========================
docker-up: ## Запуск проекта в Docker
	$(DOCKER_COMPOSE) up -d --build

docker-down: ## Остановка Docker-контейнеров
	$(DOCKER_COMPOSE) down

docker-rebuild: ## Полная пересборка Docker
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) up -d

docker-logs: ## Логи приложения
	$(DOCKER_COMPOSE) logs -f app

# =========================
# Cleanup
# =========================
clean: ## Очистка временных файлов
	rm -rf $(TMP_DIR) coverage.out