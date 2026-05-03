.PHONY: help deps proto build run docker-up docker-down test clean install-grpcurl test-register test-login list-methods stop all

# Помощь
help:
	@echo "=== Auth Service Makefile ==="
	@echo "Доступные команды:"
	@echo "  make all           - Полная установка и запуск (deps + proto + docker-up + run)"
	@echo "  make deps          - Установка всех зависимостей"
	@echo "  make proto         - Генерация gRPC кода из proto файлов"
	@echo "  make build         - Сборка бинарника"
	@echo "  make run           - Запуск сервера"
	@echo "  make docker-up     - Запуск PostgreSQL в Docker"
	@echo "  make docker-down   - Остановка PostgreSQL"
	@echo "  make test          - Запуск тестов"
	@echo "  make clean         - Очистка сгенерированных файлов"
	@echo "  make install-grpcurl - Установка grpcurl для тестирования"
	@echo "  make test-register - Быстрый тест регистрации"
	@echo "  make test-login    - Быстрый тест входа"

# Установка ВСЕХ зависимостей
deps:
	@echo "📦 Установка Go модулей..."
	go mod download
	go mod tidy
	@echo ""
	@echo "🔧 Установка protoc плагинов..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo ""
	@echo "✅ Все зависимости установлены!"

# Генерация proto файлов (зависит от deps)
proto: deps
	@echo "⚙️  Генерация gRPC кода из proto файлов..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth.proto
	@echo "✅ gRPC код сгенерирован!"

# Сборка приложения (зависит от proto)
build: proto
	@echo "🔨 Сборка приложения..."
	go build -o bin/auth-server cmd/server/main.go
	@echo "✅ Бинарник создан: bin/auth-server"

# Запуск сервера (зависит от proto)
run: proto
	@echo "🚀 Запуск Auth Service на localhost:50051..."
	go run cmd/server/main.go

# Поднять PostgreSQL в Docker
docker-up:
	@echo "🐳 Запуск PostgreSQL в Docker..."
	@if docker ps -a | grep -q auth-postgres; then \
		echo "⚠️  Контейнер auth-postgres уже существует. Удаляю..."; \
		docker rm -f auth-postgres; \
	fi
	docker run -d \
		--name auth-postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=auth_db \
		-p 5432:5432 \
		postgres:15-alpine
	@echo "⏳ Ожидание инициализации БД (5 секунд)..."
	@sleep 5
	@echo "✅ PostgreSQL запущен на localhost:5432"

# Остановить PostgreSQL
docker-down:
	@echo "🛑 Остановка PostgreSQL..."
	docker stop auth-postgres || true
	docker rm auth-postgres || true
	@echo "✅ PostgreSQL остановлен"

# Остановить сервер (Ctrl+C в терминале где запущен make run)
stop:
	@echo "Нажмите Ctrl+C в терминале где запущен сервер"

# Тесты
test:
	@echo "🧪 Запуск тестов..."
	go test -v ./...

# Очистка
clean:
	@echo "🧹 Очистка..."
	rm -rf bin/
	rm -f proto/*.pb.go
	go clean
	@echo "✅ Очистка завершена"

# Установка grpcurl для тестирования
install-grpcurl:
	@echo "📥 Установка grpcurl..."
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	@echo "✅ grpcurl установлен"

# Быстрый тест регистрации
test-register:
	@echo "🧪 Тест регистрации пользователя..."
	grpcurl -plaintext -d '{"email":"test@example.com","password":"password123","username":"testuser"}' \
		localhost:50051 auth.AuthService/Register

# Быстрый тест входа
test-login:
	@echo "🧪 Тест входа пользователя..."
	grpcurl -plaintext -d '{"email":"test@example.com","password":"password123"}' \
		localhost:50051 auth.AuthService/Login

# Список доступных методов
list-methods:
	@echo "📋 Список методов Auth Service:"
	grpcurl -plaintext localhost:50051 list auth.AuthService

# Полная установка и запуск (для первого раза)
all: deps proto docker-up
	@echo ""
	@echo "========================================="
	@echo "✅ Всё готово! Теперь запустите сервер:"
	@echo "   make run"
	@echo ""
	@echo "Для тестирования:"
	@echo "   make install-grpcurl"
	@echo "   make test-register"
	@echo "   make test-login"
	@echo "========================================="

# Docker Compose (альтернатива docker-up)
compose-up:
	@echo "🐳 Запуск через docker-compose..."
	docker-compose up -d
	@echo "✅ Сервисы запущены"

compose-down:
	@echo "🛑 Остановка docker-compose..."
	docker-compose down
	@echo "✅ Сервисы остановлены"

# Проверка статуса
status:
	@echo "📊 Статус сервисов:"
	@echo ""
	@echo "PostgreSQL:"
	@docker ps | grep auth-postgres || echo "  ❌ Не запущен"
	@echo ""
	@echo "gRPC Server:"
	@if lsof -i :50051 > /dev/null 2>&1; then \
		echo "  ✅ Запущен на localhost:50051"; \
	else \
		echo "  ❌ Не запущен"; \
	fi