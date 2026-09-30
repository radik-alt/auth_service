.PHONY: deps proto swagger clean help

PROTO_DEPS=proto_deps
PROTO_DIR=pkg/user/v1
GRPC_PROTO=$(PROTO_DIR)/user.proto

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
		--openapiv2_opt=allow_merge=true,merge_file_name=user \
		$(GRPC_PROTO)
	@echo "✅ Done!"
	@echo "📚 Swagger: $(PROTO_DIR)/user.swagger.json"

# Только Swagger
swagger: deps
	@echo "📚 Generating only Swagger..."
	protoc -I . -I $(PROTO_DEPS) \
		--openapiv2_out=$(PROTO_DIR) \
		--openapiv2_opt=allow_merge=true,merge_file_name=user \
		$(GRPC_PROTO)
	@echo "✅ Swagger: $(PROTO_DIR)/user.swagger.json"

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

# Помощь
help:
	@echo "Commands:"
	@echo "  make deps       - Download proto dependencies"
	@echo "  make proto      - Generate proto + Swagger"
	@echo "  make swagger    - Generate only Swagger"
	@echo "  make clean      - Remove generated files"
	@echo "  make clean-all  - Remove everything including deps"