# Makefile для lite-calendar
APP_NAME=lite-calendar

run:
	@echo "🚀 Running $(APP_NAME)..."
	go run ./cmd

fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

test:
	@echo "🧪 Running tests..."
	go test ./...

docker-up:
	@echo "🐳 Starting docker containers..."
	docker compose up -d

docker-down:
	@echo "🛑 Stopping containers..."
	docker compose down

docker-rebuild:
	docker compose down
	docker compose build app
	docker compose up -d

# Создать миграции
migrate-create:
	@read -p "Enter migration name: " name; \
	goose create -s -dir ./migrations $$name sql

# Применить все миграции
migrate-up:
	docker compose run --rm apply-migration

# 📊 Показать статус миграций
migrate-status:
	docker compose run --rm apply-migration -dir /app/migrations status



