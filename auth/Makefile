.PHONY: deps proto swagger clean clean-all help redis-up redis-down redis-rm redis-cli redis-logs

PROTO_DEPS=proto_deps
REDIS_CONTAINER=auth-redis
REDIS_PORT=6379
PROTO_DIR=pkg/auth/v1
GRPC_PROTO=$(PROTO_DIR)/auth.proto

# Скачивание зависимостей
deps:
	@echo "📦 Downloading proto dependencies..."
	@mkdir -p $(PROTO_DEPS)/google/api
	@if [ ! -f $(PROTO_DEPS)/google/api/annotations.proto ]; then \
		echo "📥 Downloading annotations.proto..."; \
		curl -s -o $(PROTO_DEPS)/google/api/annotations.proto \
			https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto; \
	fi
	@if [ ! -f $(PROTO_DEPS)/google/api/http.proto ]; then \
		echo "📥 Downloading http.proto..."; \
		curl -s -o $(PROTO_DEPS)/google/api/http.proto \
			https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto; \
	fi
	@echo "✅ Dependencies ready!"

# Генерация proto + Swagger
proto: deps
	@echo "🔧 Generating proto and Swagger..."
	protoc -I . -I $(PROTO_DEPS) \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=$(PROTO_DIR) \
		--openapiv2_opt=allow_merge=true,merge_file_name=auth \
		$(GRPC_PROTO)
	@echo "✅ Done!"
	@echo "📚 Swagger: $(PROTO_DIR)/auth.swagger.json"

# Только Swagger
swagger: deps
	@echo "📚 Generating only Swagger..."
	protoc -I . -I $(PROTO_DEPS) \
		--openapiv2_out=$(PROTO_DIR) \
		--openapiv2_opt=allow_merge=true,merge_file_name=auth \
		$(GRPC_PROTO)
	@echo "✅ Swagger: $(PROTO_DIR)/auth.swagger.json"

# Очистка
clean:
	@echo "🧹 Cleaning..."
	rm -f $(PROTO_DIR)/*.pb.go
	rm -f $(PROTO_DIR)/*.pb.gw.go
	rm -f $(PROTO_DIR)/*.swagger.json
	@echo "✅ Cleaned!"

# Очистка всего (включая зависимости)
clean-all: clean
	@echo "🧹 Removing dependencies..."
	rm -rf $(PROTO_DEPS)
	@echo "✅ All cleaned!"

# Поднять Redis (создаёт контейнер, если его ещё нет, иначе просто стартует)
redis-up:
	@if [ -n "$$(docker ps -q -f name=^/$(REDIS_CONTAINER)$$)" ]; then \
		echo "✅ Redis already running"; \
	elif [ -n "$$(docker ps -aq -f name=^/$(REDIS_CONTAINER)$$)" ]; then \
		echo "🔄 Starting existing Redis container..."; \
		docker start $(REDIS_CONTAINER); \
	else \
		echo "🚀 Creating Redis container..."; \
		docker run -d --name $(REDIS_CONTAINER) \
			-p $(REDIS_PORT):6379 \
			-v $(REDIS_CONTAINER)-data:/data \
			redis:7-alpine redis-server --appendonly yes; \
	fi
	@echo "✅ Redis on localhost:$(REDIS_PORT)"

# Остановить Redis (данные сохраняются)
redis-down:
	@docker stop $(REDIS_CONTAINER) 2>/dev/null || true

# Остановить и удалить контейнер Redis (volume с данными остаётся)
redis-rm: redis-down
	@docker rm $(REDIS_CONTAINER) 2>/dev/null || true

# Открыть redis-cli внутри контейнера
redis-cli:
	@docker exec -it $(REDIS_CONTAINER) redis-cli

# Логи Redis
redis-logs:
	@docker logs -f $(REDIS_CONTAINER)

# Помощь
help:
	@echo "Commands:"
	@echo "  make deps       - Download proto dependencies"
	@echo "  make proto      - Generate proto + Swagger"
	@echo "  make swagger    - Generate only Swagger"
	@echo "  make clean      - Remove generated files"
	@echo "  make clean-all  - Remove everything including deps"
	@echo "  make redis-up   - Start Redis container (creates it on first run)"
	@echo "  make redis-down - Stop Redis container"
	@echo "  make redis-rm   - Stop and remove Redis container"
	@echo "  make redis-cli  - Open redis-cli in the running container"
	@echo "  make redis-logs - Tail Redis container logs"